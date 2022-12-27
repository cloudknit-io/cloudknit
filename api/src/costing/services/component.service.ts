import { BadRequestException, Injectable, InternalServerErrorException, Logger, NotFoundException } from "@nestjs/common";
import { InjectRepository } from "@nestjs/typeorm";
import { Subject } from "rxjs";
import { Component } from "src/typeorm/component.entity";
import { Repository } from "typeorm";
import { CostingDto } from "../dtos/Costing.dto";
import { StreamDto as StreamDto } from "../streams/costing.stream";
import { MessageEvent } from "@nestjs/common";
import { Organization } from "src/typeorm";
import { SqlErrorCodes } from "src/types";
import { Environment } from "src/typeorm/reconciliation/environment.entity";
import { ComponentDto } from "../dtos/Component.dto";

@Injectable()
export class ComponentService {
  readonly stream: Subject<{}> = new Subject<{}>();
  readonly notifyStream: Subject<MessageEvent> = new Subject<MessageEvent>();
  private readonly logger = new Logger(ComponentService.name);
  
  constructor(
    @InjectRepository(Component)
    private componentRepository: Repository<Component>,
    @InjectRepository(Environment)
    private envRepo: Repository<Environment>,
  ) {
    setInterval(() => {
      this.notifyStream.next({ data: {} });
    }, 20000);
  }

  async getAll(org: Organization): Promise<Component[]> {
    const components = await this.componentRepository
      .createQueryBuilder()
      .where('organizationId = :orgId and isDestroyed = 0', {
        orgId: org.id
      })
      .getMany();

    return components;
  }
  
  async getEnvironmentCost(
    org: Organization,
    teamName: string,
    environmentName: string,
  ): Promise<number> {
    const env = await this.envRepo.findOne({
      where: {
        name: environmentName,
        teamName,
        organization: {
          id: org.id
        }
      }
    });

    if (!env) {
      this.logger.error(`could not find environment [${environmentName}] for org [${org.id} / ${org.name}]`);
      throw new NotFoundException();
    }

    const raw = await this.componentRepository
      .createQueryBuilder("components")
      .select("SUM(components.estimated_cost) as cost")
      .where(
        `components.estimated_cost != -1 and 
        components.teamName = '${teamName}' and 
        components.environmentId = '${env.id}' and 
        components.isDestroyed = 0 and 
        components.organizationId = ${org.id}`
      )
      .getRawOne();
    return Number(raw.cost || 0);
  }

  async getEnvironment(
    org: Organization,
    teamName: string,
    environmentName: string,
  ): Promise<Environment> {
    const env = await this.envRepo.findOne({
      where: {
        name: environmentName,
        teamName,
        organization: {
          id: org.id
        }
      },
      relations: {
        components: {
          environment: false
        }
      }
    });

    if (!env) {
      this.logger.error(`could not find environment [${environmentName}] for org [${org.id} / ${org.name}]`);
      throw new NotFoundException();
    }

    return env;
  }

  async getComponentCost(org: Organization, compName: string, teamName: string, envName: string): Promise<ComponentDto> {
    const component = await this.componentRepository.findOne({
      where: {
        componentName: compName,
        teamName,
        environment: {
          name: envName
        },
        organization: {
          id: org.id
        }
      }
    });

    if (!component) {
      this.logger.error({message: 'could not find component', teamName, compName, envName})
      throw new NotFoundException();
    }

    return {
      id: component.id,
      estimatedCost: component.isDestroyed ? 0 : component.estimatedCost,
      status: component.status,
      duration: component.duration,
      componentName: component.componentName,
      lastReconcileDatetime: component.lastReconcileDatetime,
      costResources: component.costResources,
    };
  }

  async getTeamCost(org: Organization, name: string): Promise<number> {
    const raw = await this.componentRepository
      .createQueryBuilder("components")
      .select("SUM(components.estimated_cost) as cost")
      .where(
        `components.teamName = '${name}' and 
        components.isDestroyed = 0 and 
        components.estimated_cost != -1 and
        components.organizationId = ${org.id}`
      )
      .getRawOne();

    if (!raw) {
      this.logger.error(`could not find cost for team ${name} for org ${org.name}`)
      throw new NotFoundException();
    }

    return Number(raw.cost || 0);
  }

  async saveOrUpdate(org: Organization, costing: CostingDto): Promise<boolean> {
    const id = `${costing.teamName}-${costing.environmentName}-${costing.component.componentName}`;

    const env = await this.envRepo.findOne({
      where: {
        name: costing.environmentName,
        teamName: costing.teamName,
        organization: {
          id: org.id
        }
      }
    });

    if (!env) {
      throw new BadRequestException(`could not find environment associated with this component ${id}`);
    }

    let savedComponent = (await this.componentRepository
      .createQueryBuilder('component')
      .leftJoinAndSelect('component.environment', 'environment')
      .leftJoinAndSelect('component.organization', 'organization')
      .where('component.organizationId = :orgId and component.id = :name', {
        orgId: org.id,
        name: id
      })
      .getOne()) || null;
    
    if (costing.component.isDestroyed && savedComponent) {
      // Update existing component
      if (costing.component.status && savedComponent.status !== costing.component.status) {
	      savedComponent.status = costing.component.status;
      }
      savedComponent = await this.softDelete(savedComponent);
    }
    else if (savedComponent) {
      savedComponent.isDestroyed = costing.component.isDestroyed ?? savedComponent.isDestroyed;
      // Update existing component
      if (costing.component.status && savedComponent.status !== costing.component.status) {
        savedComponent.status = costing.component.status;
      }
  
      if (costing.component.cost !== undefined && savedComponent.estimatedCost !== costing.component.cost) {
        savedComponent.estimatedCost = costing.component.cost;
        savedComponent.costResources = costing.component.resources;
      }

      if (costing.component.duration != undefined && savedComponent.duration !== costing.component.duration) {
        savedComponent.duration = costing.component.duration;
      }

      savedComponent = await this.componentRepository.save(savedComponent);
    } else {
      // Create new component
      const component = new Component();
      component.organization = org;
      component.teamName = costing.teamName;
      component.environment = env;
      component.id = id;
      component.status = costing.component.status;
      component.componentName = costing.component.componentName;
      component.estimatedCost = costing.component.cost;
      component.isDestroyed = costing.component.isDestroyed;

      try {
        savedComponent = await this.componentRepository.save(component);
      } catch (err) {
        if (err.code === SqlErrorCodes.NO_DEFAULT) {
          throw new BadRequestException(err.sqlMessage);
        }

        this.logger.error({ message: 'error saving new component', component, error: err.message });
        throw new InternalServerErrorException();
      }
    }

    savedComponent.environment = env;

    const data = await this.getComponentSteamDto(org, savedComponent);
    this.logger.debug({message: 'notifying cost stream', data})
    this.notifyStream.next({ data });

    return true;
  }

  async getComponentSteamDto(org: Organization, component: Component): Promise<StreamDto> {
    const data: StreamDto = {
      team: {
        teamId: component.teamName,
        cost: await this.getTeamCost(org, component.teamName),
      },
      environment: {
        environmentId: component.environment.name,
        cost: await this.getEnvironmentCost(
          org,
          component.teamName,
          component.environment.name
        ),
      },
      component: {
        id: component.id,
        estimatedCost: component.isDestroyed ? 0 : component.estimatedCost,
        status: component.status,
        duration: component.duration,
        componentName: component.componentName,
        lastReconcileDatetime: component.lastReconcileDatetime,
        costResources: component.costResources,
      },
    };
    
    return data;
  }

  async softDelete(component: Component): Promise<Component> {
    component.isDestroyed = true;
    return await this.componentRepository.save(component);
  }
}

import { BadRequestException, Injectable, Logger, NotFoundException } from "@nestjs/common";
import { InjectRepository } from "@nestjs/typeorm";
import { Subject } from "rxjs";
import { Component } from "src/typeorm/component.entity";
import { Resource } from "src/typeorm/resources/Resource.entity";
import { Repository } from "typeorm";
import { CostingDto } from "../dtos/Costing.dto";
import { StreamDto as StreamDto } from "../streams/costing.stream";
import { Mapper } from "../utilities/mapper";
import { MessageEvent } from "@nestjs/common";
import { Organization } from "src/typeorm";
import { Environment } from "src/typeorm/reconciliation/environment.entity";
import { EnvironmentDto } from "../dtos/Environment.dto";

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
    @InjectRepository(Resource)
    private resourceRepository: Repository<Resource>
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
    fullEnvName?: string
  ): Promise<number> {
    const env = await this.envRepo
      .createQueryBuilder()
      .where('name = :envName and organizationId = :orgId', {
        envName: fullEnvName ? fullEnvName : `${teamName}-${environmentName}`,
        orgId: org.id
      })
      .getOne();

    if (!env) {
      this.logger.error(`could not find environment [${environmentName}] for org [${org.id} / ${org.name}]`);
      throw new NotFoundException();
    }

    const raw = await this.componentRepository
      .createQueryBuilder("components")
      .select("SUM(components.cost) as cost")
      .where(
        `components.cost != -1 and 
        components.teamName = '${teamName}' and 
        components.environmentId = '${env.id}' and 
        components.isDestroyed = 0 and 
        components.organizationId = ${org.id}`
      )
      .getRawOne();
    return Number(raw.cost || 0);
  }

  async getComponentCost(org: Organization, componentId: string): Promise<number> {
    const raw = await this.componentRepository
      .createQueryBuilder("components")
      .select("SUM(components.cost) as cost")
      .where(
        `components.id = '${componentId}' and 
        components.isDestroyed = 0 and
        components.organizationId = ${org.id}`)
      .getRawOne();

    if (!raw) {
      this.logger.error(`could not find cost for component ${componentId} for org ${org.name}`)
      throw new NotFoundException();
    }

    return Number(raw.cost || 0);
  }

  async getTeamCost(org: Organization, name: string): Promise<number> {
    const raw = await this.componentRepository
      .createQueryBuilder("components")
      .select("SUM(components.cost) as cost")
      .where(
        `components.teamName = '${name}' and 
        components.isDestroyed = 0 and 
        components.cost != -1 and
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

    const env = await this.envRepo
      .createQueryBuilder()
      .where('name = :envName and organizationId = :orgId', {
        envName: `${costing.teamName}-${costing.environmentName}`,
        orgId: org.id
      })
      .getOne();

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
      savedComponent = await this.softDelete(savedComponent);
    }
    else if (savedComponent) {
      // Update existing component
      if (savedComponent.status !== costing.component.status) {
        savedComponent.status = costing.component.status;
      }
  
      if (savedComponent.cost !== costing.component.cost) {
        savedComponent.cost = costing.component.cost;
      }

      if (savedComponent.duration !== costing.component.duration) {
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
      component.cost = costing.component.cost;
      component.isDestroyed = costing.component.isDestroyed;

      savedComponent = await this.componentRepository.save(component);
    }

    const resources = await this.resourceRepository.save(
      Mapper.mapToResourceEntity(org, savedComponent, costing.component.resources)
    );

    savedComponent.environment = env;
    savedComponent.resources = resources;

    const data = await this.getComponentSteamDto(org, savedComponent);
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
          component.environment.name,
          // TECH DEBT: this contains the full `teamName-envName` environment name.
          // this is fixed by fixing our data model by creating a `teams` table with proper
          // relationships between teams/environment/component tables.
          // Currently, the `name` column on `environment` table is formatted like : 'teamName-envName'.
          // it should just be 'envName'
          component.environment.name 
        ),
      },
      component: {
        id: component.id,
        cost: component.isDestroyed ? 0 : component.cost,
        status: component.status,
        duration: component.duration,
        componentName: component.componentName,
        lastReconcileDatetime: component.lastReconcileDatetime,
        resources: component.resources,
      },
    };
    
    return data;
  }

  async softDelete(component: Component): Promise<Component> {
    component.isDestroyed = true;
    return await this.componentRepository.save(component);
  }

  async getResourceData(org: Organization, id: string) {
    const resultSet = await this.resourceRepository.find({
      where: {
        componentId: id,
        organization: {
          id: org.id
        }
      },
    });

    const roots = [];
    var resources = new Map<string, any>(
      resultSet.map((e) => [e.id, { ...e, subresources: [] }])
    );

    for (let i = 0; i < resultSet.length; i++) {
      const resource = resources.get(resultSet[i].id);

      if (!resource.parentId) {
        roots.push(resource);
      } else {
        resources.get(resource.parentId).subresources.push(resource);
      }
    }

    return {
      componentId: id,
      resources: roots,
    };
  }
}

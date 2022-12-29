import { BadRequestException, Injectable, InternalServerErrorException, Logger, NotFoundException } from "@nestjs/common";
import { InjectRepository } from "@nestjs/typeorm";
import { Subject } from "rxjs";
import { Repository } from "typeorm";
import { MessageEvent } from "@nestjs/common";
import { Component, Organization, Team } from "src/typeorm";
import { Environment } from "src/typeorm/environment.entity";
import { ComponentDto } from "../dtos/Component.dto";
import { EnvironmentService } from "src/environment/environment.service";

@Injectable()
export class ComponentService {
  readonly stream: Subject<{}> = new Subject<{}>();
  readonly notifyStream: Subject<MessageEvent> = new Subject<MessageEvent>();
  private readonly logger = new Logger(ComponentService.name);
  
  constructor(
    @InjectRepository(Component)
    private compRepo: Repository<Component>,
    private readonly envSvc: EnvironmentService,
  ) {
    setInterval(() => {
      this.notifyStream.next({ data: {} });
    }, 20000);
  }

  async create(org: Organization, env: Environment, name: string): Promise<Component> {
    return this.compRepo.save({
      name,
      environment: env,
      organization: org
    });
  }

  async getAllForEnvironmentById(org: Organization, env: Environment): Promise<Component[]> {
    return this.compRepo.find({
      where: {
        organization: {
          id: org.id
        },
        environment: {
          id: env.id
        }
      }
    })
  }

  async getAll(org: Organization): Promise<Component[]> {
    const components = await this.compRepo.find({
      where: {
        organization: {
          id: org.id
        }
      }
    });

    return components;
  }

  async findById(org: Organization, id: number, relations?: {}): Promise<Component> {
    return await this.compRepo.findOne({
      where: {
        id,
        organization: {
          id: org.id
        }
      },
      relations
    });
  }

  async findByName(org: Organization, env: Environment, name: string, relations?: {}): Promise<Component> {
    return await this.compRepo.findOne({
      where: {
        name,
        environment: {
          id: env.id
        },
        organization: {
          id: org.id
        }
      },
      relations
    });
  }
  
  async getEnvironmentCost(
    org: Organization,
    team: Team,
    environmentName: string,
  ): Promise<number> {
    const env = await this.envSvc.findByName(org, team, environmentName);

    if (!env) {
      this.logger.error({ message: 'could not find environment', environmentName, team });
      throw new NotFoundException();
    }

    const raw = await this.compRepo
      .createQueryBuilder("components")
      .select("SUM(components.estimated_cost) as cost")
      .where(
        `components.estimated_cost != -1 and 
        components.teamName = '${team.name}' and 
        components.environmentId = '${env.id}' and 
        components.organizationId = ${org.id}`
      )
      .getRawOne();
    return Number(raw.cost || 0);
  }

  async getCostByTeam(org: Organization, team: Team): Promise<number> {
    const raw = await this.compRepo
      .createQueryBuilder("components")
      .select("SUM(components.estimated_cost) as cost")
      .where(
        `components.teamName = '${team.name}' and 
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

  async getComponentCost(org: Organization, team: Team, compName: string, envName: string): Promise<ComponentDto> {
    const env = await this.envSvc.findByName(org, team, envName);
    const component = await this.findByName(org, env, compName);

    if (!component) {
      this.logger.error({message: 'could not find component', teamName: team.name, compName, envName})
      throw new NotFoundException();
    }

    return {
      id: component.id,
      estimatedCost: component.isDestroyed ? 0 : component.estimatedCost,
      status: component.status,
      duration: component.duration,
      name: component.name,
      lastReconcileDatetime: component.lastReconcileDatetime,
      costResources: component.costResources,
    };
  }

  // async saveOrUpdate(org: Organization, costing: CostingDto): Promise<boolean> {
  //   const id = `${costing.teamName}-${costing.environmentName}-${costing.component.componentName}`;

  //   const env = await this.envRepo.findOne({
  //     where: {
  //       name: costing.environmentName,
  //       teamName: costing.teamName,
  //       organization: {
  //         id: org.id
  //       }
  //     }
  //   });

  //   if (!env) {
  //     throw new BadRequestException(`could not find environment associated with this component ${id}`);
  //   }

  //   let savedComponent = (await this.compRepo
  //     .createQueryBuilder('component')
  //     .leftJoinAndSelect('component.environment', 'environment')
  //     .leftJoinAndSelect('component.organization', 'organization')
  //     .where('component.organizationId = :orgId and component.id = :name', {
  //       orgId: org.id,
  //       name: id
  //     })
  //     .getOne()) || null;
    
  //   if (savedComponent) {
  //     savedComponent.isDestroyed = costing.component.isDestroyed ?? savedComponent.isDestroyed;
  //     if (savedComponent.isDestroyed) {
  //       savedComponent.estimatedCost = -1;
  //     }
  //     // Update existing component
  //     if (costing.component.status && savedComponent.status !== costing.component.status) {
  //       savedComponent.status = costing.component.status;
  //     }
  
  //     if (costing.component.cost !== undefined && savedComponent.estimatedCost !== costing.component.cost) {
  //       savedComponent.estimatedCost = costing.component.cost;
  //       savedComponent.costResources = costing.component.resources;
  //     }

  //     if (costing.component.duration != undefined && savedComponent.duration !== costing.component.duration) {
  //       savedComponent.duration = costing.component.duration;
  //     }

  //     savedComponent = await this.compRepo.save(savedComponent);
  //   } else {
  //     // Create new component
  //     const component = new Component();
  //     component.organization = org;
  //     component.teamName = costing.teamName;
  //     component.environment = env;
  //     component.id = id;
  //     component.status = costing.component.status;
  //     component.componentName = costing.component.componentName;
  //     component.estimatedCost = costing.component.cost;
  //     component.isDestroyed = costing.component.isDestroyed;

  //     try {
  //       savedComponent = await this.compRepo.save(component);
  //     } catch (err) {
  //       if (err.code === SqlErrorCodes.NO_DEFAULT) {
  //         throw new BadRequestException(err.sqlMessage);
  //       }

  //       this.logger.error({ message: 'error saving new component', component, error: err.message });
  //       throw new InternalServerErrorException();
  //     }
  //   }

  //   savedComponent.environment = env;

  //   const data = await this.getComponentSteamDto(org, savedComponent);
  //   this.logger.debug({message: 'notifying cost stream', data})
  //   this.notifyStream.next({ data });

  //   return true;
  // }

  // async getComponentSteamDto(org: Organization, component: Component): Promise<StreamDto> {
  //   const data: StreamDto = {
  //     team: {
  //       teamId: component.teamName,
  //       cost: await this.getTeamCost(org, component.teamName),
  //     },
  //     environment: {
  //       environmentId: `${component.teamName}-${component.environment.name}`,
  //       cost: await this.getEnvironmentCost(
  //         org,
  //         component.teamName,
  //         component.environment.name
  //       ),
  //     },
  //     component: {
  //       id: component.id,
  //       estimatedCost: component.estimatedCost,
  //       status: component.status,
  //       duration: component.duration,
  //       componentName: component.componentName,
  //       lastReconcileDatetime: component.lastReconcileDatetime,
  //       costResources: component.costResources,
  //       isDestroyed: component.isDestroyed
  //     },
  //   };
    
  //   return data;
  // }

  async softDelete(component: Component): Promise<Component> {
    component.isDestroyed = true;
    component.estimatedCost = -1;
    return await this.compRepo.save(component);
  }
}

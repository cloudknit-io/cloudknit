import { BadRequestException, Injectable, Logger, NotFoundException } from "@nestjs/common";
import { InjectRepository } from "@nestjs/typeorm";
import { Subject } from "rxjs";
import { Component } from "src/typeorm/costing/Component";
import { Resource } from "src/typeorm/resources/Resource.entity";
import { Repository } from "typeorm";
import { CostingDto } from "../dtos/Costing.dto";
import { CostingStreamDto } from "../streams/costing.stream";
import { Mapper } from "../utilities/mapper";
import { MessageEvent } from "@nestjs/common";
import { Organization } from "src/typeorm";
import { Environment } from "src/typeorm/reconciliation/environment.entity";

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
      .where('organizationId = :orgId and isDeleted = 0', {
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
        components.isDeleted = 0 and 
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
        components.isDeleted = 0 and
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
        components.isDeleted = 0 and 
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

  async saveComponents(org: Organization, costing: CostingDto): Promise<boolean> {
    const id = `${costing.teamName}-${costing.environmentName}-${costing.component.componentName}`;

    let savedComponent = (await this.componentRepository.findOne({ 
      where: { 
        id,
        organization: org
      },
      relations: {
        environment: true
      } })) || null;
    
    if (costing.component.isDeleted && savedComponent) {
      savedComponent = await this.softDelete(savedComponent);
    } else {
      const env = await this.envRepo
      .createQueryBuilder()
      .where('name = :envName and organizationId = :orgId', {
        envName: `${costing.teamName}-${costing.environmentName}`,
        orgId: org.id
      })
      .getOne();

      const component = new Component();
      component.organization = org;
      component.teamName = costing.teamName;
      component.environment = env;
      component.id = id;
      component.componentName = costing.component.componentName;
      component.cost = costing.component.cost;
      component.isDeleted = costing.component.isDeleted;

      await this.componentRepository.delete({
        id: id,
        organization: org
      });

      savedComponent = await this.componentRepository.save(component);
      savedComponent.environment = env;

      const resources = await this.resourceRepository.save(
        Mapper.mapToResourceEntity(org, component, costing.component.resources)
      );

      savedComponent.resources = resources;
    }

    const data = await this.getCostingStreamDto(savedComponent);
    this.notifyStream.next({ data });

    return true;
  }

  async getCostingStreamDto(component: Component): Promise<CostingStreamDto> {
    const data: CostingStreamDto = {
      team: {
        teamId: component.teamName,
        cost: await this.getTeamCost(component.organization, component.teamName),
      },
      environment: {
        environmentId: `${component.teamName}-${component.environment.name}`,
        cost: await this.getEnvironmentCost(
          component.organization,
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
        componentId: component.id,
        cost: component.isDeleted ? 0 : component.cost,
      },
    };
    
    return data;
  }

  async softDelete(component: Component): Promise<Component> {
    component.isDeleted = true;
    return await this.componentRepository.save(component);
  }

  async getResourceData(org: Organization, id: string) {
    const resultSet = await this.resourceRepository.find({
      where: {
        componentId: id,
        organization: org
      },
    });

    const roots = [];
    const resources = new Map<string, any>(
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

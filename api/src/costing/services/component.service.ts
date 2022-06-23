import { Injectable } from "@nestjs/common";
import { InjectRepository } from "@nestjs/typeorm";
import { Subject } from "rxjs";
import { Component } from "src/typeorm/costing/entities/Component";
import { CostComponent, Resource } from "src/typeorm/resources/Resource.entity";
import { Connection, Repository } from "typeorm";
import { CostingDto } from "../dtos/Costing.dto";
import { CostingStreamDto } from "../streams/costing.stream";
import { Mapper } from "../utilities/mapper";
import { MessageEvent } from "@nestjs/common";

@Injectable()
export class ComponentService {
  readonly stream: Subject<{}> = new Subject<{}>();
  readonly notifyStream: Subject<MessageEvent> = new Subject<MessageEvent>();
  constructor(
    private readonly connection: Connection,
    @InjectRepository(Component)
    private componentRepository: Repository<Component>,
    @InjectRepository(Resource)
    private resourceRepository: Repository<Resource>,
    @InjectRepository(CostComponent)
    private costComponentRepository: Repository<CostComponent>
  ) {
    setInterval(() => {
      this.notifyStream.next({ data: {} });
    }, 20000);
  }

  async getAll(): Promise<Component[]> {
    const components = await this.componentRepository.find({
      where: {
        isDeleted: false,
      },
    });
    return components;
  }
  async getEnvironmentCost(
    teamName: string,
    environmentName: string
  ): Promise<number> {
    const raw = await this.componentRepository
      .createQueryBuilder("components")
      .select("SUM(components.cost) as cost")
      .where(
        `components.cost != -1 and components.teamName = '${teamName}' and components.environmentName = '${environmentName}' and components.isDeleted = 0`
      )
      .getRawOne();
    return Number(raw.cost || 0);
  }

  async getComponentCost(componentId: string): Promise<number> {
    const raw = await this.componentRepository
      .createQueryBuilder("components")
      .select("SUM(components.cost) as cost")
      .where(`components.id = '${componentId}' and components.isDeleted = 0`)
      .getRawOne();
    return Number(raw.cost || 0);
  }

  async getTeamCost(name: string): Promise<number> {
    const raw = await this.componentRepository
      .createQueryBuilder("components")
      .select("SUM(components.cost) as cost")
      .where(
        `components.teamName = '${name}' and components.isDeleted = 0 and components.cost != -1`
      )
      .getRawOne();
    return Number(raw.cost || 0);
  }

  async saveComponents(costing: CostingDto): Promise<boolean> {
    const id = `${costing.teamName}-${costing.environmentName}-${costing.component.componentName}`;
    let savedComponent = (await this.componentRepository.findOne({ where: { id } })) || null;
    if (costing.component.isDeleted && savedComponent) {
      savedComponent = await this.softDelete(savedComponent);
    } else {
      const component = new Component();
      component.teamName = costing.teamName;
      component.environmentName = costing.environmentName;
      component.id = id;
      component.componentName = costing.component.componentName;
      component.cost = costing.component.cost;
      component.isDeleted = costing.component.isDeleted;
      await this.componentRepository.delete({
        id: id,
      });
      savedComponent = await this.componentRepository.save(component);
      const resources = await this.resourceRepository.save(
        Mapper.mapToResourceEntity(component, costing.component.resources)
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
        cost: await this.getTeamCost(component.teamName),
      },
      environment: {
        environmentId: `${component.teamName}-${component.environmentName}`,
        cost: await this.getEnvironmentCost(
          component.teamName,
          component.environmentName
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

  async getResourceData(id: string) {
    const resultSet = await this.resourceRepository.find({
      where: {
        componentId: id,
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

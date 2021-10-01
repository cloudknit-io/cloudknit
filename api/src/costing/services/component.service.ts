import { Injectable } from "@nestjs/common";
import { InjectRepository } from "@nestjs/typeorm";
import { Subject } from "rxjs";
import { Component } from "src/typeorm/costing/entities/Component";
import { CostComponent, Resource } from "src/typeorm/resources/Resource.entity";
import { Connection, Repository } from "typeorm";
import { CostingDto } from "../dtos/Costing.dto";
import { Mapper } from "../utilities/mapper";

@Injectable()
export class ComponentService {
  readonly stream: Subject<{}> = new Subject<{}>();
  readonly notifyStream: Subject<{}> = new Subject<{}>();
  constructor(
    private readonly connection: Connection,
    @InjectRepository(Component)
    private componentRepository: Repository<Component>,
    @InjectRepository(Resource)
    private resourceRepository: Repository<Resource>,
    @InjectRepository(CostComponent)
    private costComponentRepository: Repository<CostComponent>
  ) {}

  async getAll(): Promise<Component[]> {
    const components = await this.componentRepository.find({
      where: {
        isDeleted: 0,
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
        `components.teamName = '${teamName}' and components.environmentName = '${environmentName}' and components.isDeleted = 0`
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
      .where(`components.teamName = '${name}' and components.isDeleted = 0`)
      .getRawOne();
    return Number(raw.cost || 0);
  }

  async saveComponents(costing: CostingDto): Promise<boolean> {
    const id = `${costing.teamName}-${costing.environmentName}-${costing.component.componentName}`;
    let savedComponent = null;
    if (costing.component.isDeleted) {
      savedComponent = await this.softDelete(id);
      console.log('getting saved component upon delete ---------> ', JSON.stringify(savedComponent));
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
    this.notifyStream.next(savedComponent);
    return true;
  }
  
  async softDelete(id: string): Promise<Component> {
    const component = this.componentRepository.findOne(id);
    return await this.componentRepository.save({
      ...component,
      id: id,
      isDeleted: true,
    });
  }

  async getResourceData(id: string) {
    const resultSet = await this.resourceRepository.find({
      where: {
        componentId: id,
      },
    });
    const roots = [];
    const resources = new Map<string, any>();
    for (let i = 0; i < resultSet.length; i++) {
      resultSet[i].subresources = [];
      const resource = resultSet[i];
      if (!resultSet[i].parentId) {
        roots.push(resource);
        resources.set(resource.id, resource);
      } else {
        resources.set(resource.id, resource);
        resources.get(resource.parentId).subresources.push(resource);
      }
    }
    return {
      componentId: id,
      resources: roots,
    };
  }

  async execute(query: string) {
    return await this.connection.query(query);
  }
}

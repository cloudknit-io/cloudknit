import { Injectable } from '@nestjs/common'
import { InjectRepository } from '@nestjs/typeorm'
import { Subject } from 'rxjs'
import { Component } from 'src/typeorm/costing/entities/Component'
import { Resource } from 'src/typeorm/resources/Resource.entity'
import { Repository } from 'typeorm'
import { CostingDto } from '../dtos/Costing.dto'
import { Mapper } from '../utilities/mapper'

@Injectable()
export class ComponentService {
  private readonly recursiveResourceQuery = (
    componentId: string,
  ) => `with recursive cte (name, hourlyCost, monthlyCost, resourceName) as (
    select     name,
               hourlyCost, 
               monthlyCost,
               resourceName
    from       resources
    where      componentId = '${componentId}'
    union all
    select     r.name,
               r.hourlyCost,
               r.monthlyCost,
               r.resourceName
    from       resources r
    inner join cte
            on r.resourceName = cte.name
  )
  select * from cte;`
  
  readonly stream: Subject<{}> = new Subject<{}>()
  readonly notifyStream: Subject<{}> = new Subject<{}>()
  constructor(
    @InjectRepository(Component)
    private componentRepository: Repository<Component>,
    @InjectRepository(Resource)
    private resourceRepository: Repository<Resource>,
  ) {}

  async getAll(): Promise<{}> {
    const components = await this.componentRepository.find()
    this.stream.next(components)
    return components
  }
  async getEnvironmentCost(
    teamName: string,
    environmentName: string,
  ): Promise<number> {
    const raw = await this.componentRepository
      .createQueryBuilder('components')
      .select('SUM(components.cost) as cost')
      .where(
        `components.teamName = '${teamName}' and components.environmentName = '${environmentName}'`,
      )
      .getRawOne()
    return Number(raw.cost || 0)
  }

  async getComponentCost(componentId: string): Promise<number> {
    const raw = await this.componentRepository
      .createQueryBuilder('components')
      .select('SUM(components.cost) as cost')
      .where(`components.id = '${componentId}'`)
      .getRawOne()
    return Number(raw.cost || 0)
  }

  async getTeamCost(name: string): Promise<number> {
    const raw = await this.componentRepository
      .createQueryBuilder('components')
      .select('SUM(components.cost) as cost')
      .where(`components.teamName = '${name}'`)
      .getRawOne()
    return Number(raw.cost || 0)
  }

  async saveComponents(costing: CostingDto): Promise<boolean> {
    const entry: Component = {
      teamName: costing.teamName,
      environmentName: costing.environmentName,
      id: `${costing.teamName}-${costing.environmentName}-${costing.component.componentName}`,
      componentName: costing.component.componentName,
      cost: costing.component.cost,
      resources: costing.component.resources,
    }
    const savedComponent = await this.componentRepository.save(entry)
    this.notifyStream.next(savedComponent)
    return true
  }

  async getResourceData(id: string) {
    const resultSet = await this.resourceRepository.query(
      this.recursiveResourceQuery(id),
    )
    const roots = []
    const resources = new Map<string, Resource>()
    for (let i = 0; i < resultSet.length; i++) {
      const resource = Mapper.getResource(resultSet[i])
      if (!resultSet[i].resourceName) {
        roots.push(resource)
        resources.set(resource.name, resource)
      } else {
        resources.set(resource.name, resource)
        resources.get(resource.resourceName).subresources.push(resource)
      }
    }
    return {
      componentId: id,
      resources: roots,
    }
  }
}

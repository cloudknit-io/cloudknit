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

  private readonly resourceQuery = (componentId: string) => {
    return ''.concat(
      'select r1.name as n0, r2.name as n1, r3.name as n2, r4.name as n3, r5.name as n4 ',
      'from resources r1 left join resources r2 on (r1.name = r2.resourceName) ',
      'left join resources r3 on (r2.name = r3.resourceName) ',
      'left join resources r4 on (r3.name = r4.resourceName) ',
      'left join resources r5 on (r4.name = r5.resourceName) ',
      `where r1.componentId = "${componentId}";`,
    )
  }

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
    // return await this.resourceRepository.query(this.recursiveResourceQuery(id));
    const getResourceInfo = (resourceName: string) =>
      this.resourceRepository.query(
        `Select name, hourlyCost, monthlyCost from resources where name = '${resourceName}'`,
      )
    const query = this.resourceQuery(id)
    const resultSet: [] = (await this.resourceRepository.query(query)) || []
    const resources = new Map<string, Resource>()
    for (let i = 0; i < resultSet.length; i++) {
      const hierarchy = Object.keys(resultSet[i])
        .filter((c) => resultSet[i][c])
        .map((c) => resultSet[i][c])
      let par = null
      if (resources.has(hierarchy[0])) {
        par = resources.get(hierarchy[0])
      } else {
        par = Mapper.getResource(await getResourceInfo(hierarchy[0]))
        resources.set(hierarchy[0], par)
      }
      for (let j = 1; j < hierarchy.length; j++) {
        const data = Mapper.getResource(await getResourceInfo(hierarchy[j]))
        par.subresources.push(data)
        par = data
      }
    }
    
    return {
      componentId: id,
      resources: [...resources.values()],
    }
  }
}

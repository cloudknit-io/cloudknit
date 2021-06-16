import { Injectable } from '@nestjs/common'
import { InjectRepository } from '@nestjs/typeorm'
import { Subject } from 'rxjs'
import { Component } from 'src/typeorm/costing/entities/Component'
import { Repository } from 'typeorm'
import { CostingDto } from '../dtos/Costing.dto'

@Injectable()
export class ComponentService {
  readonly stream: Subject<{}> = new Subject<{}>()
  constructor(
    @InjectRepository(Component)
    private componentRepository: Repository<Component>,
  ) {}

  async getAll(): Promise<{}> {
    const components = await this.componentRepository.find()
    this.stream.next(components)
    return components
  }

  async getEnvironmentCost(teamName: string, environmentName: string): Promise<number> {
    const components = await this.componentRepository.find({
      where: {
        teamName: teamName,
        environmentName: environmentName,
      },
    })
    this.stream.next(components)
    return components.reduce((p, c, _i) => p + Number(c.cost), 0)
  }

  async getComponentCost(componentId: string): Promise<number> {
    const components = await this.componentRepository.find({
      where: {
        id: componentId,
      },
    })
    return components.reduce((p, c, _i) => p + Number(c.cost), 0)
  }

  async getTeamCost(name: string): Promise<number> {
    const components = await this.componentRepository.find({
      where: {
        teamName: name,
      },
    })
    this.stream.next(components)
    return components.reduce((p, c, _i) => p + Number(c.cost), 0)
  }

  async saveComponents(costing: CostingDto): Promise<boolean> {
    const entry: Component = {
      teamName: costing.teamName,
      environmentName: costing.environmentName,
      id: `${costing.teamName}-${costing.environmentName}-${costing.component.componentName}`,
      componentName: costing.component.componentName,
      cost: costing.component.cost,
    }
    const savedComponent = await this.componentRepository.save(entry)
    this.stream.next([savedComponent])
    return true
  }
}

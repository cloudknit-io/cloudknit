import { Injectable } from '@nestjs/common'
import { InjectRepository } from '@nestjs/typeorm'
import { Subject } from 'rxjs'
import { Component } from 'src/typeorm/costing/entities/Component'
import { Repository } from 'typeorm'
import { CostingDto } from '../dtos/Costing.dto'

@Injectable()
export class ComponentService {
  readonly stream: Subject<{}> = new Subject<{}>()
  readonly notifyStream: Subject<{}> = new Subject<{}>()
  constructor(
    @InjectRepository(Component)
    private componentRepository: Repository<Component>,
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
    return Number(raw.cost || 0);
  }

  async getComponentCost(componentId: string): Promise<number> {
    const raw = await this.componentRepository
      .createQueryBuilder('components')
      .select('SUM(components.cost) as cost')
      .where(
        `components.id = '${componentId}'`,
      )
      .getRawOne()
    return Number(raw.cost || 0);
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
    }
    const savedComponent = await this.componentRepository.save(entry)
    this.notifyStream.next(savedComponent)
    return true
  }
}

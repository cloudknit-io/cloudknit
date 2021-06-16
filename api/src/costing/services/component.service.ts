import { Injectable } from '@nestjs/common'
import { InjectRepository } from '@nestjs/typeorm'
import { Component } from 'src/typeorm/costing/entities/Component'
import { Repository } from 'typeorm'
import { CostingDto } from '../dtos/Costing.dto'

@Injectable()
export class ComponentService {
  constructor(
    @InjectRepository(Component)
    private componentRepository: Repository<Component>,
  ) {}

  async getAll(): Promise<{}> {
    const components = await this.componentRepository.find();
    return components;
  }

  async getEnvironmentCost(teamName: string, environmentName: string): Promise<number> {
    const components = await this.componentRepository.find({
      where: {
        teamName: teamName,
        environmentName: environmentName,
      },
    })
    return (components).reduce((p, c, _i) => p + Number(c.cost), 0)
  }

  async getComponentCost(componentId: string): Promise<number> {
    const components = await this.componentRepository.find({
      where: {
        id: componentId,
      },
    })
    return (components).reduce((p, c, _i) => p + Number(c.cost), 0)
  }

  async getTeamCost(name: string): Promise<number> {
    const components = await this.componentRepository.find({
      where: {
        teamName: name,
      },
    })
    return (components).reduce((p, c, _i) => p + Number(c.cost), 0)
  }

  async saveComponents(costing: CostingDto): Promise<boolean> {
    const entry: Component = {
      teamName: costing.teamName,
      environmentName: costing.environmentName,
      id : `${costing.teamName}-${costing.environmentName}-${costing.component.componentName}`,
      componentName: costing.component.componentName,
      cost: costing.component.cost,
    };
    await this.componentRepository.save(entry);
    return true
  }
}

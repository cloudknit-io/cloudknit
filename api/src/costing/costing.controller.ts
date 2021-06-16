import { Body, Controller, Get, Param, Post } from '@nestjs/common'
import { CostingDto } from './dtos/Costing.dto';
import { ComponentService } from './services/component.service'

@Controller({
  path: 'costing/api/v1',
})
export class CostingController {
  constructor(
    private readonly componentService: ComponentService,
  ) {}

  @Get('all')
  async getAll(): Promise<{}> {
    return await this.componentService.getAll();
  }
  
  @Get('team/:name')
  async getTeamCost(@Param('name') name: string): Promise<number> {
    return await this.componentService.getTeamCost(name);
  }
  
  @Get('environment/:teamName/:environmentName')
  async getEnvironmentCost(@Param('teamName') teamName: string, @Param('environmentName') environmentName: string): Promise<number> {
    return await this.componentService.getEnvironmentCost(teamName, environmentName);
  }
  
  @Get('component/:componentId')
  async getComponentCost(@Param('componentId') componentId: string): Promise<number> {
    return await this.componentService.getComponentCost(componentId);
  }

  @Post('saveComponent')
  async saveComponent(@Body() costing: CostingDto): Promise<boolean> {
    return await this.componentService.saveComponents(costing);
  }
}

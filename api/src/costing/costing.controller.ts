import { Body, Controller, Get, Param, Post } from '@nestjs/common'
import { CostingDto } from './dtos/Costing.dto';
import { ComponentService } from './services/component.service'

@Controller({
  path: 'costing',
})
export class CostingController {
  constructor(
    private readonly componentService: ComponentService,
  ) {}

  // @Get('team/:name')
  // async getTeam(@Param('name') name: string): Promise<TeamDto> {
  //   return await this.teamService.getTeam(name)
  // }

  // @Get('environment/:name')
  // async getEnvironment(@Param('name') name: string): Promise<EnvironmentDto> {
  //   return await this.environmentService.getEnvironment(name)
  // }

  @Get('all')
  async getAll(): Promise<{}> {
    return await this.componentService.getAll();
  }

  @Get('team-cost/:name')
  async getComponent(@Param('name') name: string): Promise<number> {
    return await this.componentService.getTeamCost(name);
  }
  
  @Get('env-name/:name')
  async getEnvironmentCost(@Param('name') name: string): Promise<number> {
    return await this.componentService.getEnvironmentCost(name);
  }
  
  @Get('config-cost/:name')
  async getComponentCost(@Param('name') name: string): Promise<number> {
    return await this.componentService.getComponentCost(name);
  }

  @Post('saveComponent')
  async saveComponent(@Body() costing: CostingDto): Promise<boolean> {
    return await this.componentService.saveComponents(costing);
  }
}

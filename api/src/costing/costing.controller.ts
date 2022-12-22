import { Body, Controller, Get, Param, Post, Query, Request } from '@nestjs/common'
import { RequiredQueryValidationPipe, TeamEnvQueryParams } from 'src/reconciliation/validationPipes';
import { ComponentDto } from './dtos/Component.dto';
import { CostingDto } from './dtos/Costing.dto'
import { ComponentService } from './services/component.service'

@Controller({
  version: '1'
})
export class CostingController {
  constructor(private readonly componentService: ComponentService) {}

  @Get('all')
  async getAll(@Request() req): Promise<{}> {
    return await this.componentService.getAll(req.org);
  }

  @Get('team/:name')
  async getTeamCost(@Request() req, @Param('name') name: string): Promise<number> {
    return await this.componentService.getTeamCost(req.org, name);
  }

  @Get('environment')
  async getEnvironmentCost(
    @Request() req,
    @Query(new RequiredQueryValidationPipe()) te: TeamEnvQueryParams
  ): Promise<number> {
    return await this.componentService.getEnvironmentCost(
      req.org,
      te.teamName,
      te.envName,
    )
  }

  @Get('component/:componentId')
  async getComponentCost(
    @Request() req,
    @Param('componentId') componentId: string,
  ): Promise<ComponentDto> {
    return await this.componentService.getComponentCost(req.org, componentId);
  }

  @Post('saveComponent')
  async saveComponent(@Request() req, @Body() costing: CostingDto): Promise<boolean> {
    return await this.componentService.saveOrUpdate(req.org, costing)
  }
}

import { Controller, Get, Query, Request } from '@nestjs/common'
import { RequiredQueryValidationPipe, TeamEnvCompQueryParams, TeamEnvQueryParams } from 'src/reconciliation/validationPipes';
import { ComponentDto } from './dtos/Component.dto';
import { ComponentService } from './services/component.service'

@Controller({
  version: '1'
})
export class CostingController {
  constructor(
    private readonly compSvc: ComponentService,
  ) {}

  @Get('all')
  async getAll(@Request() req): Promise<{}> {
    return await this.compSvc.getAll(req.org);
  }

  @Get('environment')
  async getEnvironmentCost(
    @Request() req,
    @Query(new RequiredQueryValidationPipe()) te: TeamEnvQueryParams
  ): Promise<number> {
    const {org, team } = req;

    return await this.compSvc.getEnvironmentCost(
      org,
      team,
      te.envName,
    )
  }

  @Get('component')
  async getComponentCost(
    @Request() req,
    @Query(new RequiredQueryValidationPipe()) tec: TeamEnvCompQueryParams
  ): Promise<ComponentDto> {
    const {org, team, env } = req;
    return await this.compSvc.getComponentCost(org, team, env, tec.compName);
  }
}

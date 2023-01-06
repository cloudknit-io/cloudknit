import { Controller, Get, Body, Patch, Delete, Request } from '@nestjs/common';
import { EnvironmentService } from './environment.service';
import { UpdateEnvironmentDto } from './dto/update-environment.dto';
import { APIRequest } from 'src/types';
import { ReconciliationService } from 'src/reconciliation/reconciliation.service';

@Controller({
  version: '1'
})
export class EnvironmentController {
  constructor(
    private readonly envSvc: EnvironmentService,
    private readonly reconSvc: ReconciliationService
    ) {}

  @Get()
  async findOne(@Request() req) {
    const {org, team, env} = req;

    return this.envSvc.findById(org, env.id);
  }

  @Get('dag')
  async getDag(@Request() req: APIRequest) {
    const { env } = req;

    return env.dag;
  }

  @Patch()
  async update(@Request() req: APIRequest, @Body() updateEnvDto: UpdateEnvironmentDto) {
    const { org, env } = req;

    return this.envSvc.updateById(org, env.id, updateEnvDto);
  }

  @Delete()
  remove(@Request() req: APIRequest) {
    const { org, env } = req;

    return this.envSvc.remove(org, env.id);
  }

  @Get('audit')
  async getAudits(@Request() req: APIRequest) {
    const { org, env } = req;

    return this.reconSvc.getEnvironmentAuditList(org, env);
  }
}

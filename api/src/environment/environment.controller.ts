import { Controller, Get, Post, Body, Patch, Param, Delete, Request } from '@nestjs/common';
import { EnvironmentService } from './environment.service';
import { UpdateEnvironmentDto } from './dto/update-environment.dto';
import { ComponentService } from 'src/costing/services/component.service';
import { In } from 'typeorm';
import { APIRequest } from 'src/types';

@Controller({
  version: '1'
})
export class EnvironmentController {
  constructor(
    private readonly envSvc: EnvironmentService,
    private readonly compSvc: ComponentService
    ) {}

  @Get()
  async findOne(@Request() req) {
    const {org, team, env} = req;

    // TODO : Cost?
    // this.compSvc.getAllForEnvironmentById(org, env);

    return this.envSvc.findById(org, env.id, team);
  }

  @Get('dag')
  async getDag(@Request() req: APIRequest) {
    const { env } = req;

    return env.dag;
  }

  @Get('cost')
  async getCost(@Request() req: APIRequest) {
    const {org, team, env} = req;
    
    return this.compSvc.getEnvironmentCost(org, team, env.name);
  }

  @Patch()
  async update(@Request() req: APIRequest, @Param('id') id: number, @Body() updateEnvDto: UpdateEnvironmentDto) {
    const {org } = req;

    return this.envSvc.update(org, id, updateEnvDto);
  }

  @Delete()
  remove(@Request() req: APIRequest, @Param('id') id: number) {
    const { org } = req;

    return this.envSvc.remove(org, id);
  }
}

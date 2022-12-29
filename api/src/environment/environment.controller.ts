import { Controller, Get, Post, Body, Patch, Param, Delete, Request } from '@nestjs/common';
import { EnvironmentService } from './environment.service';
import { UpdateEnvironmentDto } from './dto/update-environment.dto';
import { ComponentService } from 'src/costing/services/component.service';

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

    return this.envSvc.findById(org, env.id, team);
  }

  @Get('/cost')
  async getTeamCost(@Request() req) {
    const {org, team, env} = req;
    
    return this.compSvc.getEnvironmentCost(org, team, env.name);
  }

  @Patch()
  async update(@Request() req, @Param('id') id: number, @Body() updateEnvDto: UpdateEnvironmentDto) {
    const {org, team, env} = req;

    return this.envSvc.update(org, id, updateEnvDto);
  }

  @Delete()
  remove(@Request() req, @Param('id') id: number) {
    const {org, team, env} = req;

    return this.envSvc.remove(org, id);
  }
}

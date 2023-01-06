import { Controller, Get, Body, Patch, Param, Delete, Request, Query } from '@nestjs/common';
import { TeamService } from './team.service';
import { UpdateTeamDto } from './dto/update-team.dto';
import { ComponentService } from 'src/costing/services/component.service';

@Controller({
  version: '1'
})
export class TeamController {
  constructor(
    private readonly teamSvc: TeamService,
    private readonly compSvc: ComponentService) {}

  @Get()
  async findOne(@Request() req) {
    return this.teamSvc.findById(req.org, req.team.id);
  }

  @Get('/cost')
  async getTeamCost(@Request() req) {
    const { org, team } = req;
    
    // TODO : Add get cost by team
    
    return team;
  }

  @Patch()
  async update(@Request() req, @Param('id') id: number, @Body() updateTeamDto: UpdateTeamDto) {
    const org = req.org;

    return this.teamSvc.update(org, id, updateTeamDto);
  }

  @Delete()
  remove(@Request() req, @Param('id') id: number) {
    const org = req.org;

    return this.teamSvc.remove(org, id);
  }
}

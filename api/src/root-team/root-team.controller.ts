import { Controller, Get, Post, Body, Request } from '@nestjs/common';
import { RootTeamService } from './root-team.service';
import { CreateTeamDto } from 'src/team/dto/create-team.dto';

@Controller({
  version: '1'
})
export class RootTeamController {
  constructor(private readonly rootTeamService: RootTeamService) {}

  @Post()
  create(@Request() req, @Body() createTeam: CreateTeamDto) {
    const org = req.org;

    createTeam.organization = org;

    return this.rootTeamService.create(createTeam);
  }

  @Get()
  async findAll(@Request() req) {
    const org = req.org;

    return this.rootTeamService.findAll(org);
  }
}

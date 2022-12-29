import { Controller, Get, Post, Body, Request, BadRequestException, Logger, InternalServerErrorException } from '@nestjs/common';
import { RootTeamService } from './root-team.service';
import { CreateTeamDto } from 'src/team/dto/create-team.dto';
import { SqlErrorCodes } from 'src/types';

@Controller({
  version: '1'
})
export class RootTeamController {
  private readonly logger = new Logger(RootTeamController.name);

  constructor(
    private readonly rootTeamSvc: RootTeamService
    ) {}

  @Post()
  async create(@Request() req, @Body() createTeam: CreateTeamDto) {
    const org = req.org;

    createTeam.organization = org;

    try {
      return await this.rootTeamSvc.create(createTeam);
    } catch (ex) {
      if (ex.code === SqlErrorCodes.DUP_ENTRY) {
        throw new BadRequestException('team already exists');
      }
      
      this.logger.error({ message: 'could not create team', createTeam, ex });
      throw new InternalServerErrorException('could not create team');
    }
  }

  @Get()
  async findAll(@Request() req) {
    const org = req.org;

    return this.rootTeamSvc.findAll(org);
  }
}

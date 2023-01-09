import { Controller, Get, Post, Body, Request, BadRequestException, Logger, InternalServerErrorException } from '@nestjs/common';
import { RootTeamService } from './root.team.service';
import { CreateTeamDto } from 'src/team/dto/create-team.dto';
import { APIRequest } from 'src/types';
import { TeamSpecDto } from 'src/team/dto/team-spec.dto';
import { TeamService } from 'src/team/team.service';
import { getGithubOrgFromRepoUrl } from 'src/organization/utilities';
import { handleSqlErrors } from 'src/utilities/errorHandler';

@Controller({
  version: '1'
})
export class RootTeamController {
  private readonly logger = new Logger(RootTeamController.name);
  private TeamNameRegex = /^[a-zA-Z]+[a-zA-Z0-9]*(-[a-zA-Z0-9]+)*$/;

  constructor(
    private readonly rootTeamSvc: RootTeamService,
    private readonly teamSvc: TeamService
    ) {}

  @Post()
  async spec(@Request() req: APIRequest, @Body() spec: TeamSpecDto) {
    const { org } = req;

    let team = await this.teamSvc.findByName(org, spec.teamName);

    if (!team) {
      return await this.createTeam(req, {
        name: spec.teamName,
        organization: org,
        repo: spec.configRepo.source,
        repo_path: spec.configRepo.path
      });
    } else {
      return this.teamSvc.update(org, team.id, {
        name: spec.teamName,
        repo: spec.configRepo.source,
        repo_path: spec.configRepo.path
      });
    }
  }

  async createTeam(@Request() req: APIRequest, @Body() createTeam: CreateTeamDto) {
    if (!createTeam.name || !this.TeamNameRegex.test(createTeam.name) || createTeam.name.length > 63) {
      throw new BadRequestException("team name is invalid");
    }

    // validate git repo
    const orgName = getGithubOrgFromRepoUrl(createTeam.repo);
    if (!orgName) {
      throw new BadRequestException('bad github repo url');
    }

    const org = req.org;
    createTeam.organization = org;

    try {
      return await this.rootTeamSvc.create(createTeam);
    } catch (err) {
      handleSqlErrors(err, 'team already exists');
      
      this.logger.error({ message: 'could not create team', createTeam, err });
      throw new InternalServerErrorException('could not create team');
    }
  }

  @Get()
  async findAll(@Request() req) {
    const org = req.org;

    return this.rootTeamSvc.findAll(org);
  }
}

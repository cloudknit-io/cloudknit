import {
  Controller,
  Get,
  Body,
  Patch,
  Delete,
  Request,
  Query,
  Logger,
  Post,
  BadRequestException,
  InternalServerErrorException,
} from '@nestjs/common';
import { TeamService } from './team.service';
import { UpdateTeamDto } from './dto/update-team.dto';
import { APIRequest, OrgApiParam, TeamApiParam } from 'src/types';
import { TeamSpecDto } from './dto/team-spec.dto';
import { CreateTeamDto } from './dto/create-team.dto';
import { getGithubOrgFromRepoUrl } from 'src/organization/utilities';
import { handleSqlErrors } from 'src/utilities/errorHandler';
import { TeamQueryParams } from './team.dto';
import { Team } from 'src/typeorm';

@Controller({
  version: '1',
})
export class TeamController {
  private readonly logger = new Logger(TeamController.name);
  private TeamNameRegex = /^[a-zA-Z]+[a-zA-Z0-9]*(-[a-zA-Z0-9]+)*$/;

  constructor(private readonly teamSvc: TeamService) {}

  @Post()
  @OrgApiParam()
  async spec(@Request() req: APIRequest, @Body() spec: TeamSpecDto) {
    const { org } = req;

    let team = await this.teamSvc.findByName(org, spec.teamName);

    if (!team) {
      return await this.createTeam(req, {
        name: spec.teamName,
        organization: org,
        repo: spec.configRepo.source,
        repo_path: spec.configRepo.path,
      });
    } else {
      return this.teamSvc.update(org, team.id, {
        name: spec.teamName,
        repo: spec.configRepo.source,
        repo_path: spec.configRepo.path,
      });
    }
  }

  async createTeam(
    @Request() req: APIRequest,
    @Body() createTeam: CreateTeamDto
  ) {
    if (
      !createTeam.name ||
      !this.TeamNameRegex.test(createTeam.name) ||
      createTeam.name.length > 63
    ) {
      throw new BadRequestException('team name is invalid');
    }

    // validate git repo
    const orgName = getGithubOrgFromRepoUrl(createTeam.repo);
    if (!orgName) {
      throw new BadRequestException('bad github repo url');
    }

    const org = req.org;
    createTeam.organization = org;

    try {
      return await this.teamSvc.create(createTeam);
    } catch (err) {
      handleSqlErrors(err, 'team already exists');

      this.logger.error({ message: 'could not create team', createTeam, err });
      throw new InternalServerErrorException('could not create team');
    }
  }

  @Get()
  @OrgApiParam()
  async findAll(
    @Request() req: APIRequest,
    @Query() qParams: TeamQueryParams
  ): Promise<Team[]> {
    const org = req.org;
    const withEnv = qParams.withEnvironments.toLowerCase() === 'true';
    const getEnvs = withEnv;

    return this.teamSvc.findAll(org, getEnvs);
  }

  @Get('/:teamId')
  @TeamApiParam()
  async findOne(@Request() req: APIRequest) {
    return this.teamSvc.findById(req.org, req.team.id);
  }

  @Get('/:teamId/cost')
  @TeamApiParam()
  async getTeamCost(@Request() req: APIRequest) {
    const { org, team } = req;

    // TODO : Add get cost by team

    return team;
  }

  @Patch('/:teamId')
  @TeamApiParam()
  async update(
    @Request() req: APIRequest,
    @Body() updateTeamDto: UpdateTeamDto
  ) {
    const { org, team } = req;

    return this.teamSvc.update(org, team.id, updateTeamDto);
  }

  @Delete('/:teamId')
  @TeamApiParam()
  remove(@Request() req: APIRequest) {
    const { org, team } = req;

    return this.teamSvc.remove(org, team.id);
  }
}

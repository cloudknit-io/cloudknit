import {
  BadRequestException, Body, Controller, Delete, Get, InternalServerErrorException, Logger, Patch, Post, Query, Request
} from '@nestjs/common';
import { OnEvent } from '@nestjs/event-emitter';
import { ApiTags } from '@nestjs/swagger';
import { getGithubOrgFromRepoUrl } from 'src/organization/utilities';
import { Team } from 'src/typeorm';
import {
  APIRequest, EnvironmentReconCostUpdateEvent,
  InternalEventType,
  OrgApiParam,
  TeamApiParam
} from 'src/types';
import { handleSqlErrors } from 'src/utilities/errorHandler';
import { CreateTeamDto } from './dto/create-team.dto';
import { TeamSpecDto } from './dto/team-spec.dto';
import { UpdateTeamDto } from './dto/update-team.dto';
import { TeamQueryParams } from './team.dto';
import { TeamService } from './team.service';

@Controller({
  version: '1',
})
@ApiTags('teams')
export class TeamController {
  private readonly logger = new Logger(TeamController.name);
  private TeamNameRegex = /^[a-zA-Z]+[a-zA-Z0-9]*(-[a-zA-Z0-9]+)*$/;

  constructor(private readonly teamSvc: TeamService) {}

  @Post()
  @OrgApiParam()
  async spec(@Request() req: APIRequest, @Body() spec: TeamSpecDto) {
    const { org } = req;
    console.log(spec);
    let team = await this.teamSvc.findByName(org, spec.teamName);
    if (!team) {
      return await this.createTeam(req, {
        name: spec.teamName,
        organization: org,
        repo: spec.configRepo.source,
        repo_path: spec.configRepo.path,
        teardownProtection: spec.teardownProtection
      });
    } else {
      return this.teamSvc.update(org, team.id, {
        name: spec.teamName,
        repo: spec.configRepo.source,
        repo_path: spec.configRepo.path,
        teardownProtection: Boolean(spec.teardownProtection)
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
    const withComps = qParams.withComponents.toLowerCase() === 'true';
    const getEnvs = withEnv;

    return this.teamSvc.findAll(org, getEnvs, withComps);
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

  @OnEvent(InternalEventType.EnvironmentReconCostUpdate, { async: true })
  async environmentCostUpdateListener(evt: EnvironmentReconCostUpdateEvent) {
    const envRecon = evt.payload;

    await this.teamSvc.updateCost(envRecon.organization, envRecon.teamId);
  }
}

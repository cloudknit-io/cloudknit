import {
  Controller,
  Get,
  Body,
  Patch,
  Delete,
  Request,
  Post,
  Logger,
  InternalServerErrorException,
} from '@nestjs/common';
import { EnvironmentService } from './environment.service';
import { UpdateEnvironmentDto } from './dto/update-environment.dto';
import { APIRequest, EnvironmentApiParam, TeamApiParam } from 'src/types';
import { handleSqlErrors } from 'src/utilities/errorHandler';
import { ComponentService } from 'src/component/component.service';
import { EnvSpecComponentDto, EnvSpecDto } from './dto/env-spec.dto';
import { Component, Environment, Organization, Team } from 'src/typeorm';
import { CreateEnvironmentDto } from './dto/create-environment.dto';
import { ReconciliationService } from 'src/reconciliation/reconciliation.service';
import { ApiTags } from '@nestjs/swagger';

@Controller({
  version: '1',
})
@ApiTags('environments')
export class EnvironmentController {
  private readonly logger = new Logger(EnvironmentController.name);

  constructor(
    private readonly envSvc: EnvironmentService,
    private readonly reconSvc: ReconciliationService,
    private readonly compSvc: ComponentService
  ) {}

  @Post()
  @TeamApiParam()
  async saveOrUpdate(@Request() req: APIRequest, @Body() body: EnvSpecDto) {
    const { org, team } = req;

    let env = await this.envSvc.findByName(org, team, body.envName);

    if (!env) {
      return this.createEnv(org, team, {
        name: body.envName,
        dag: body.components,
      });
    }

    const currentComps: Component[] =
      await this.compSvc.getAllForEnvironmentById(org, env);
    const incoming: EnvSpecComponentDto[] = body.components;

    const newComponents: EnvSpecComponentDto[] = incoming.filter((inc) => {
      return !currentComps.find((comp) => comp.name === inc.name);
    });
    const missingComponents: Component[] = [];
    const existingComponents: EnvSpecComponentDto[] = [];

    for (const comp of currentComps) {
      const found = incoming.find((i) => comp.name === i.name);

      if (!found) {
        missingComponents.push(comp);
        continue;
      }

      existingComponents.push(found);
    }

    const dag: EnvSpecComponentDto[] = []
      .concat(existingComponents)
      .concat(newComponents);

    env = await this.envSvc.updateById(org, env.id, {
      dag,
      name: env.name,
      duration: env.duration,
      isDeleted: env.isDeleted,
    });

    // create new components
    await this.batchCreateComponents(org, env, newComponents);

    // delete missing
    await this.batchDeleteComponents(org, env, missingComponents);

    return env;
  }

  @Get()
  @TeamApiParam()
  async findAll(@Request() req: APIRequest): Promise<Environment[]> {
    const { org, team } = req;

    return this.envSvc.findAll(org, team);
  }

  async createEnv(
    org: Organization,
    team: Team,
    createEnv: CreateEnvironmentDto
  ): Promise<Environment> {
    let env: Environment;

    try {
      env = await this.envSvc.create(org, team, createEnv);
      this.logger.log({ message: `created new environment`, env });
    } catch (err) {
      handleSqlErrors(err, 'environment already exists');

      this.logger.error({
        message: 'could not create environment',
        createEnv,
        err,
      });
      throw new InternalServerErrorException('could not create environment');
    }

    if (!createEnv.dag || createEnv.dag.length == 0) {
      return env;
    }

    // create all new components
    await this.batchCreateComponents(org, env, createEnv.dag);

    return env;
  }

  async batchCreateComponents(
    org: Organization,
    env: Environment,
    comps: EnvSpecComponentDto[]
  ) {
    if (!comps || comps.length === 0) {
      return;
    }

    try {
      const res = await this.compSvc.batchCreate(
        org,
        env,
        comps.map((comp) => comp.name)
      );
      this.logger.log({
        message: `created ${res.identifiers.length} new components`,
        env,
      });
    } catch (err) {
      handleSqlErrors(err, 'component already exists');

      this.logger.error({
        message:
          'could not batch create components during environment creation',
        err,
      });
      throw new InternalServerErrorException('could not create components');
    }
  }

  async batchDeleteComponents(
    org: Organization,
    env: Environment,
    comps: Component[]
  ) {
    if (!comps || comps.length === 0) {
      return;
    }

    try {
      const res = await this.compSvc.batchDelete(org, env, comps);
      this.logger.log({ message: `deleted ${res.affected} components`, env });
    } catch (err) {
      handleSqlErrors(err);

      this.logger.error({
        message:
          'could not batch delete components during environment spec reconciliation',
        err,
      });
      throw new InternalServerErrorException('could not delete components');
    }
  }

  @Get('/:environmentId')
  @EnvironmentApiParam()
  async findOne(@Request() req: APIRequest) {
    const { org, team, env } = req;

    return this.envSvc.findById(org, env.id);
  }

  @Get('/:environmentId/dag')
  @EnvironmentApiParam()
  async getDag(@Request() req: APIRequest) {
    const { env } = req;

    return env.dag;
  }

  @Patch('/:environmentId')
  @EnvironmentApiParam()
  async update(
    @Request() req: APIRequest,
    @Body() updateEnvDto: UpdateEnvironmentDto
  ) {
    const { org, env } = req;

    return this.envSvc.updateById(org, env.id, updateEnvDto);
  }

  @Delete('/:environmentId')
  @EnvironmentApiParam()
  remove(@Request() req: APIRequest) {
    const { org, env } = req;

    return this.envSvc.remove(org, env.id);
  }

  @Get('/:environmentId/audit')
  @EnvironmentApiParam()
  async getAudits(@Request() req: APIRequest) {
    const { org, env } = req;

    return this.reconSvc.getEnvironmentAuditList(org, env);
  }
}

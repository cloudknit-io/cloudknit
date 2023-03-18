import {
  Body, Controller, Delete, Get, InternalServerErrorException, Logger, Patch, Post, Request
} from '@nestjs/common';
import { OnEvent } from '@nestjs/event-emitter';
import { ApiTags } from '@nestjs/swagger';
import { ComponentService } from 'src/component/component.service';
import { ReconciliationService } from 'src/reconciliation/reconciliation.service';
import { Component, Environment, Organization, Team } from 'src/typeorm';
import {
  APIRequest,
  ComponentCostUpdateEvent,
  EnvironmentApiParam,
  EnvironmentReconCostUpdateEvent,
  EnvironmentReconEnvUpdateEvent,
  InternalEventType,
  TeamApiParam
} from 'src/types';
import { handleSqlErrors } from 'src/utilities/errorHandler';
import { CreateEnvironmentDto } from './dto/create-environment.dto';
import { EnvSpecComponentDto, EnvSpecDto } from './dto/env-spec.dto';
import { UpdateEnvironmentDto } from './dto/update-environment.dto';
import { EnvironmentService } from './environment.service';

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
      isDeleted: env.isDeleted,
    });

    // create new components
    await this.batchCreateComponents(org, env, newComponents);

    // delete missing
    await this.batchDeleteComponents(org, env, missingComponents);

    return env;
  }

  @Patch('/:environmentId/reconcile')
  @EnvironmentApiParam()
  async Update(@Request() req: APIRequest) {
    const { org, team, env } = req;

    const envRecon = await this.createEnvRecon(org, team, env, {
      dag: env.dag,
      name: env.name,
    });

    envRecon.environment = null;

    const newEnv = await this.envSvc.updateById(org, env.id, {
      latestEnvRecon: envRecon,
    });

    return newEnv;
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

  async createEnvRecon(
    org: Organization,
    team: Team,
    env: Environment,
    createEnv: CreateEnvironmentDto
  ) {
    return this.reconSvc.createEnvRecon(org, team, env, {
      components: createEnv.dag,
      name: env.name,
      startDateTime: new Date().toISOString(),
      teamName: team.name,
    });
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

    return env.latestEnvRecon.dag;
  }

  @Get('/:environmentId/audit/latest')
  @EnvironmentApiParam()
  async getLatestRecon(@Request() req: APIRequest) {
    const { env } = req;

    return { reconcileId: env.latestEnvRecon.reconcileId };
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

  @OnEvent(InternalEventType.ComponentCostUpdate, { async: true })
  async compCostUpdateListener(evt: ComponentCostUpdateEvent) {
    const comp = evt.payload;
    let env = comp.environment;

    if (!env) {
      env = await this.envSvc.findById(
        comp.organization,
        comp.envId,
        false,
        true
      );
    }

    await this.reconSvc.updateCost(env);
  }

  @OnEvent(InternalEventType.EnvironmentReconCostUpdate, { async: true })
  async envReconCostUpdateListener(evt: EnvironmentReconCostUpdateEvent) {
    const envRecon = evt.payload;

    const env = await this.envSvc.findById(
      envRecon.organization,
      envRecon.envId
    );

    await this.envSvc.mergeAndSaveEnv(envRecon.organization, env, {
      lastReconcileDatetime: new Date().toISOString(),
    });
  }

  @OnEvent(InternalEventType.EnvironmentReconEnvUpdate, { async: true })
  async envReconEnvUpdateListener(evt: EnvironmentReconEnvUpdateEvent) {
    const envRecon = evt.payload;

    const env = await this.envSvc.findById(
      envRecon.organization,
      envRecon.envId
    );

    await this.envSvc.mergeAndSaveEnv(envRecon.organization, env, {
      lastReconcileDatetime: new Date().toISOString(),
    });
  }
}

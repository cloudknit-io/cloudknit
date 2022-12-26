import { BadRequestException, Injectable, InternalServerErrorException, Logger, NotFoundException } from "@nestjs/common";
import { InjectDataSource, InjectRepository } from "@nestjs/typeorm";
import { Subject } from "rxjs";
import { get } from "src/config";
import { Mapper } from "src/costing/utilities/mapper";
import { S3Handler } from "src/utilities/s3Handler";
import { Organization } from "src/typeorm";
import { ComponentReconcile } from "src/typeorm/reconciliation/component-reconcile.entity";
import { Component } from "src/typeorm/component.entity";
import { EnvironmentReconcile } from "src/typeorm/reconciliation/environment-reconcile.entity";
import { Environment } from "src/typeorm/reconciliation/environment.entity";
import { In, IsNull, Like, Not } from "typeorm";
import { Repository } from "typeorm/repository/Repository";
import { ComponentDto } from "./dtos/component.dto";
import { ApprovedByDto, ComponentAudit } from "./dtos/componentAudit.dto";
import { EnvironmentDto } from "./dtos/environment.dto";
import { EnvironmentAudit } from "./dtos/environmentAudit.dto";
import { EvnironmentReconcileDto } from "./dtos/reconcile.Dto";
import { EnvironmentService } from "./environment.service";
import { SSEService } from "./sse.service";

@Injectable()
export class ReconciliationService {
  private readonly s3h = S3Handler.instance();
  private readonly config = get();
  private readonly ckEnvironment = this.config.environment;
  private readonly logger = new Logger(ReconciliationService.name);

  constructor(
    @InjectRepository(EnvironmentReconcile)
    private readonly environmentReconcileRepository: Repository<EnvironmentReconcile>,
    @InjectRepository(ComponentReconcile)
    private readonly componentReconcileRepository: Repository<ComponentReconcile>,
    @InjectRepository(Component)
    private readonly componentRepository: Repository<Component>,
    @InjectRepository(Environment)
    private readonly environmentRepository: Repository<Environment>,
    private readonly envSvc: EnvironmentService,
    private readonly sseSvc: SSEService
  ) { }

  async putComponent(org: Organization, component: ComponentDto, env: Environment, teamName: string) {
    if (typeof env === 'undefined') {
      env = await this.envSvc.getEnvironment(org, component.environmentName, teamName);
    }

    const componentName = `${teamName}-${component.environmentName}-${component.name}`;

    const existing = await this.componentRepository
      .createQueryBuilder()
      .where('id = :id and component_name = :name and environmentId = :envId', {
        id: componentName,
        name: component.name,
        envId: env.id
      })
      .getOne();

    this.logger.log({ message: "PutComponent", existing, id: componentName, teamName, component, env});

    if (!existing) {
      const entry = await this.componentRepository.save({
        id: componentName,
        componentName: component.name,
        duration: component.duration,
        teamName: teamName,
        environment: env,
        status: component.status,
        organization: org
      });
      return entry;
    }

    existing.duration = component.duration;
    const entry = await this.componentRepository.save(existing);
    
    return entry;
  }

  async getComponent(org: Organization, compName: string, envName: string, teamName: string) {
    const comp = await this.componentRepository.findOne({
      where: {
        componentName: compName,
        organization: {
          id: org.id
        },
        environment: {
          name : envName,
          teamName: teamName
        }
      }
    });

    if (!comp) {
      this.logger.error(`Bad component query [${compName}] for org [${org.id} / ${org.name}]`);
      throw new NotFoundException('could not find component');
    }
    return comp;
  }

  async saveOrUpdateEnvironment(org: Organization, runData: EvnironmentReconcileDto) {
    const reconcileId = Number.isNaN(parseInt(runData.reconcileId))
      ? null
      : parseInt(runData.reconcileId);

    let savedEntry: EnvironmentReconcile = null;
    let env = await this.envSvc.getEnvironment(org, runData.name, runData.teamName);

    if (!env) {
      try {
        env = await this.envSvc.putEnvironment(org, {name: runData.name, teamName: runData.teamName, duration: -1});

        this.logger.log({message: 'created environment', runData });
      } catch (e) {
        this.logger.error({message: `could not create environment ${runData.name}`, error: e.message });
        throw new InternalServerErrorException('could not create environment');
      }
    }

    this.logger.log({ message: 'saveOrUpdateEnvironment', reconcileId, environment: env, runData });

    if (reconcileId) {
      const existingEntry = await this.environmentReconcileRepository
        .createQueryBuilder()
        .where('reconcile_id = :reconcileId and environmentId = :envId',
          {
            reconcileId,
            envId: env.id
          }
        )
        .getOne();

      if (!existingEntry) {
        this.logger.error({message: `could not find environment reconcile entry`, reconcileId, environment: env})
        throw new BadRequestException(`could not find environment reconcile entry`);
      }
      
      existingEntry.end_date_time = runData.endDateTime ? runData.endDateTime : "";
      existingEntry.status = runData.status;
        
      this.logger.log({ message: 'updating existing environment reconcile entry', existingEntry });
      
      savedEntry = await this.environmentReconcileRepository.save(
        existingEntry
      );

      const ed = new Date(savedEntry.end_date_time).getTime();
      const sd = new Date(savedEntry.start_date_time).getTime();
      const duration = ed - sd;

      await this.envSvc.putEnvironment(org, {
        name: savedEntry.name,
        teamName: runData.teamName,
        duration,
      });
    } else {
      const entry: EnvironmentReconcile = {
        reconcile_id: null,
        environment: env,
        name: runData.name,
        start_date_time: runData.startDateTime,
        team_name: runData.teamName,
        status: runData.status,
        end_date_time: runData.endDateTime,
        organization: org
      };

      // QUESTION : This queries all previously unfinished (status == null) EnvironmentReconcile's 
      // by environment name and sets them to "skipped".
      // What it should do is query Reconcile table to get the most recent run for the Environment
      // Then query the EnvironmentReconcile table by ReconcileId and set those to "skipped"
      //
      // Querying by `name` is the wrong way to do this.
      try {
        await this.updateSkippedWorkflows<EnvironmentReconcile>(
          entry.name,
          this.environmentReconcileRepository
        );
      } catch (err) {
        this.logger.error('could not update skipped workflows', err);
        throw new InternalServerErrorException();
      }

      this.logger.log({ message: 'creating new environmentReconcileEntry', entry });
      savedEntry = await this.environmentReconcileRepository.save(entry);

      await this.envSvc.putEnvironment(org, {
        name: savedEntry.name,
        teamName: runData.teamName,
        duration: -1,
      });
    }
    savedEntry.organization = org;
    this.sseSvc.sendEnvironmentReconcile(savedEntry);

    return savedEntry.reconcile_id;
  }

  async saveOrUpdateComponent(org: Organization, envReconcile: EvnironmentReconcileDto) {
    const reconcileId = Number.isNaN(parseInt(envReconcile.reconcileId))
      ? null
      : parseInt(envReconcile.reconcileId);

    if (!reconcileId) {
      this.logger.error({ message: 'No reconcileId found when trying to save or update a component', org, envReconcile });
      throw new BadRequestException("reconcileId is mandatory to save or update component")
    }

    this.logger.log({message: 'save or update component', reconcileId, envReconcile});

    const env = await this.envSvc.getEnvironment(org, envReconcile.name, envReconcile.teamName);

    if (!env) {
      throw new BadRequestException(`could not find environment ${envReconcile.name}`);
    }

    this.logger.log({message: 'found environment', reconcileId, env});

    const envRecEntry = await this.environmentReconcileRepository
      .createQueryBuilder()
      .where('reconcile_id = :reconcileId and environmentId = :envId',
        {
          reconcileId,
          envId: env.id
        }
      )
      .getOne();

    if (!envRecEntry) {
      throw new BadRequestException(`could not find environmentReconcileEntry for environment ${envReconcile.name}`);
    }

    this.logger.log({message: 'found environment reconcile entry', reconcileId, envRecEntry});

    let componentEntry: ComponentReconcile = Mapper.mapToComponentReconcile(
      org,
      envRecEntry,
      envReconcile.componentReconciles
    )[0];

    this.logger.log({message: 'created component entry', reconcileId, env, componentEntry});

    if (!componentEntry.reconcile_id) {
      await this.updateSkippedWorkflows<ComponentReconcile>(
        componentEntry.name,
        this.componentReconcileRepository
      );
    }

    let duration = -1;
    if (componentEntry.reconcile_id) {
      const existingEntry = await this.componentReconcileRepository.findOne(
        {
          where: {
            reconcile_id: componentEntry.reconcile_id
          }
        }
      );
      existingEntry.end_date_time = componentEntry.end_date_time;
      existingEntry.status = componentEntry.status;
      componentEntry = existingEntry;
      const ed = new Date(existingEntry.end_date_time).getTime();
      const sd = new Date(existingEntry.start_date_time).getTime();
      duration = ed - sd;
    }

    if (envReconcile.componentReconciles.length == 0) {
      throw new BadRequestException("no component reconciles");
    }

    const comp = envReconcile.componentReconciles[0];

    await this.putComponent(org, {
      name: comp.name,
      duration,
      teamName: comp.teamName,
      environmentName: envRecEntry.name,
      status: comp.status,
    }, env, envReconcile.teamName);

    const entry = await this.componentReconcileRepository.save(componentEntry);
    entry.organization = org;
    entry.environmentReconcile = new EnvironmentReconcile();
    entry.environmentReconcile.name = envReconcile.name;
    entry.environmentReconcile.team_name = envReconcile.teamName;
    entry.status = componentEntry.status;
    this.sseSvc.sendComponentReconcile(entry);

    this.logger.log({message: 'created component reconcile entry', reconcileId, compReconcileEntry: entry});

    return entry.reconcile_id;
  }

  async updateSkippedWorkflows<T>(name: string, repo: Repository<ComponentReconcile|EnvironmentReconcile>) {
    const entries = await repo.find({
      where: {
        name: name,
        end_date_time: IsNull(),
      },
    });

    if (entries.length > 0) {
      const newEntries = entries.map((entry) => ({
        ...entry,
        status: "skipped_reconcile",
      })) as any;

      this.logger.log({message: 'updating skipped workflows', newEntries});

      await repo.save(newEntries);
    }
  }

  async getComponentAuditList(org: Organization, id: string, envName: string, teamName: string): Promise<ComponentAudit[]> {
    const envAuditList = await this.getEnvironmentAuditList(org, envName, teamName);
    const reconcileIds = envAuditList.map(e => e.reconcileId);
    const components = await this.componentReconcileRepository.find({
      where: {
        name: id,
        organization: {
          id: org.id
        },
        environmentReconcile: {
          reconcile_id: In(reconcileIds)
        }
      }
    });

    return Mapper.getComponentAuditList(components);
  }

  async getEnvironmentAuditList(org: Organization, envName: string, teamName: string): Promise<EnvironmentAudit[]> {
    const env = await this.envSvc.getEnvironment(org, envName, teamName);

    if (!env) {
      this.logger.error({ message: `Could not find environment with name ${envName}`, org })
      throw new BadRequestException(`Could not find environment`);
    }

    const environments = await this.environmentReconcileRepository
      .createQueryBuilder()
      .where('name = :name and organizationId = :orgId', {
        name: envName,
        orgId: org.id
      })
      .getMany();

    return Mapper.getEnvironmentAuditList(environments);
  }

  async getLogs(
    org: Organization,
    team: string,
    environment: string,
    component: string,
    id: number
  ) {
    const prefix = `${team}/${environment}/${component}/${id}/`;
    const bucket = `zlifecycle-${this.ckEnvironment}-tfplan-${org.name}`;
    
    try {
      const objects = await this.s3h.getObjects(
        bucket,
        prefix
      );

      return objects.map((o) => ({
        key: o.key,
        body: o.data.Body.toString(),
      }));
    } catch (err) {
      this.logger.error({ message: 'error getting S3 terraform logs', prefix, bucket }, err);
      return err;
    }
  }

  async getApplyLogs(
    org: Organization,
    team: string,
    environment: string,
    component: string,
    id: number,
    latest: boolean
  ) {
    const logs =
      latest === true
        ? await this.getLatestLogs(org, team, environment, component)
        : await this.getLogs(org, team, environment, component, id);

    if (Array.isArray(logs)) {
      return logs.filter((e) => e.key.includes("apply_output"));
    }

    return logs;
  }

  async getPlanLogs(
    org: Organization,
    team: string,
    environment: string,
    component: string,
    id: number,
    latest: boolean
  ) {
    const logs =
      latest === true
        ? await this.getLatestLogs(org, team, environment, component)
        : await this.getLogs(org, team, environment, component, id);

    if (Array.isArray(logs)) {
      return logs.filter((e) => e.key.includes("plan_output"));
    }

    return logs;
  }

  async getLatestLogs(
    org: Organization,
    team: string,
    environment: string,
    component: string
  ) {
    const latestAudit = await this.getLatestCompReconcile(org, component, environment, team);

    if (!latestAudit) {
      this.logger.error({ message: 'could not find latestAuditId to get latest logs', component, environment, team });
      throw new NotFoundException(`could not find audit logs for ${team}-${environment}-${component}`);
    }

    return await this.getLogs(
      org,
      team,
      environment,
      component,
      latestAudit.reconcile_id
    );
  }

  async patchApprovedBy(org: Organization, body: ApprovedByDto) {
    const latestAudit = await this.getLatestCompReconcile(org, body.compName, body.envName, body.teamName);

    if (!latestAudit) {
      throw new NotFoundException(`could not find latest audit for component ${body.compName} in env ${body.envName} for team ${body.teamName}`);
    }

    latestAudit.approved_by = body.email;

    return await this.componentReconcileRepository.save(latestAudit);
  }

  async getApprovedBy(org: Organization, body: ApprovedByDto) {
    if (body.rid < 0) {
      const latestAudit = await this.getLatestCompReconcile(org, body.compName, body.envName, body.teamName);

      if (!latestAudit) {
        throw new NotFoundException(`could not find latest audit for component ${body.compName} in env ${body.envName} for team ${body.teamName}`);
      }

      return latestAudit;
    }

    return await this.componentReconcileRepository.findOne({
      where : {
        reconcile_id : body.rid
      }
    });
  }

  async getStateFile(
    org: Organization,
    team: string,
    environment: string,
    component: string
  ) {
    const prefix = `${team}/${environment}/${component}/terraform.tfstate`;
    const resp = await this.s3h.getObject(
      `zlifecycle-${this.ckEnvironment}-tfstate-${org.name}`,
      prefix
    );

    return {
      ...resp,
      data: (resp.data?.Body || "").toString() || "",
    };
  }

  private async getLatestCompReconcile(org: Organization, compName: string, envName: string, teamName: string): Promise<ComponentReconcile> {
    try {
      // const latestAuditId = this.componentReconcileRepository
      // .createQueryBuilder()
      // .where('name = :name and organizationId = :orgId', {
      //   name: componentId,
      //   orgId: org.id,
      //   status: Not(Like("skipped%")),
      // })
      // .orderBy('start_date_time', 'DESC')
      // .getOne();

      const envRecon = await this.environmentReconcileRepository.findOne({
        where: {
          name: envName,
          team_name: teamName,
          organization: {
            id : org.id
          }
        },
        order: {
          start_date_time: -1,
        }
      });
  
      const latestAuditId = await this.componentReconcileRepository.findOne({
        where: {
          name: compName,
          status: Not(Like("skipped%")),
          organization: {
            id : org.id
          },
          environmentReconcile: {
            reconcile_id: envRecon.reconcile_id
          }
        },
        order: {
          start_date_time: -1,
        }
      });

      this.logger.debug({message: `getting latest component reconcile id`, compName, envName, teamName, envRecon});
  
      return latestAuditId;
    } catch (err) {
      this.logger.error('could not get latestAuditId', err);
      throw new InternalServerErrorException('could not get latest audit id');
    }
  }
}

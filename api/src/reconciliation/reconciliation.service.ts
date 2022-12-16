import { BadRequestException, Injectable, InternalServerErrorException, Logger, NotFoundException } from "@nestjs/common";
import { InjectRepository } from "@nestjs/typeorm";
import { Subject } from "rxjs";
import { get } from "src/config";
import { Mapper } from "src/costing/utilities/mapper";
import { S3Handler } from "src/utilities/s3Handler";
import { Organization } from "src/typeorm";
import { ComponentReconcile } from "src/typeorm/reconciliation/component-reconcile.entity";
import { Component } from "src/typeorm/component.entity";
import { EnvironmentReconcile } from "src/typeorm/reconciliation/environment-reconcile.entity";
import { Environment } from "src/typeorm/reconciliation/environment.entity";
import { IsNull, Like, Not } from "typeorm";
import { Repository } from "typeorm/repository/Repository";
import { ComponentDto } from "./dtos/component.dto";
import { ComponentAudit } from "./dtos/componentAudit.dto";
import { EnvironmentDto } from "./dtos/environment.dto";
import { EnvironmentAudit } from "./dtos/environmentAudit.dto";
import { EvnironmentReconcileDto } from "./dtos/reconcile.Dto";

@Injectable()
export class ReconciliationService {
  readonly notifyStream: Subject<{}> = new Subject<{}>();
  readonly applicationStream: Subject<any> = new Subject<any>();
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
    private readonly environmentRepository: Repository<Environment>
  ) {
    setInterval(() => {
      this.notifyStream.next({});
      this.applicationStream.next({});
    }, 20000);
  }

  async putEnvironment(org: Organization, environment: EnvironmentDto) {
    const existing = await this.getEnvironment(org, environment.name);

    if (!existing) {
      let env = new Environment();
      env.name = environment.name;
      env.organization = org;
      env.duration = environment.duration;

      return await this.environmentRepository.save(env);
    }

    existing.duration = environment.duration;
    const entry = await this.environmentRepository.save(existing);
    
    this.notifyApplications(entry.name);

    return entry;
  }

  async putComponent(org: Organization, component: ComponentDto, env: Environment, teamName: string) {
    if (typeof env === 'undefined') {
      env = await this.getEnvironment(org, component.environmentName);
    }

    const id = `${teamName}-${env.name}-${component.componentName}`;

    const existing = await this.componentRepository
      .createQueryBuilder()
      .where('id = :id and component_name = :name and environmentId = :envId', {
        id,
        name: component.componentName,
        envId: env.id
      })
      .getOne();

    this.logger.log({ message: "PutComponent", existing, teamName});

    if (!existing) {
      const entry = await this.componentRepository.save({
        id,
        componentName: component.componentName,
        duration: component.duration,
        teamName: teamName,
        environment: env
      });
      this.notifyApplications(component.environmentName);
      return entry;
    }

    existing.duration = component.duration;
    const entry = await this.componentRepository.save(existing);
    this.notifyApplications(component.environmentName);
    
    return entry;
  }

  async getComponent(org: Organization, id: string) {
    const comp = await this.componentRepository.findOne({
      where: {
        componentName: id,
      },
      relations: {
        environment: true
      }
    });

    if (!comp) {
      this.logger.error(`Bad component query [${id}] for org [${org.id} / ${org.name}]`);
      throw new NotFoundException('could not find component');
    }

    const env = await this.environmentRepository.findOne({
      where: {
        id: comp.environment.id
      },
      relations: {
        organization: true
      }
    });

    if (env.organization.id === org.id) {
      return comp;
    }

    this.logger.error(`Could not find component [${id}] for org [${org.id} / ${org.name}]`);
    throw new NotFoundException('could not find component');
  }

  async getEnvironment(org: Organization, id: string) {
    return await this.environmentRepository
      .createQueryBuilder()
      .where('organizationId = :orgId and name = :name', {
        orgId: org.id,
        name: id
      })
      .getOne();
  }

  async saveOrUpdateEnvironment(org: Organization, runData: EvnironmentReconcileDto) {
    const reconcileId = Number.isNaN(parseInt(runData.reconcileId))
      ? null
      : parseInt(runData.reconcileId);

    let savedEntry: EnvironmentReconcile = null;
    let env = await this.getEnvironment(org, runData.name);

    if (!env) {      
      try {
        env = await this.putEnvironment(org, {name: runData.name, duration: -1});

        this.logger.log(`created environment ${runData.name} for org ${org.name}`);
      } catch (e) {
        this.logger.error(`could not create environment ${runData.name} for org ${org.name}`, e.message, org.id);
        throw new InternalServerErrorException('could not create environment');
      }
    }

    this.logger.log(`saveOrUpdateEnvironment reconcileId: ${reconcileId}, environment: ${JSON.stringify(env)}, runData: ${JSON.stringify(runData)}`);

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
        
      this.logger.log(`updating existing environment reconcile entry ${JSON.stringify(existingEntry)}`);
      
      savedEntry = await this.environmentReconcileRepository.save(
        existingEntry
      );

      const ed = new Date(savedEntry.end_date_time).getTime();
      const sd = new Date(savedEntry.start_date_time).getTime();
      const duration = ed - sd;

      await this.putEnvironment(org, {
        name: savedEntry.name,
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

      this.logger.log(`creating new environmentReconcileEntry ${JSON.stringify(entry)}`);
      savedEntry = await this.environmentReconcileRepository.save(entry);

      await this.putEnvironment(org, {
        name: savedEntry.name,
        duration: -1,
      });
    }

    this.notifyStream.next(savedEntry);

    return savedEntry.reconcile_id;
  }

  async saveOrUpdateComponent(org: Organization, runData: EvnironmentReconcileDto) {
    const reconcileId = Number.isNaN(parseInt(runData.reconcileId))
      ? null
      : parseInt(runData.reconcileId);

    if (!reconcileId) {
      this.logger.error(`No reconcileId found when trying to save or update a component. org [${org.id} / ${org.name}] env [${runData.name}]`);
      throw new BadRequestException("reconcileId is mandatory to save or update component")
    }

    this.logger.log(`reconcileId: ${reconcileId}, component: ${runData.name} save or update ${org.name}`);

    const env = await this.getEnvironment(org, runData.name);

    if (!env) {
      throw new BadRequestException(`could not find environment ${runData.name}`);
    }

    this.logger.log(`reconcileId ${reconcileId}, component: ${runData.name} found environment ${JSON.stringify(env)} ${org.name}`);

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
      throw new BadRequestException(`could not find environmentReconcileEntry for environment ${runData.name}`);
    }

    this.logger.log(`reconcileId ${reconcileId}, envRecEntry: ${JSON.stringify(envRecEntry)} for org ${org.name}`);

    let componentEntry: ComponentReconcile = Mapper.mapToComponentReconcile(
      org,
      envRecEntry,
      runData.componentReconciles
    )[0];

    this.logger.log(`reconcileId ${reconcileId}, componentEntry ${JSON.stringify(componentEntry)} for org ${org.name}`);

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

    await this.putComponent(org, {
      componentName: componentEntry.name,
      duration,
      environmentName: envRecEntry.name,
    }, env, runData.teamName);

    const entry = await this.componentReconcileRepository.save(componentEntry);
    this.notifyStream.next(entry);

    this.logger.log(`reconcileId ${reconcileId}, entry: ${JSON.stringify(entry)}`);

    return entry.reconcile_id;
  }

  async updateSkippedWorkflows<T>(name: string, repo: Repository<ComponentReconcile|EnvironmentReconcile>) {
    const entries = await repo.find({
      where: {
        name: name,
        end_date_time: IsNull(),
      },
    });

    this.logger.log(`updateSkippedWorkflows ${name} entries ${JSON.stringify(entries)}`);

    if (entries.length > 0) {
      const newEntries = entries.map((entry) => ({
        ...entry,
        status: "skipped_reconcile",
      })) as any;

      this.logger.log(`updateSkippedWorkflows ${name} newEntries ${JSON.stringify(newEntries)}`);

      await repo.save(newEntries);
    }
  }

  async getComponentAuditList(org: Organization, id: string): Promise<ComponentAudit[]> {
    const components = await this.componentReconcileRepository
      .createQueryBuilder()
      .where('name = :name and organizationId = :orgId', {
        name: id,
        orgId: org.id
      }).getMany();

    return Mapper.getComponentAuditList(components);
  }

  async getEnvironmentAuditList(org: Organization, id: string): Promise<EnvironmentAudit[]> {
    const env = await this.getEnvironment(org, id);

    if (!env) {
      this.logger.error(`Could not find environment with name [${id}] for org [${org.id} / ${org.name}]`)
      throw new BadRequestException(`Could not find environment`);
    }

    const environments = await this.environmentReconcileRepository
      .createQueryBuilder()
      .where('name = :name and organizationId = :orgId', {
        name: id,
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
    try {
      const prefix = `${team}/${environment}/${component}/${id}/`;
      const bucket = `zlifecycle-${this.ckEnvironment}-tfplan-${org.name}`;

      this.logger.debug(`getLogs prefix ${prefix}`);
      this.logger.debug(`getLogs bucket ${bucket}`);

      const objects = await this.s3h.getObjects(
        bucket,
        prefix
      );

      return objects.map((o) => ({
        key: o.key,
        body: o.data.Body.toString(),
      }));
    } catch (err) {
      this.logger.error('error getting S3 plan/apply logs', err);
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
    const latestAudit = await this.getLatestAudit(org,
      `${team}-${environment}-${component}`
    );

    if (!latestAudit) {
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

  async patchApprovedBy(org: Organization, email: string, componentId: string) {
    const latestAudit = await this.getLatestAudit(org, componentId);

    if (!latestAudit) {
      throw new NotFoundException(`could not find latest audit for component ${componentId}`);
    }

    latestAudit.approved_by = email;

    return await this.componentReconcileRepository.save(latestAudit);
  }

  async getApprovedBy(org: Organization, id: string, rid: string) {
    if (rid === "-1") {
      const latestAudit = await this.getLatestAudit(org, id);

      if (!latestAudit) {
        throw new NotFoundException(`could not find latest audit for component ${id}`);
      }

      return latestAudit;
    }

    return await this.componentReconcileRepository.findOne({
      where : {
        reconcile_id : parseInt(rid)
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

  private async notifyApplications(environmentName: string) {
    const apps = await this.environmentRepository.find({
      where: {
        name: environmentName,
      },
    });
    this.applicationStream.next(apps);
  }

  private async getLatestAudit(org: Organization, componentId): Promise<ComponentReconcile> {
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
  
      const latestAuditId = await this.componentReconcileRepository.find({
        where: {
          name: componentId,
          status: Not(Like("skipped%")),
          organization: {
            id : org.id
          }
        },
        order: {
          start_date_time: -1,
        },
        take: 1,
      });

      this.logger.debug(`latestAuditId ${JSON.stringify(latestAuditId)} - component: ${componentId}`);
  
      return latestAuditId.length > 0 ? latestAuditId[0] : null;
    } catch (err) {
      this.logger.error('could not get latestAuditId', err);
    }
  }
}

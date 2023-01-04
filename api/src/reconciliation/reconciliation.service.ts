import { BadRequestException, Injectable, Logger } from "@nestjs/common";
import { InjectRepository } from "@nestjs/typeorm";
import { get } from "src/config";
import { Mapper } from "src/costing/utilities/mapper";
import { S3Handler } from "src/utilities/s3Handler";
import { Organization, Team } from "src/typeorm";
import { ComponentReconcile } from "src/typeorm/component-reconcile.entity";
import { Component } from "src/typeorm/component.entity";
import { EnvironmentReconcile } from "src/typeorm/environment-reconcile.entity";
import { Environment } from "src/typeorm/environment.entity";
import { Equal, In, IsNull, Like, Not } from "typeorm";
import { Repository } from "typeorm/repository/Repository";
import { ApprovedByDto, ComponentAudit } from "./dtos/componentAudit.dto";
import { EnvironmentAudit } from "./dtos/environmentAudit.dto";
import { EnvironmentService } from "src/environment/environment.service";
import { TeamService } from "src/team/team.service";
import { CreateEnvironmentReconciliationDto, UpdateEnvironmentReconciliationDto, CreateComponentReconciliationDto, UpdateComponentReconciliationDto } from "./dtos/reconciliation.dto";

@Injectable()
export class ReconciliationService {
  private readonly s3h = S3Handler.instance();
  private readonly config = get();
  private readonly ckEnvironment = this.config.environment;
  private readonly logger = new Logger(ReconciliationService.name);

  constructor(
    @InjectRepository(EnvironmentReconcile)
    private readonly envReconRepo: Repository<EnvironmentReconcile>,
    @InjectRepository(ComponentReconcile)
    private readonly compReconRepo: Repository<ComponentReconcile>,
    @InjectRepository(Component)
    private readonly envSvc: EnvironmentService,
    private readonly teamSvc: TeamService,
  ) { }

  async createEnvRecon(org: Organization, team: Team, env: Environment, createEnv: CreateEnvironmentReconciliationDto): Promise<EnvironmentReconcile> {
    return this.envReconRepo.save({
      startDateTime: createEnv.startDateTime,
      environment: env,
      team,
      status: "initializing",
      organization: org
    });
  }

  async updateEnvRecon(envRecon: EnvironmentReconcile, mergeRecon: UpdateEnvironmentReconciliationDto): Promise<EnvironmentReconcile> {
    this.envReconRepo.merge(envRecon, mergeRecon);

    return this.envReconRepo.save(envRecon);
  }

  async getEnvReconByReconcileId(org: Organization, reconcileId: number, withEnv: boolean = false) {
    return this.envReconRepo.findOne({
      where: {
        reconcileId,
        organization: {
          id: org.id
        }
      },
      relations: {
        environment: withEnv
      }
    })
  }

  async getEnvReconByEnv(org: Organization, env: Environment) {
    return this.envReconRepo.findOne({
      where: {
        environment: {
          id: env.id
        },
        organization: {
          id: org.id
        }
      },
      order: {
        startDateTime: -1
      }
    })
  }

  async getSkippedEnvironments(org: Organization, team: Team, env: Environment, ignoreReconcileIds: number[]) {
    return this.envReconRepo.find({
      where: {
        endDateTime: IsNull(),
        status: Not(Equal('skipped_reconcile')),
        reconcileId: Not(In(ignoreReconcileIds)),
        team: {
          id: team.id
        },
        organization: {
          id: org.id
        },
        environment: {
          id: env.id
        }
      },
    });
  }

  async bulkUpdateEnvironmentEntries(entries: EnvironmentReconcile[], status: string) {
    if (entries.length > 0) {
      const newEntries = entries.map((entry) => ({
        ...entry,
        status
      })) as any;

      this.logger.log({message: 'updating skipped workflows', environment: entries[0].environment});

      await this.envReconRepo.save(newEntries);
    }
  }

  async createCompRecon(org: Organization, envRecon: EnvironmentReconcile, createComp: CreateComponentReconciliationDto): Promise<ComponentReconcile> {
    return this.compReconRepo.save({
      name: createComp.name,
      status: 'initializing',
      organization: org,
      startDateTime: createComp.startDateTime,
      environmentReconcile: envRecon
    });
  }

  async updateCompRecon(compRecon: ComponentReconcile, mergeRecon: UpdateComponentReconciliationDto): Promise<ComponentReconcile> {
    this.compReconRepo.merge(compRecon, mergeRecon);

    return this.compReconRepo.save(compRecon);
  }

  async getCompReconById(org: Organization, reconcileId: number, withEnvRecon: boolean = false): Promise<ComponentReconcile> {
    return this.compReconRepo.findOne({
      where: {
        reconcileId,
        organization: {
          id: org.id
        }
      },
      relations: {
        environmentReconcile: withEnvRecon
      }
    })
  }

  async getCompReconByName(org: Organization, env: Environment, compName: string): Promise<ComponentReconcile> {
    return this.compReconRepo.findOne({
      where: {
        name: compName,
        environmentReconcile: {
          environment: {
            id: env.id
          }
        },
        organization: {
          id: org.id
        }
      },
      order: {
        startDateTime: -1
      },
    })
  }

  async getSkippedComponents(org: Organization, env: Environment, ignoreReconcileIds: number[]) {
    return await this.compReconRepo.find({
      where: {
        endDateTime: IsNull(),
        status: Not(Equal('skipped_reconcile')),
        reconcileId: Not(In(ignoreReconcileIds)),
        environmentReconcile: {
          environment: {
            id: env.id
          }
        },
        organization: {
          id: org.id
        }
      },
    });
  }

  async bulkUpdateComponentEntries(entries: ComponentReconcile[], status: string) {
    if (entries.length > 0) {
      const newEntries = entries.map((entry) => ({
        ...entry,
        status
      })) as any;

      await this.compReconRepo.save(newEntries);
    }
  }

  async patchApprovedBy(org: Organization, body: ApprovedByDto): Promise<ComponentReconcile> {
    const envRecon = await this.getEnvReconByReconcileId(org, body.envReconcileId);
    const compRecon = await this.getLatestCompReconcile(org, envRecon, body.compName);

    this.compReconRepo.merge(compRecon, { approved_by: body.email })

    return this.compReconRepo.save(compRecon);
  }

  // async putComponent(org: Organization, team: Team, component: ComponentDto, env: Environment) {
  //   if (typeof env === 'undefined') {
  //     env = await this.envSvc.findByName(org, component.environmentName, team);
  //   }

  //   const componentName = `${teamName}-${component.environmentName}-${component.name}`;

  //   const existing = await this.componentRepository
  //     .createQueryBuilder()
  //     .where('id = :id and component_name = :name and environmentId = :envId', {
  //       id: componentName,
  //       name: component.name,
  //       envId: env.id
  //     })
  //     .getOne();

  //   this.logger.log({ message: "PutComponent", existing, id: componentName, teamName, component, env});

  //   if (!existing) {
  //     const entry = await this.componentRepository.save({
  //       id: componentName,
  //       componentName: component.name,
  //       duration: component.duration,
  //       teamName: teamName,
  //       environment: env,
  //       status: component.status,
  //       organization: org
  //     });
  //     return entry;
  //   }

  //   existing.duration = component.duration;
  //   const entry = await this.componentRepository.save(existing);
    
  //   return entry;
  // }

  // async saveOrUpdateComponent(org: Organization, envReconcile: EvnironmentReconcileDto) {
  //   const reconcileId = Number.isNaN(parseInt(envReconcile.reconcileId))
  //     ? null
  //     : parseInt(envReconcile.reconcileId);

  //   if (!reconcileId) {
  //     this.logger.error({ message: 'No reconcileId found when trying to save or update a component', org, envReconcile });
  //     throw new BadRequestException("reconcileId is mandatory to save or update component")
  //   }

  //   this.logger.log({message: 'save or update component', reconcileId, envReconcile});

  //   const env = await this.envSvc.findById(org, envReconcile.name, envReconcile.teamName);

  //   if (!env) {
  //     throw new BadRequestException(`could not find environment ${envReconcile.name}`);
  //   }

  //   this.logger.log({message: 'found environment', reconcileId, env});

  //   const envRecEntry = await this.environmentReconcileRepository
  //     .createQueryBuilder()
  //     .where('reconcile_id = :reconcileId and environmentId = :envId',
  //       {
  //         reconcileId,
  //         envId: env.id
  //       }
  //     )
  //     .getOne();

  //   if (!envRecEntry) {
  //     throw new BadRequestException(`could not find environmentReconcileEntry for environment ${envReconcile.name}`);
  //   }

  //   this.logger.log({message: 'found environment reconcile entry', reconcileId, envRecEntry});

  //   let componentEntry: ComponentReconcile = Mapper.mapToComponentReconcile(
  //     org,
  //     envRecEntry,
  //     envReconcile.componentReconciles
  //   )[0];

  //   this.logger.log({message: 'created component entry', reconcileId, env, componentEntry});

  //   if (!componentEntry.reconcileId) {
  //     await this.updateSkippedWorkflowsForComponent(
  //       envReconcile.teamName,
  //       envReconcile.name,
  //       componentEntry.name,
  //       this.componentReconcileRepository
  //     );
  //   }

  //   let duration = -1;
  //   if (componentEntry.reconcileId) {
  //     const existingEntry = await this.componentReconcileRepository.findOne(
  //       {
  //         where: {
  //           reconcile_id: componentEntry.reconcileId
  //         }
  //       }
  //     );
  //     existingEntry.end_date_time = componentEntry.end_date_time;
  //     existingEntry.status = componentEntry.status;
  //     componentEntry = existingEntry;
  //     const ed = new Date(existingEntry.end_date_time).getTime();
  //     const sd = new Date(existingEntry.start_date_time).getTime();
  //     duration = ed - sd;
  //   }

  //   if (envReconcile.componentReconciles.length == 0) {
  //     throw new BadRequestException("no component reconciles");
  //   }

  //   const comp = envReconcile.componentReconciles[0];

  //   await this.putComponent(org, {
  //     name: comp.name,
  //     duration,
  //     teamName: comp.teamName,
  //     environmentName: envRecEntry.name,
  //     status: comp.status,
  //   }, env, envReconcile.teamName);

  //   const entry = await this.componentReconcileRepository.save(componentEntry);
  //   entry.organization = org;
  //   entry.environmentReconcile = new EnvironmentReconcile();
  //   entry.environmentReconcile.name = envReconcile.name;
  //   entry.environmentReconcile.team_name = envReconcile.teamName;
  //   entry.status = componentEntry.status;
  //   this.sseSvc.sendComponentReconcile(entry);

  //   this.logger.log({message: 'created component reconcile entry', reconcileId, compReconcileEntry: entry});

  //   return entry.reconcile_id;
  // }

  // async updateSkippedWorkflows(entries: ComponentReconcile[]|EnvironmentReconcile[], repo: Repository<ComponentReconcile|EnvironmentReconcile>) {
  //   if (entries.length > 0) {
  //     const newEntries = entries.map((entry) => ({
  //       ...entry,
  //       status: "skipped_reconcile",
  //     })) as any;

  //     this.logger.log({message: 'updating skipped workflows', newEntries});

  //     await repo.save(newEntries);
  //   }
  // }

  async getComponentAuditList(org: Organization, team: Team, id: string, envName: string): Promise<ComponentAudit[]> {
    const envAuditList = await this.getEnvironmentAuditList(org, team, envName);
    const reconcileIds = envAuditList.map(e => e.reconcileId);
    const components = await this.compReconRepo.find({
      where: {
        name: id,
        organization: {
          id: org.id
        },
        environmentReconcile: {
          reconcileId: In(reconcileIds)
        }
      }
    });

    return Mapper.getComponentAuditList(components);
  }

  async getEnvironmentAuditList(org: Organization, team: Team, envName: string): Promise<EnvironmentAudit[]> {
    const env = await this.envSvc.findByName(org, team, envName);

    if (!env) {
      this.logger.error({ message: `Could not find environment with name ${envName}`, org })
      throw new BadRequestException(`Could not find environment`);
    }

    const environments = await this.envReconRepo.find({
      where: {
        environment: {
          id: env.id
        },
        organization: {
          id: org.id
        }
      }
    });

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

  // async getApplyLogs(
  //   org: Organization,
  //   team: string,
  //   environment: string,
  //   component: string,
  //   id: number,
  //   latest: boolean
  // ) {
  //   const logs =
  //     latest === true
  //       ? await this.getLatestLogs(org, team, environment, component)
  //       : await this.getLogs(org, team, environment, component, id);

  //   if (Array.isArray(logs)) {
  //     return logs.filter((e) => e.key.includes("apply_output"));
  //   }

  //   return logs;
  // }

  // async getPlanLogs(
  //   org: Organization,
  //   team: string,
  //   environment: string,
  //   component: string,
  //   id: number,
  //   latest: boolean
  // ) {
  //   const logs =
  //     latest === true
  //       ? await this.getLatestLogs(org, team, environment, component)
  //       : await this.getLogs(org, team, environment, component, id);

  //   if (Array.isArray(logs)) {
  //     return logs.filter((e) => e.key.includes("plan_output"));
  //   }

  //   return logs;
  // }

  // async getLatestLogs(
  //   org: Organization,
  //   env: Environment,
  //   compName: string
  // ) {    
  //   const envRecon = await this.getEnvReconByEnv(org, env);
  //   const compRecon = await this.getLatestCompReconcile(org, envRecon, compName);

  //   if (!compRecon) {
  //     this.logger.error({ message: 'could not find latestAuditId to get latest logs', component: compName, environment, team });
  //     throw new NotFoundException(`could not find audit logs for ${team}-${environment}-${compName}`);
  //   }

  //   return await this.getLogs(
  //     org,
  //     team,
  //     environment,
  //     compName,
  //     compRecon.reconcileId
  //   );
  // }

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

  private async getLatestCompReconcile(org: Organization, envRecon: EnvironmentReconcile, compName: string): Promise<ComponentReconcile> {
    return await this.compReconRepo.findOne({
      where: {
        name: compName,
        status: Not(Like("skipped%")),
        organization: {
          id : org.id
        },
        environmentReconcile: {
          reconcileId: envRecon.reconcileId
        }
      },
      order: {
        startDateTime: -1,
      }
    });
  }
}

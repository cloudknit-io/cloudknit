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

  async getComponentAuditList(org: Organization, env: Environment, comp: Component): Promise<ComponentAudit[]> {
    const envAuditList = await this.getEnvironmentAuditList(org, env);
    const reconcileIds = envAuditList.map(e => e.reconcileId);
    
    const components = await this.compReconRepo.find({
      where: {
        name: Equal(comp.name),
        organization: {
          id: org.id
        },
        environmentReconcile: {
          reconcileId: In(reconcileIds)
        }
      },
      order: {
        startDateTime: -1
      }
    });

    return Mapper.getComponentAuditList(components);
  }

  async getEnvironmentAuditList(org: Organization, env: Environment): Promise<EnvironmentAudit[]> {
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

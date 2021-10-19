import { Injectable } from "@nestjs/common";
import { InjectRepository } from "@nestjs/typeorm";
import { Subject } from "rxjs";
import { Mapper } from "src/costing/utilities/mapper";
import { S3Handler } from "src/costing/utilities/s3Handler";
import { ComponentReconcile } from "src/typeorm/reconciliation/component-reconcile.entity";
import { EnvironmentReconcile } from "src/typeorm/reconciliation/environment-reconcile.entity";
import { Notification } from "src/typeorm/reconciliation/notification.entity";
import { Repository } from "typeorm/repository/Repository";
import { ComponentAudit } from "../dtos/componentAudit.dto";
import { EnvironmentAudit } from "../dtos/environmentAudit.dto";
import { NotificationDto } from "../dtos/notification.dto";
import { EvnironmentReconcileDto } from "../dtos/reconcile.Dto";

@Injectable()
export class ReconciliationService {
  readonly notifyStream: Subject<{}> = new Subject<{}>();
  readonly notificationStream: Subject<Notification> =
    new Subject<Notification>();
  private readonly s3h = S3Handler.instance();
  private readonly notFound = "";
  constructor(
    @InjectRepository(EnvironmentReconcile)
    private readonly environmentReconcileRepository: Repository<EnvironmentReconcile>,
    @InjectRepository(ComponentReconcile)
    private readonly componentReconcileRepository: Repository<ComponentReconcile>,
    @InjectRepository(Notification)
    private readonly notificationRepository: Repository<Notification>
  ) {}

  async saveOrUpdateEnvironment(runData: EvnironmentReconcileDto) {
    const reconcileId = Number.isNaN(parseInt(runData.reconcileId))
      ? null
      : parseInt(runData.reconcileId);

    let savedEntry = null;
    if (reconcileId) {
      const existingEntry = await this.environmentReconcileRepository.findOne(
        reconcileId
      );
      existingEntry.end_date_time = runData.endDateTime;
      existingEntry.status = runData.status;
      savedEntry = await this.environmentReconcileRepository.save(
        existingEntry
      );
    } else {
      const entry: EnvironmentReconcile = {
        reconcile_id: reconcileId,
        name: runData.name,
        start_date_time: runData.startDateTime,
        team_name: runData.teamName,
        status: runData.status,
        end_date_time: runData.endDateTime,
      };
      await this.updateSkippedWorkflows<EnvironmentReconcile>(
        entry.name,
        this.environmentReconcileRepository
      );
      savedEntry = await this.environmentReconcileRepository.save(entry);
    }
    this.notifyStream.next(savedEntry);
    return savedEntry.reconcile_id;
  }

  async saveOrUpdateComponent(runData: EvnironmentReconcileDto) {
    const reconcileId = Number.isNaN(parseInt(runData.reconcileId))
      ? null
      : parseInt(runData.reconcileId);
    if (!reconcileId) {
      throw "Reconcile Id is mandatory to save or update component";
    }
    const savedEntry = await this.environmentReconcileRepository.findOne(
      reconcileId
    );

    let componentEntry: ComponentReconcile = Mapper.mapToComponentReconcile(
      savedEntry,
      runData.componentReconciles
    )[0];

    if (!componentEntry.reconcile_id) {
      await this.updateSkippedWorkflows<ComponentReconcile>(
        componentEntry.name,
        this.componentReconcileRepository
      );
    }

    if (componentEntry.reconcile_id) {
      const existingEntry = await this.componentReconcileRepository.findOne(
        componentEntry.reconcile_id
      );
      existingEntry.end_date_time = componentEntry.end_date_time;
      existingEntry.status = componentEntry.status;
      componentEntry = existingEntry;
    }

    const entry = await this.componentReconcileRepository.save(componentEntry);
    this.notifyStream.next(entry);
    return entry.reconcile_id;
  }

  async updateSkippedWorkflows<T>(id: any, repo: Repository<T>) {
    const entries = await repo.find({
      where: {
        name: id,
        end_date_time: null,
      },
    });
    if (entries.length > 0) {
      const newEntries = entries.map((entry) => ({
        ...entry,
        status: "Skipped",
      }));
      await repo.save(newEntries);
    }
  }

  async getComponentAuditList(id: string): Promise<ComponentAudit[]> {
    const components = await this.componentReconcileRepository.find({
      where: {
        name: id,
      },
    });
    return Mapper.getComponentAuditList(components);
  }

  async getEnvironmentAuditList(id: string): Promise<EnvironmentAudit[]> {
    const environments = await this.environmentReconcileRepository.find({
      where: {
        name: id,
      },
    });
    return Mapper.getEnvironmentAuditList(environments);
  }

  async getLogs(
    companyId: string,
    team: string,
    environment: string,
    component: string,
    id: number
  ) {
    try {
      const prefix = `${team}/${environment}/${component}/${id}/`;
      const objects = await this.s3h.getObjects(
        `zlifecycle-tfplan-${companyId}`,
        prefix
      );
      return objects.map((o) => ({
        key: o.key,
        body: o.data.Body.toString(),
      }));
    } catch (err) {
      if (err === this.notFound) {
        return err;
      }
    }
  }

  async getApplyLogs(
    companyId: string,
    team: string,
    environment: string,
    component: string,
    id: number,
    latest: boolean
  ) {
    const logs =
      latest === true
        ? await this.getLatestLogs(companyId, team, environment, component)
        : await this.getLogs(companyId, team, environment, component, id);
    if (Array.isArray(logs)) {
      return logs.filter((e) => e.key.includes("apply_output"));
    }
    return logs;
  }

  async getPlanLogs(
    companyId: string,
    team: string,
    environment: string,
    component: string,
    id: number,
    latest: boolean
  ) {
    const logs =
      latest === true
        ? await this.getLatestLogs(companyId, team, environment, component)
        : await this.getLogs(companyId, team, environment, component, id);
    if (Array.isArray(logs)) {
      return logs.filter((e) => e.key.includes("plan_output"));
    }
    return logs;
  }

  async getLatestLogs(
    companyId: string,
    team: string,
    environment: string,
    component: string
  ) {
    const latestAuditId = await this.componentReconcileRepository.find({
      where: {
        name: `${team}-${environment}-${component}`,
      },
      order: {
        start_date_time: -1,
      },
      take: 1,
    });
    if (latestAuditId.length === 0) {
      return this.notFound;
    }
    const logs = await this.getLogs(
      companyId,
      team,
      environment,
      component,
      latestAuditId[0].reconcile_id
    );
    if (Array.isArray(logs)) {
      return logs;
    }
    return logs;
  }

  async getStateFile(
    companyId: string,
    team: string,
    environment: string,
    component: string
  ) {
    const prefix = `${team}/${environment}/${component}/terraform.tfstate`;
    const resp = await this.s3h.getObject(
      `zlifecycle-tfstate-${companyId}`,
      prefix
    );

    return {
      ...resp,
      data: (resp.data?.Body || "").toString() || "",
    };
  }

  async saveNotification(notification: NotificationDto) {
    const notificationEntity: Notification = {
      company_id: notification.companyId,
      environment_name: notification.environmentName,
      message: notification.message,
      team_name: notification.teamName,
      timestamp: notification.timestamp,
      message_type: notification.messageType,
    };
    const savedEntity = await this.notificationRepository.save(
      notificationEntity
    );
    this.notificationStream.next(savedEntity);
  }

  getNotification(companyId: string, teamName: string) {
    this.notificationRepository
      .find({
        where: {
          company_id: companyId,
          team_name: teamName,
          seen: false,
        },
      })
      .then((notification) => {
        notification.forEach((e) => this.notificationStream.next(e));
      });
  }

  async getAllNotification(companyId: string, teamName: string) {
    return await this.notificationRepository.find({
      where: {
        company_id: companyId,
        team_name: teamName,
      },
      order: {
        timestamp: "DESC",
      },
    });
  }

  async setSeenStatusForNotification(notificationId: number) {
    await this.notificationRepository.save({
      notification_id: notificationId,
      seen: true,
    });
  }
}

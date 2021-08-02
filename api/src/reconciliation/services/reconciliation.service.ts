import { Injectable } from "@nestjs/common";
import { InjectRepository } from "@nestjs/typeorm";
import { S3 } from "aws-sdk";
import { Subject } from "rxjs";
import { Mapper } from "src/costing/utilities/mapper";
import { S3Handler } from "src/costing/utilities/s3Handler";
import { ComponentReconcile } from "src/typeorm/reconciliation/component-reconcile.entity";
import { EnvironmentReconcile } from "src/typeorm/reconciliation/environment-reconcile.entity";
import { Repository } from "typeorm/repository/Repository";
import { ComponentAudit } from "../dtos/componentAudit.dto";
import { EnvironmentAudit } from "../dtos/environmentAudit.dto";
import { EvnironmentReconcileDto } from "../dtos/reconcile.Dto";

@Injectable()
export class ReconciliationService {
  readonly notifyStream: Subject<{}> = new Subject<{}>();
  private readonly s3h = S3Handler.instance();
  constructor(
    @InjectRepository(EnvironmentReconcile)
    private readonly environmentReconcileRepository: Repository<EnvironmentReconcile>,
    @InjectRepository(ComponentReconcile)
    private readonly componentReconcileRepository: Repository<ComponentReconcile>
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
    team: string,
    environment: string,
    component: string,
    id: number
  ) {
    try {
      const prefix = `${team}/${environment}/${component}/${id}/`;
      const objects = await this.s3h.getObjects(
        "zlifecycle-tfplan-zmart",
        prefix
      );
      return objects.map((o) => ({
        key: o.key,
        body: o.data.Body.toString(),
      }));
    } catch (err) {
      if (err === "No Object was found") {
        return err;
      }
    }
  }

  async getLatestLogs(team: string, environment: string, component: string) {
    const latestAuditId = await this.componentReconcileRepository.find({
      where: {
        name: component,
      },
      order: {
        start_date_time: -1,
      },
      take: 1
    });
    if (latestAuditId.length > 0) {
      try {
        const prefix = `${team}/${environment}/${component}/${latestAuditId[0].reconcile_id}/`;
        const objects = await this.s3h.getObjects(
          "zlifecycle-tfplan-zmart",
          prefix
        );
        return objects.map((o) => ({
          key: o.key,
          body: o.data.Body.toString(),
        }));
      } catch (err) {
        if (err === "No Object was found") {
          return err;
        }
      }
    }
    return "No Object was found";
  }
}

import { ComponentAudit } from "src/reconciliation/dtos/componentAudit.dto";
import { EnvironmentAudit } from "src/reconciliation/dtos/environmentAudit.dto";
import { ComponentReconcile } from "src/typeorm/component-reconcile.entity";
import { EnvironmentReconcile } from "src/typeorm/environment-reconcile.entity";

export class Mapper {
  static getComponentAuditList(
    components: ComponentReconcile[]
  ): ComponentAudit[] {
    return components.map((c) => {
      let diff = -1;
      if (c.endDateTime) {
        const ed = new Date(c.endDateTime).getTime();
        const sd = new Date(c.startDateTime).getTime()
        diff = ed - sd;
      }

      return {
        reconcileId: c.reconcileId,
        startDateTime: c.startDateTime,
        duration: diff,
        status: c.status,
        approvedBy: c.approved_by,
      };
    });
  }

  static getEnvironmentAuditList(
    components: EnvironmentReconcile[]
  ): EnvironmentAudit[] {
    return components.map((c) => {
      let diff = -1;
      if (c.endDateTime) {
        const ed = new Date(c.endDateTime).getTime();
        const sd = new Date(c.startDateTime).getTime()
        diff = ed - sd;
      }

      return {
        reconcileId: c.reconcileId,
        startDateTime: c.startDateTime,
        duration: diff,
        status: c.status,
      };
    });
  }
}

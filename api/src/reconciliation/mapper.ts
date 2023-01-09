import { ComponentReconcileWrap } from "src/reconciliation/dtos/componentAudit.dto";
import { EnvironmentReconcileWrap } from "src/reconciliation/dtos/environmentAudit.dto";
import { ComponentReconcile } from "src/typeorm/component-reconcile.entity";
import { EnvironmentReconcile } from "src/typeorm/environment-reconcile.entity";

export class Mapper {
  static getComponentAuditList(
    components: ComponentReconcile[]
  ): ComponentReconcileWrap[] {
    return components.map((c) => {
      let diff = -1;
      if (c.endDateTime) {
        const ed = new Date(c.endDateTime).getTime();
        const sd = new Date(c.startDateTime).getTime()
        diff = ed - sd;
      }

      return {
        ...c,
        duration: diff
      };
    });
  }

  static getEnvironmentAuditList(
    components: EnvironmentReconcile[]
  ): EnvironmentReconcileWrap[] {
    return components.map((c) => {
      let diff = -1;
      if (c.endDateTime) {
        const ed = new Date(c.endDateTime).getTime();
        const sd = new Date(c.startDateTime).getTime()
        diff = ed - sd;
      }

      return {
        ...c,
        duration: diff
      };
    });
  }
}

import { ComponentAudit } from "src/reconciliation/dtos/componentAudit.dto";
import { EnvironmentAudit } from "src/reconciliation/dtos/environmentAudit.dto";
import { ComponentReconcileDto } from "src/reconciliation/dtos/reconcile.Dto";
import { Organization } from "src/typeorm";
import { ComponentReconcile } from "src/typeorm/component-reconcile.entity";
import { EnvironmentReconcile } from "src/typeorm/environment-reconcile.entity";

export class Mapper {
  static getResource(data: any) {
    return {
      name: data["name"],
      hourlyCost: data["hourlyCost"],
      monthlyCost: data["monthlyCost"],
      subresources: [],
      resourceName: data.resourceName,
      costComponents: [],
    };
  }

  static getCostComponent(data: any) {
    return {
      hourlyCost: data["hourlyCost"],
      hourlyQuantity: data["hourlyQuantity"],
      monthlyCost: data["monthlyCost"],
      monthlyQuantity: data["monthlyQuantity"],
      name: data["name"],
      price: data["price"],
      unit: data["unit"],
    };
  }

  static getComponentAuditList(
    components: ComponentReconcile[]
  ): ComponentAudit[] {
    return components.map((c) => {
      let diff = -1;
      if (c.end_date_time) {
        const ed = new Date(c.end_date_time).getTime();
        const sd = new Date(c.start_date_time).getTime()
        diff = ed - sd;
      }

      return {
        reconcileId: c.reconcileId,
        startDateTime: c.start_date_time,
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
      if (c.end_date_time) {
        const ed = new Date(c.end_date_time).getTime();
        const sd = new Date(c.start_date_time).getTime()
        diff = ed - sd;
      }

      return {
        reconcileId: c.reconcileId,
        startDateTime: c.start_date_time,
        duration: diff,
        status: c.status,
      };
    });
  }
}

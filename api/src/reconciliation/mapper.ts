import { ComponentReconcileWrap } from 'src/reconciliation/dtos/componentAudit.dto';
import { EnvironmentReconcileWrap } from 'src/reconciliation/dtos/environmentAudit.dto';
import { ComponentReconcile } from 'src/typeorm/component-reconcile.entity';
import { EnvironmentReconcile } from 'src/typeorm/environment-reconcile.entity';

export class Mapper {
  static wrapComponentRecon(compRecon: ComponentReconcile): ComponentReconcileWrap {
    let diff = -1;
    if (compRecon.endDateTime) {
      const ed = new Date(compRecon.endDateTime).getTime();
      const sd = new Date(compRecon.startDateTime).getTime();
      diff = ed - sd;
    }

    return {
      ...compRecon,
      duration: diff,
    };
  }

  static getComponentAuditList(compRecons: ComponentReconcile[]): ComponentReconcileWrap[] {
    return compRecons.map(c => this.wrapComponentRecon(c));
  }

  static wrapEnvironmentRecon(envRecon: EnvironmentReconcile): EnvironmentReconcileWrap {
    let diff = -1;
    if (envRecon.endDateTime) {
      const ed = new Date(envRecon.endDateTime).getTime();
      const sd = new Date(envRecon.startDateTime).getTime();
      diff = ed - sd;
    }

    return {
      ...envRecon,
      duration: diff,
    };
  }

  static getEnvironmentAuditList(envRecons: EnvironmentReconcile[]): EnvironmentReconcileWrap[] {
    return envRecons.map(envRecon => this.wrapEnvironmentRecon(envRecon));
  }
}

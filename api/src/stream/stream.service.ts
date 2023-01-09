import { Injectable } from "@nestjs/common";
import { Subject } from "rxjs";
import { ComponentReconcileWrap } from "src/reconciliation/dtos/componentAudit.dto";
import { EnvironmentReconcileWrap } from "src/reconciliation/dtos/environmentAudit.dto";
import { Mapper } from "src/reconciliation/mapper";
import { Component, ComponentReconcile, EnvironmentReconcile, Team } from "src/typeorm";
import { Environment } from "src/typeorm/environment.entity";

@Injectable()
export class StreamService {
  readonly envStream: Subject<Environment> = new Subject<Environment>();
  readonly compStream: Subject<Component> = new Subject<Component>();
  readonly reconcileStream: Subject<AuditWrapper> = new Subject<AuditWrapper>();

  constructor() { }

  normalizeOrg(obj: Environment|Component|ComponentReconcile|EnvironmentReconcile) {
    if (obj && obj.orgId) {
      delete obj.organization;

      return obj;
    }

    if (obj.organization) {
      obj.orgId = obj.organization.id;
      delete obj.organization;

      return obj;
    }

    return obj;
  }

  sendEnvironment(env: Environment) {
    this.envStream.next(this.normalizeOrg(env) as Environment);
  }

  sendComponent(comp: Component) {
    this.compStream.next(this.normalizeOrg(comp) as Component);
  }

  sendCompReconcile(compRecon: ComponentReconcile) {
    const data = Mapper.wrapComponentRecon(this.normalizeOrg(compRecon) as ComponentReconcile);

    this.reconcileStream.next({data, type: 'ComponentReconcile'});
  }

  sendEnvReconcile(envRecon: EnvironmentReconcile) {
    const data = Mapper.wrapEnvironmentRecon(this.normalizeOrg(envRecon) as EnvironmentReconcile);

    this.reconcileStream.next({data, type: 'EnvironmentReconcile'});
  }
}

export class AuditWrapper {
  data: EnvironmentReconcileWrap|ComponentReconcileWrap;
  type: 'EnvironmentReconcile'|'ComponentReconcile';
}

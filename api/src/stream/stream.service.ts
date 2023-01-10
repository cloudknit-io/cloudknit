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
    if (env.team) {
      env.teamId = env.team.id;

      delete env.team;
    }

    this.envStream.next(this.normalizeOrg(env) as Environment);
  }

  sendComponent(comp: Component) {
    if (comp.environment) {
      comp.envId = comp.environment.id;

      delete comp.environment;
    }

    this.compStream.next(this.normalizeOrg(comp) as Component);
  }

  sendCompReconcile(compRecon: ComponentReconcile) {
    const data = Mapper.wrapComponentRecon(this.normalizeOrg(compRecon) as ComponentReconcile);

    if (data.environmentReconcile) {
      data.envReconId = data.environmentReconcile.reconcileId;

      delete data.environmentReconcile;
    }

    if (data.component) {
      data.compId = data.component.id;

      delete data.component;
    }

    this.reconcileStream.next({data, type: 'ComponentReconcile'});
  }

  sendEnvReconcile(envRecon: EnvironmentReconcile) {
    const data = Mapper.wrapEnvironmentRecon(this.normalizeOrg(envRecon) as EnvironmentReconcile);

    if (data.environment) {
      data.envId = data.environment.id;

      delete data.environment;
    }

    if (data.team) {
      data.teamId = data.team.id;

      delete data.team;
    }

    this.reconcileStream.next({data, type: 'EnvironmentReconcile'});
  }
}

export class AuditWrapper {
  data: EnvironmentReconcileWrap|ComponentReconcileWrap;
  type: 'EnvironmentReconcile'|'ComponentReconcile';
}

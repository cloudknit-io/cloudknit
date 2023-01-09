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

  sendEnvironment(env: Environment) {
    this.envStream.next(env);
  }

  sendComponent(comp: Component) {
    this.compStream.next(comp);
  }

  sendCompReconcile(compRecon: ComponentReconcile) {
    const data = Mapper.wrapComponentRecon(compRecon);

    this.reconcileStream.next({data, type: 'ComponentReconcile'});
  }

  sendEnvReconcile(envRecon: EnvironmentReconcile) {
    const data = Mapper.wrapEnvironmentRecon(envRecon);

    this.reconcileStream.next({data, type: 'EnvironmentReconcile'});
  }
}

export class AuditWrapper {
  data: EnvironmentReconcileWrap|ComponentReconcileWrap;
  type: 'EnvironmentReconcile'|'ComponentReconcile';
}

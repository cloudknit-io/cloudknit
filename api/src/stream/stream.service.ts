import { Injectable } from "@nestjs/common";
import { Subject } from "rxjs";
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
    this.reconcileStream.next({data: compRecon, type: 'ComponentReconcile'});
  }

  sendEnvReconcile(envRecon: EnvironmentReconcile) {
    this.reconcileStream.next({data: envRecon, type: 'EnvironmentReconcile'});
  }
}

export class AuditWrapper {
  data: EnvironmentReconcile|ComponentReconcile;
  type: 'EnvironmentReconcile'|'ComponentReconcile';
}

import { Injectable } from "@nestjs/common";
import { Subject } from "rxjs";
import { Component, ComponentReconcile, EnvironmentReconcile, Team } from "src/typeorm";
import { Environment } from "src/typeorm/environment.entity";

@Injectable()
export class StreamService {
  readonly notifyStream: Subject<{}> = new Subject<{}>();
  readonly envStream: Subject<Environment> = new Subject<Environment>();
  readonly compStream: Subject<Component> = new Subject<Component>();
  readonly teamStream: Subject<Team> = new Subject<Team>();
  readonly reconcileStream: Subject<ComponentReconcile | EnvironmentReconcile> = new Subject<ComponentReconcile | EnvironmentReconcile>();

  constructor() {
    // setInterval(() => {
    //   this.notifyStream.next({});
    //   this.envStream.next(null);
    // }, 20000);
  }

  sendEnvironment(env: Environment) {
    this.envStream.next(env);
  }

  sendComponent(comp: Component) {
    this.compStream.next(comp);
  }

  sendTeam(team: Team) {
    this.teamStream.next(team);
  }

  sendCompReconcile(compRecon: ComponentReconcile) {
    this.reconcileStream.next(compRecon);
  }

  sendEnvReconcile(envRecon: EnvironmentReconcile) {
    this.reconcileStream.next(envRecon);
  }
}

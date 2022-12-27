import { Injectable, Logger } from "@nestjs/common";
import { Subject } from "rxjs";
import { get } from "src/config";
import { ComponentReconcile } from "src/typeorm/component-reconcile.entity";
import { EnvironmentReconcile } from "src/typeorm/environment-reconcile.entity";
import { Environment } from "src/typeorm/environment.entity";

@Injectable()
export class SSEService {
  readonly notifyStream: Subject<{}> = new Subject<{}>();
  readonly applicationStream: Subject<any> = new Subject<any>();
  private readonly config = get();
  private readonly ckEnvironment = this.config.environment;
  private readonly logger = new Logger(SSEService.name);

  constructor() {
    setInterval(() => {
      this.notifyStream.next({});
      this.applicationStream.next({});
    }, 20000);
  }

  public async sendEnvironment(env: Environment) {
    this.applicationStream.next(env);
  }

  public async sendComponentReconcile(compRecon: ComponentReconcile) {
    this.notifyStream.next(compRecon);
  }

  public async sendEnvironmentReconcile(envRecon: EnvironmentReconcile) {
    this.notifyStream.next(envRecon);
  }
}

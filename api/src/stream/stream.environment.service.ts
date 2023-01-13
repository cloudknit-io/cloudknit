import { Inject, Injectable, Logger } from "@nestjs/common";
import { Environment } from "src/typeorm/environment.entity";
import { Connection, EntitySubscriberInterface, InsertEvent, RemoveEvent, UpdateEvent } from "typeorm";
import { StreamService } from "./stream.service";

@Injectable()
export class StreamEnvironmentService implements EntitySubscriberInterface<Environment> {
  private readonly logger = new Logger(StreamEnvironmentService.name); 

  constructor(
    @Inject(Connection) conn: Connection,
    private readonly sseSvc: StreamService,
  ) {
    conn.subscribers.push(this);
  }

  listenTo(): string | Function {
    return Environment;
  }

  afterInsert(event: InsertEvent<Environment>) {
    this.validateAndSend(event.entity, 'afterInsert');
  }

  afterUpdate(event: UpdateEvent<Environment>): void | Promise<Environment> {
    this.validateAndSend(event.entity as Environment, 'afterUpdate');
  }

  afterRemove(event: RemoveEvent<Environment>): void | Promise<any> {
    this.validateAndSend(event.entity, 'afterRemove');
  }

  validateAndSend(env: Environment, operation: string) {
    if (env.organization || env.orgId) {
      this.sseSvc.sendEnvironment(env);
      return;
    }

    this.logger.error({message: 'environment stream object has no organization', env, operation}); 
  }
}

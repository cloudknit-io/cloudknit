import { Inject, Injectable, Logger } from '@nestjs/common';
import { EnvironmentReconcile } from 'src/typeorm';
import {
  Connection,
  EntitySubscriberInterface,
  InsertEvent,
  RemoveEvent,
  UpdateEvent,
} from 'typeorm';
import { StreamService } from './stream.service';

@Injectable()
export class StreamEnvironmentReconcileService
  implements EntitySubscriberInterface<EnvironmentReconcile>
{
  private readonly logger = new Logger(StreamEnvironmentReconcileService.name);

  constructor(
    @Inject(Connection) conn: Connection,
    private readonly sseSvc: StreamService
  ) {
    conn.subscribers.push(this);
  }

  listenTo(): string | Function {
    return EnvironmentReconcile;
  }

  afterInsert(event: InsertEvent<EnvironmentReconcile>) {
    this.validateAndSend(event.entity, 'afterInsert');
  }

  afterUpdate(
    event: UpdateEvent<EnvironmentReconcile>
  ): void | Promise<EnvironmentReconcile> {
    this.validateAndSend(event.entity as EnvironmentReconcile, 'afterUpdate');
  }

  afterRemove(event: RemoveEvent<EnvironmentReconcile>): void | Promise<any> {
    // this.validateAndSend(event.databaseEntity, 'afterRemove');
  }

  validateAndSend(envRecon: EnvironmentReconcile, operation: string) {
    if (envRecon.organization || envRecon.orgId) {
      this.sseSvc.sendEnvReconcile(envRecon);
      return;
    }

    this.logger.error({
      message: 'component stream object has no organization',
      comp: envRecon,
      operation,
    });
  }
}

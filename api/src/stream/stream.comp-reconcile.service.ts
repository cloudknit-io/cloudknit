import { Inject, Injectable, Logger } from '@nestjs/common';
import { ComponentReconcile } from 'src/typeorm';
import { Connection, EntitySubscriberInterface, InsertEvent, RemoveEvent, UpdateEvent } from 'typeorm';
import { StreamService } from './stream.service';

@Injectable()
export class StreamComponentReconcileService implements EntitySubscriberInterface<ComponentReconcile> {
  private readonly logger = new Logger(StreamComponentReconcileService.name);

  constructor(@Inject(Connection) conn: Connection, private readonly sseSvc: StreamService) {
    conn.subscribers.push(this);
  }

  listenTo(): string | Function {
    return ComponentReconcile;
  }

  afterInsert(event: InsertEvent<ComponentReconcile>) {
    this.validateAndSend(event.entity, 'afterInsert');
  }

  afterUpdate(event: UpdateEvent<ComponentReconcile>): void | Promise<ComponentReconcile> {
    this.validateAndSend(event.entity as ComponentReconcile, 'afterUpdate');
  }

  afterRemove(event: RemoveEvent<ComponentReconcile>): void | Promise<any> {
    // this.validateAndSend(event.databaseEntity, 'afterRemove');
  }

  validateAndSend(compRecon: ComponentReconcile, operation: string) {
    if (compRecon.organization || compRecon.orgId) {
      this.sseSvc.sendCompReconcile(compRecon);
      return;
    }

    this.logger.error({
      message: 'component stream object has no organization',
      comp: compRecon,
      operation,
    });
  }
}

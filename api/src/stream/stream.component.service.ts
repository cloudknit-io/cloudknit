import { Inject, Injectable, Logger } from '@nestjs/common';
import { Component } from 'src/typeorm';
import { Connection, EntitySubscriberInterface, InsertEvent, RemoveEvent, UpdateEvent } from 'typeorm';
import { StreamService } from './stream.service';

@Injectable()
export class StreamComponentService implements EntitySubscriberInterface<Component> {
  private readonly logger = new Logger(StreamComponentService.name);

  constructor(@Inject(Connection) conn: Connection, private readonly sseSvc: StreamService) {
    conn.subscribers.push(this);
  }

  listenTo(): string | Function {
    return Component;
  }

  afterInsert(event: InsertEvent<Component>) {
    this.validateAndSend(event.entity, 'afterInsert');
  }

  afterUpdate(event: UpdateEvent<Component>): void | Promise<Component> {
    this.validateAndSend(event.entity as Component, 'afterUpdate');
  }

  afterRemove(event: RemoveEvent<Component>): void | Promise<any> {
    // this.validateAndSend(event.databaseEntity, 'afterRemove');
  }

  validateAndSend(comp: Component, operation: string) {
    if (comp.organization || comp.orgId) {
      this.sseSvc.sendComponent(comp);
      return;
    }

    this.logger.error({
      message: 'component stream object has no organization',
      comp,
      operation,
    });
  }
}

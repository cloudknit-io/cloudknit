import { Inject, Injectable, Logger } from '@nestjs/common';
import { EventEmitter2 } from '@nestjs/event-emitter';
import { Component } from 'src/typeorm';
import {
  Connection,
  EntitySubscriberInterface,
  EventSubscriber,
  InsertEvent,
  RemoveEvent,
  UpdateEvent,
} from 'typeorm';
import { StreamService } from './stream.service';

@Injectable()
@EventSubscriber()
export class StreamComponentService
  implements EntitySubscriberInterface<Component>
{
  private readonly logger = new Logger(StreamComponentService.name);

  constructor(
    @Inject(Connection) private conn: Connection,
    private readonly sseSvc: StreamService,
    private evtEmitter: EventEmitter2
  ) {
    conn.subscribers.push(this);
  }

  listenTo(): string | Function {
    return Component;
  }

  afterInsert(event: InsertEvent<Component>) {
    const comp = event.entity as Component;

    this.validateAndSend(comp, 'afterInsert');
  }

  async afterUpdate(event: UpdateEvent<Component>): Promise<void> {
    const comp = event.entity as Component;
    // TODO: Find an alternate solution, to the below fix.
    const repo = await this.conn.getRepository<Component>(Component);
    const compWithLatestRecon = await repo.findOne({
      where: {
        id: comp.id,
      },
      relations: {
        latestCompRecon: true,
      },
    });
    this.validateAndSend(compWithLatestRecon, 'afterUpdate');
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

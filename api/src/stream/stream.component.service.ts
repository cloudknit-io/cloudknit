import { Inject, Injectable, Logger } from '@nestjs/common';
import { EnvironmentService } from 'src/environment/environment.service';
import { Component } from 'src/typeorm';
import {
  Connection,
  EntitySubscriberInterface,
  InsertEvent,
  RemoveEvent,
  UpdateEvent,
} from 'typeorm';
import { StreamService } from './stream.service';

@Injectable()
export class StreamComponentService
  implements EntitySubscriberInterface<Component>
{
  private readonly logger = new Logger(StreamComponentService.name);

  constructor(
    @Inject(Connection) conn: Connection,
    private readonly sseSvc: StreamService,
    private readonly envSvc: EnvironmentService
  ) {
    conn.subscribers.push(this);
  }

  listenTo(): string | Function {
    return Component;
  }

  afterInsert(event: InsertEvent<Component>) {
    const comp = event.entity as Component;

    this.envSvc.updateCost(comp.organization, comp.environment);

    this.validateAndSend(event.entity, 'afterInsert');
  }

  afterUpdate(event: UpdateEvent<Component>): void | Promise<Component> {
    const comp = event.entity as Component;

    for (const col of event.updatedColumns) {
      if (col.propertyName === 'estimatedCost') {
        this.envSvc.updateCost(comp.organization, comp.environment);
        break;
      }
    }

    this.validateAndSend(comp, 'afterUpdate');
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

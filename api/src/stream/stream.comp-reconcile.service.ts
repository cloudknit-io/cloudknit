import { Inject, Injectable, Logger } from '@nestjs/common';
import { ComponentReconcile } from 'src/typeorm';
import {
  Connection,
  EntitySubscriberInterface,
  InsertEvent,
  RemoveEvent,
  UpdateEvent,
} from 'typeorm';
import { StreamService } from './stream.service';
import { EventEmitter2 } from '@nestjs/event-emitter';
import {
  ComponentReconcileCostUpdateEvent,
  ComponentReconcileEntityUpdateEvent,
  InternalEventType,
} from 'src/types';

@Injectable()
export class StreamComponentReconcileService
  implements EntitySubscriberInterface<ComponentReconcile>
{
  private readonly logger = new Logger(StreamComponentReconcileService.name);
  private readonly eventColumns = new Set([
    'estimatedCost',
    'status',
    'costResources',
    'isDestroyed',
  ]);

  constructor(
    @Inject(Connection) conn: Connection,
    private readonly sseSvc: StreamService,
    private evtEmitter: EventEmitter2
  ) {
    conn.subscribers.push(this);
  }

  listenTo(): string | Function {
    return ComponentReconcile;
  }

  afterInsert(event: InsertEvent<ComponentReconcile>) {
    this.validateAndSend(event.entity, 'afterInsert');
  }

  afterUpdate(
    event: UpdateEvent<ComponentReconcile>
  ): void | Promise<ComponentReconcile> {
    const compRecon = event.entity as ComponentReconcile;

    this.logger.log({
      message: `********* After Update in DB, recon ${compRecon.reconcileId}`,
      columns: event.updatedColumns.map((c) => c.propertyName),
    });

    if (
      event.updatedColumns.find((col) =>
        this.eventColumns.has(col.propertyName)
      )
    ) {
      if (
        event.updatedColumns.find((col) => col.propertyName === 'estimatedCost')
      ) {
        // this.evtEmitter.emit(
        //   InternalEventType.ComponentReconcileCostUpdate,
        //   new ComponentReconcileCostUpdateEvent({ ...compRecon })
        // );
      }
      // this.evtEmitter.emit(
      //   InternalEventType.ComponentReconcileEntityUpdate,
      //   new ComponentReconcileEntityUpdateEvent({
      //     ...compRecon,
      //   })
      // );
    }
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

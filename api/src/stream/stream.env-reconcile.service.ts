import { Inject, Injectable, Logger } from '@nestjs/common';
import { EventEmitter2 } from '@nestjs/event-emitter';
import { EnvironmentReconcile } from 'src/typeorm';
import {
  EnvironmentReconCostUpdateEvent,
  EnvironmentReconEnvUpdateEvent,
  InternalEventType
} from 'src/types';
import {
  Connection,
  EntitySubscriberInterface,
  InsertEvent,
  RemoveEvent,
  UpdateEvent
} from 'typeorm';
import { StreamService } from './stream.service';

@Injectable()
export class StreamEnvironmentReconcileService
  implements EntitySubscriberInterface<EnvironmentReconcile>
{
  private readonly logger = new Logger(StreamEnvironmentReconcileService.name);

  constructor(
    @Inject(Connection) conn: Connection,
    private readonly sseSvc: StreamService,
    private evtEmitter: EventEmitter2
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
    if (event.updatedColumns.length === 0) return;
    const envRecon = event.entity as EnvironmentReconcile;
    const costUpdated = event.updatedColumns.find(
      (col) => col.propertyName === 'estimatedCost'
    );
    const updateEnv = event.updatedColumns.some((col) =>
      ['status', 'endDateTime', 'errorMessage'].includes(col.propertyName)
    );

    if (costUpdated) {
      this.evtEmitter.emit(
        InternalEventType.EnvironmentReconCostUpdate,
        new EnvironmentReconCostUpdateEvent({ ...envRecon })
      );
    }

    if (updateEnv) {
      this.evtEmitter.emit(
        InternalEventType.EnvironmentReconEnvUpdate,
        new EnvironmentReconEnvUpdateEvent({ ...envRecon })
      );
    }
    
    this.validateAndSend(envRecon, 'afterUpdate');
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

import { Inject, Injectable, Logger } from '@nestjs/common';
import { EventEmitter2 } from '@nestjs/event-emitter';
import { Environment } from 'src/typeorm/environment.entity';
import { EnvironmentCostUpdateEvent, InternalEventType } from 'src/types';
import {
  Connection,
  EntitySubscriberInterface,
  InsertEvent,
  RemoveEvent,
  UpdateEvent,
} from 'typeorm';
import { StreamService } from './stream.service';

@Injectable()
export class StreamEnvironmentService
  implements EntitySubscriberInterface<Environment>
{
  private readonly logger = new Logger(StreamEnvironmentService.name);

  constructor(
    @Inject(Connection) conn: Connection,
    private readonly sseSvc: StreamService,
    private evtEmitter: EventEmitter2
  ) {
    conn.subscribers.push(this);
  }

  listenTo(): string | Function {
    return Environment;
  }

  afterInsert(event: InsertEvent<Environment>) {
    this.validateAndSend(event.entity, 'afterInsert');
  }

  async afterUpdate(
    event: UpdateEvent<Environment>
  ): Promise<Environment | void> {
    const env = event.entity as Environment;

    if (
      event.updatedColumns.find((col) => col.propertyName === 'estimatedCost')
    ) {
      this.evtEmitter.emit(
        InternalEventType.EnvironmentCostUpdate,
        new EnvironmentCostUpdateEvent({ ...env })
      );
    }

    this.validateAndSend(env, 'afterUpdate');
  }

  afterRemove(event: RemoveEvent<Environment>): void | Promise<any> {
    this.validateAndSend(event.entity, 'afterRemove');
  }

  validateAndSend(env: Environment, operation: string) {
    if (env.organization || env.orgId) {
      this.sseSvc.sendEnvironment(env);
      return;
    }

    this.logger.error({
      message: 'environment stream object has no organization',
      env,
      operation,
    });
  }
}

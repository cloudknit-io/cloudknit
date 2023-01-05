import { Inject } from "@nestjs/common";
import { StreamService } from "src/stream/stream.service";
import { EntitySubscriberInterface, EventSubscriber, InsertEvent, UpdateEvent } from "typeorm"
import { Environment } from "./environment.entity"

@EventSubscriber()
export class EnvironmentSubscriber implements EntitySubscriberInterface<Environment> {

  listenTo(): string | Function {
    return Environment;
  }

  afterInsert(event: InsertEvent<Environment>) {
    console.log('after insert');
  }

  afterUpdate(event: UpdateEvent<Environment>): void | Promise<any> {
  }
}

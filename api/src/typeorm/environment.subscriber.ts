import { EntitySubscriberInterface, EventSubscriber, InsertEvent } from "typeorm"
import { Environment } from "./environment.entity"

@EventSubscriber()
export class EnvironmentSubscriber implements EntitySubscriberInterface<Environment> {

  listenTo(): string | Function {
    return Environment;
  }
  /**
   * Called after entity insertion.
   */
  afterInsert(event: InsertEvent<Environment>) {
      console.log(`ENVIRONMENT INSERTED: `, event.entity)
  }
}

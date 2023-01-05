import { EntitySubscriberInterface, EventSubscriber, InsertEvent } from "typeorm"
import { Component } from "./component.entity"

@EventSubscriber()
export class ComponentSubscriber implements EntitySubscriberInterface<Component> {

  listenTo(): string | Function {
    return Component;
  }
  /**
   * Called after entity insertion.
   */
  afterInsert(event: InsertEvent<Component>) {
      console.log(`Component INSERTED: `, event.entity)
  }
}

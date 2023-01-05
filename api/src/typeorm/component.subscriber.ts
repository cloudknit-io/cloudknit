import { Injectable } from "@nestjs/common";
import { ModuleRef } from "@nestjs/core";
import { StreamService } from "src/stream/stream.service";
import { EntitySubscriberInterface, EventSubscriber, InsertEvent, UpdateEvent } from "typeorm"
import { Component } from "./component.entity"

@EventSubscriber()
@Injectable()
export class ComponentSubscriber implements EntitySubscriberInterface<Component> {
  constructor() {
  }

  listenTo(): string | Function {
    return Component;
  }

  afterInsert(event: InsertEvent<Component>) {
  }

  afterUpdate(event: UpdateEvent<Component>): void | Promise<any> {
  }
}

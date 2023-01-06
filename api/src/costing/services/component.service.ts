import { Injectable, Logger, NotFoundException } from "@nestjs/common";
import { InjectRepository } from "@nestjs/typeorm";
import { Subject } from "rxjs";
import { In, Repository } from "typeorm";
import { MessageEvent } from "@nestjs/common";
import { Component, Organization, Team } from "src/typeorm";
import { Environment } from "src/typeorm/environment.entity";
import { ComponentDto } from "../dtos/Component.dto";
import { EnvironmentService } from "src/environment/environment.service";

@Injectable()
export class ComponentService {
  readonly stream: Subject<{}> = new Subject<{}>();
  readonly notifyStream: Subject<MessageEvent> = new Subject<MessageEvent>();
  private readonly logger = new Logger(ComponentService.name);
  
  constructor(
    @InjectRepository(Component)
    private compRepo: Repository<Component>,
    private readonly envSvc: EnvironmentService,
  ) {
    setInterval(() => {
      this.notifyStream.next({ data: {} });
    }, 20000);
  }

  async batchCreate(org: Organization, env: Environment, names: string[]) {
    return await this.compRepo
    .createQueryBuilder()
    .useTransaction(true)
    .insert()
    .into(Component)
    .values(names.map(name => {
      return {
        name,
        environment: env,
        organization: org
      }
    }))
    .execute();
  }

  async batchDelete(org: Organization, env: Environment, comps: Component[]) {
    return this.compRepo.delete({
      id: In(comps.map(c => c.id)),
      organization: {
        id: org.id
      },
      environment: {
        id: env.id
      }
    })
  }

  async create(org: Organization, env: Environment, name: string): Promise<Component> {
    return this.compRepo.save({
      name,
      environment: env,
      organization: org
    });
  }

  async getAllForEnvironmentById(org: Organization, env: Environment): Promise<Component[]> {
    return this.compRepo.find({
      where: {
        organization: {
          id: org.id
        },
        environment: {
          id: env.id
        }
      }
    })
  }

  async getAll(org: Organization, isDestroyed: boolean = false): Promise<Component[]> {
    const components = await this.compRepo.find({
      where: {
        organization: {
          id: org.id
        },
        isDestroyed
      }
    });

    return components;
  }

  async findById(org: Organization, id: number, isDestroyed: boolean = false, relations?: {}): Promise<Component> {
    return await this.compRepo.findOne({
      where: {
        id,
        organization: {
          id: org.id
        },
        isDestroyed
      },
      relations
    });
  }

  async findByName(org: Organization, env: Environment, name: string, isDestroyed: boolean = false, relations?: {}): Promise<Component> {
    return await this.compRepo.findOne({
      where: {
        name,
        environment: {
          id: env.id
        },
        organization: {
          id: org.id
        },
        isDestroyed
      },
      relations
    });
  }

  async softDelete(component: Component): Promise<Component> {
    component.isDestroyed = true;
    component.estimatedCost = -1;
    return await this.compRepo.save(component);
  }
}

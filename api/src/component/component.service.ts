import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { ReconciliationService } from 'src/reconciliation/reconciliation.service';
import {
  Component, Environment,
  Organization
} from 'src/typeorm';
import { Equal, In, Repository } from 'typeorm';
import { ComponentWrap } from './component.dto';
import { UpdateComponentDto } from './dto/update-component.dto';

@Injectable()
export class ComponentService {
  constructor(
    @InjectRepository(Component)
    private compRepo: Repository<Component>,
    private recSvc: ReconciliationService
  ) {}

  async batchCreate(org: Organization, env: Environment, names: string[]) {
    return await this.compRepo
      .createQueryBuilder()
      .useTransaction(true)
      .insert()
      .into(Component)
      .values(
        names.map((name) => {
          return {
            name,
            environment: env,
            organization: org,
          };
        })
      )
      .execute();
  }

  async batchDelete(org: Organization, env: Environment, comps: Component[]) {
    return this.compRepo.delete({
      id: In(comps.map((c) => c.id)),
      organization: {
        id: org.id,
      },
      environment: {
        id: env.id,
      },
    });
  }

  async create(
    org: Organization,
    env: Environment,
    name: string
  ): Promise<Component> {
    return this.compRepo.save({
      name,
      environment: env,
      organization: org,
    });
  }

  async getAllForEnvironmentById(
    org: Organization,
    env: Environment,
    withEnv: boolean = false
  ): Promise<Component[]> {
    return this.compRepo.find({
      where: {
        organization: {
          id: org.id,
        },
        environment: {
          id: env.id,
        },
      },
      relations: {
        environment: withEnv,
      },
    });
  }

  async getAll(
    org: Organization,
    withEnv: boolean = false
  ): Promise<Component[]> {
    const components = await this.compRepo.find({
      where: {
        organization: {
          id: org.id,
        },
      },
      relations: {
        environment: withEnv,
      },
    });

    return components;
  }

  async findById(
    org: Organization,
    id: number,
    withEnv: boolean = false
  ): Promise<Component> {
    return await this.compRepo.findOne({
      where: {
        id,
        organization: {
          id: org.id,
        },
      },
      relations: {
        environment: withEnv,
      },
    });
  }

  async findByName(
    org: Organization,
    env: Environment,
    name: string,
    withEnv: boolean = false
  ): Promise<Component> {
    return await this.compRepo.findOne({
      where: {
        name,
        environment: {
          id: env.id,
        },
        organization: {
          id: org.id,
        },
      },
      relations: {
        environment: withEnv,
      },
    });
  }

  async findByNameWithTeamName(
    org: Organization,
    teamName: string,
    envName: string,
    name: string,
    withEnv: boolean = false
  ): Promise<Component> {
    return await this.compRepo.findOne({
      where: {
        name: Equal(name),
        environment: {
          name: Equal(envName),
          team: {
            name: Equal(teamName),
          },
        },
        organization: {
          id: org.id,
        },
      },
      relations: {
        environment: withEnv,
      },
    });
  }

  async findAll(org: Organization, env: Environment, withEnv: boolean = false) {
    return this.compRepo.find({
      where: {
        environment: {
          id: env.id,
        },
        organization: {
          id: org.id,
        },
      },
      relations: {
        environment: withEnv,
      },
    });
  }

  async findAllWithLastCompRecon(org: Organization, env: Environment) {
    type compReconType = {
      id: number;
      status: string;
    };
    const components = await this.findAll(org, env);
    const compRecons: compReconType[] =
      await this.recSvc.getLatestCompReconByComponentIds(
        components.map((c) => c.id)
      );
    return components.map((c) => {
      const cw: ComponentWrap = {
        ...c,
        lastAuditStatus: compRecons.find((e) => e.id === c.id)?.status,
      };
      return cw;
    });
  }

  async update(
    comp: Component,
    mergeComp: UpdateComponentDto
  ): Promise<Component> {
    this.compRepo.merge(comp, mergeComp);

    return this.compRepo.save(comp);
  }
}

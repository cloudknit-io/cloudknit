import { Injectable, Logger } from '@nestjs/common';
import { Subject } from 'rxjs';
import { Mapper } from 'src/reconciliation/mapper';
import {
  Component,
  ComponentReconcile,
  Environment,
  EnvironmentReconcile,
  Team,
} from 'src/typeorm';
import { StreamItem, StreamTypeEnum } from './dto/stream-item.dto';
import { RedisClientType, createClient } from 'redis';
import { get } from 'src/config';

@Injectable()
export class StreamService {
  readonly webStream: Subject<StreamItem> = new Subject<StreamItem>();
  readonly logger = new Logger('StreamService');
  redisClient: any;
  apiStreamChannel = 'api-stream-channel';

  constructor() {
    this.connectToRedis().then(() => {
      this.startPublishingToRedis();
    });
  }

  normalizeOrg(
    obj:
      | Team
      | Environment
      | Component
      | ComponentReconcile
      | EnvironmentReconcile
  ) {
    if (obj && obj.orgId) {
      delete obj.organization;

      return obj;
    }

    if (obj.organization) {
      obj.orgId = obj.organization.id;
      delete obj.organization;

      return obj;
    }

    return obj;
  }

  sendEnvironment(env: Environment) {
    if (env.team) {
      env.teamId = env.team.id;

      delete env.team;
    }

    const payload = {
      data: this.normalizeOrg(env) as Environment,
      type: StreamTypeEnum.Environment,
    };

    this.webStream.next(payload);
  }

  sendTeam(team: Team) {
    const payload = {
      data: this.normalizeOrg(team) as Team,
      type: StreamTypeEnum.Team,
    };

    this.webStream.next(payload);
  }

  sendComponent(comp: Component) {
    if (comp.environment) {
      comp.envId = comp.environment.id;

      delete comp.environment;
    }

    const payload = {
      data: this.normalizeOrg(comp) as Component,
      type: StreamTypeEnum.Component,
    };

    this.webStream.next(payload);
  }

  sendCompReconcile(compRecon: ComponentReconcile) {
    const data = Mapper.wrapComponentRecon(
      this.normalizeOrg(compRecon) as ComponentReconcile
    );

    if (data.environmentReconcile) {
      data.envReconId = data.environmentReconcile.reconcileId;

      delete data.environmentReconcile;
    }

    if (data.component) {
      data.compId = data.component.id;

      delete data.component;
    }

    const payload = {
      data,
      type: StreamTypeEnum.ComponentReconcile,
    };

    this.webStream.next(payload);
  }

  sendEnvReconcile(envRecon: EnvironmentReconcile) {
    const data = Mapper.wrapEnvironmentRecon(
      this.normalizeOrg(envRecon) as EnvironmentReconcile
    );

    if (data.environment) {
      data.envId = data.environment.id;

      delete data.environment;
    }

    if (data.team) {
      data.teamId = data.team.id;

      delete data.team;
    }

    const payload = {
      data,
      type: StreamTypeEnum.EnvironmentReconcile,
    };

    this.webStream.next(payload);
  }

  async connectToRedis() {
    try {
      const { url, password } = get().redis;
      const client = createClient({
        url,
        password,
      });
      client.on('error', (err) => this.logger.error('Redis Client Error', err));
      await client.connect();
      if (client.isReady) {
        this.redisClient = client;
      }
    } catch (err) {
      this.logger.error('Could not connect to redis server');
    }
  }

  startPublishingToRedis() {
    if (!this.redisClient) return;
    this.webStream.subscribe((item: StreamItem) => {
      let msg = null;
      if (!item || !item.data) {
        msg = {
          data: {},
          type: StreamTypeEnum.Empty,
        };
      }

      msg = {
        data: item.data,
        type: item.type,
      };

      this.publishToRedis(msg);
    });
  }

  publishToRedis(msg: any) {
    this.redisClient.PUBLISH(this.apiStreamChannel, JSON.stringify(msg));
  }
}

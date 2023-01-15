import { Injectable } from '@nestjs/common';
import { default as Redis } from 'ioredis';

@Injectable()
export class CachingService {
  private readonly redis: Redis;

  constructor() {
    // this.redis = new Redis({
    //   port: 6379,
    //   host: 'redis',
    // });
  }

  public async publish(msg: any) {
    const channel = `api-channel-1`;

    this.redis.publish(channel, JSON.stringify(msg));
  }
}

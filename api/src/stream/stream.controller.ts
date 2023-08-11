import { Controller, Sse, Request } from '@nestjs/common';
import { StreamService } from './stream.service';
import { from, map, Observable } from 'rxjs';
import { OrgApiParam } from 'src/types';
import { ApiTags } from '@nestjs/swagger';
import { StreamItem, StreamTypeEnum } from './dto/stream-item.dto';
import { createClient } from 'redis';
import { RedisClient } from 'ioredis/built/connectors/SentinelConnector/types';

@Controller({
  version: '1',
})
@ApiTags('stream')
export class StreamController {
  constructor(private readonly sseSvc: StreamService) {}
}

interface MessageEvent {
  data: string | object;
  id?: string;
  type?: string;
  retry?: number;
}

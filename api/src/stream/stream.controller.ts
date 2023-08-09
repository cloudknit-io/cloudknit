import { Controller, Sse, Request } from '@nestjs/common';
import { StreamService } from './stream.service';
import { from, map, Observable } from 'rxjs';
import { OrgApiParam } from 'src/types';
import { ApiTags } from '@nestjs/swagger';
import { StreamItem, StreamTypeEnum } from './dto/stream-item.dto';
import { createClient } from 'redis';

@Controller({
  version: '1',
})
@ApiTags('stream')
export class StreamController {
  constructor(private readonly sseSvc: StreamService) {
    const redis = createClient({
      url: "redis://172.16.244.62:6379"
    });
    redis.connect().then(() => {
      this.sseSvc.webStream.subscribe((item: StreamItem) => {
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

        this.sseSvc.publishToRedis(msg, redis);
      });
    });
  }

  @Sse()
  @OrgApiParam()
  orgStream(@Request() req): Observable<MessageEvent> {
    const { org } = req;

    return from(this.sseSvc.webStream).pipe(
      map((item: StreamItem) => {
        // if (!item || !item.data || item.data.orgId !== org.id) {
        return {
          data: {},
          type: StreamTypeEnum.Empty,
        };
        // }

        return {
          data: item.data,
          type: item.type,
        };
      })
    );
  }
}

interface MessageEvent {
  data: string | object;
  id?: string;
  type?: string;
  retry?: number;
}

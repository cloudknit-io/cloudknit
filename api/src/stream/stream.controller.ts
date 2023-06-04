import { Controller, Sse, Request } from '@nestjs/common';
import { StreamService } from './stream.service';
import { from, map, Observable } from 'rxjs';
import { OrgApiParam } from 'src/types';
import { ApiTags } from '@nestjs/swagger';
import { StreamItem, StreamTypeEnum } from './dto/stream-item.dto';

@Controller({
  version: '1',
})
@ApiTags('stream')
export class StreamController {
  constructor(private readonly sseSvc: StreamService) {}

  @Sse()
  @OrgApiParam()
  orgStream(@Request() req): Observable<MessageEvent> {
    const { org } = req;

    return from(this.sseSvc.webStream).pipe(
      map((item: StreamItem) => {
        if (!item || !item.data || item.data.orgId !== org.id) {
          return {
            data: {},
            type: StreamTypeEnum.Empty,
          };
        }

        if (item.type === StreamTypeEnum.Component) {
          console.log(`${new Date().toISOString()}: stream svc: component: `, item.data);
        }
        if (item.type === StreamTypeEnum.ComponentReconcile) {
          console.log(`${new Date().toISOString()}: stream svc: component reconcile: `, item.data);
        }

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

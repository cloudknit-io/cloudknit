import { Controller, Sse, Request } from '@nestjs/common';
import { AuditWrapper, StreamService } from './stream.service';
import { from, map, Observable } from 'rxjs';
import { Component, Environment } from 'src/typeorm';
import { OrgApiParam } from 'src/types';

@Controller({
  version: '1',
})
export class StreamController {
  constructor(private readonly sseSvc: StreamService) {}

  @Sse('component')
  @OrgApiParam()
  components(@Request() req): Observable<MessageEvent> {
    const { org } = req;

    return from(this.sseSvc.compStream).pipe(
      map((comp: Component) => {
        if (!comp || comp.orgId !== org.id) {
          return {
            data: {},
            type: 'Component',
          };
        }

        return {
          data: comp,
          type: 'Component',
        };
      })
    );
  }

  @Sse('audit')
  @OrgApiParam()
  componentReconcile(@Request() req): Observable<MessageEvent> {
    const { org } = req;

    return from(this.sseSvc.reconcileStream).pipe(
      map((item: AuditWrapper) => {
        const { data } = item;

        if (!data || data.orgId !== org.id) {
          return {
            data,
            type: item.type,
          };
        }

        return {
          data,
          type: item.type,
        };
      })
    );
  }

  @Sse('environment')
  @OrgApiParam()
  environments(@Request() req): Observable<MessageEvent> {
    const { org } = req;

    return from(this.sseSvc.envStream).pipe(
      map((env: Environment) => {
        if (!env || env.orgId !== org.id) {
          return {
            data: {},
            type: 'Environment',
          };
        }

        delete env.team;

        return {
          data: env,
          type: 'Environment',
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

import { Controller, Sse, Request } from '@nestjs/common';
import { AuditWrapper, StreamService } from './stream.service';
import { from, map, Observable } from 'rxjs';
import { Component, ComponentReconcile, Environment, EnvironmentReconcile } from 'src/typeorm';

@Controller({
  version: '1'
})
export class StreamController {
  constructor(private readonly sseSvc: StreamService) {}

  @Sse("component")
  components(@Request() req): Observable<MessageEvent> {
    const { org } = req;

    return from(this.sseSvc.compStream).pipe(
      map((comp: Component) => {
        if (!comp || comp.orgId !== org.id) {
          return;
        }

        return {
          data: comp,
          type: 'Component'
        }
      })
    );
  }

  @Sse("audit")
  componentReconcile(@Request() req): Observable<MessageEvent> {
    const { org } = req;

    return from(this.sseSvc.reconcileStream).pipe(
      map((item: AuditWrapper) => {
        const { data } = item;

        if (!data || data.orgId !== org.id) {
          return;
        }

        return {
          data,
          type: item.type
        }
      })
    );
  }

  @Sse("environment")
  environments(@Request() req): Observable<MessageEvent> {
    const { org } = req;

    return from(this.sseSvc.envStream).pipe(
      map((env: Environment) => {
        if (!env || env.orgId !== org.id) {
          return;
        }

        delete env.team;

        return {
          data: env,
          type: 'Environment'
        }
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

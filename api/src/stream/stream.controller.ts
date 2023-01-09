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
          if (comp.organization && comp.organization.id === org.id) {
            comp.orgId = comp.organization.id;
          } else {
            return;
          }
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
          if (data.organization && data.organization.id === org.id) {
            data.orgId = data.organization.id;
          } else {
            return;
          }
        }

        // @ts-ignore
        delete data.environment;

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
          if (env.organization && env.organization.id === org.id) {
            env.orgId = env.organization.id;
          } else {
            return;
          }
        }

        delete env.team;
        delete env.organization;

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

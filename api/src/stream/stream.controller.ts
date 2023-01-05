import { Controller, Sse, Request, Query } from '@nestjs/common';
import { StreamService } from './stream.service';
import { from, map, noop, Observable } from 'rxjs';
import { RequiredQueryValidationPipe, TeamEnvCompQueryParams, TeamEnvQueryParams } from 'src/reconciliation/validationPipes';
import { Component, Environment } from 'src/typeorm';

@Controller({
  version: '1'
})
export class StreamController {
  constructor(private readonly sseSvc: StreamService) {}

  @Sse("components")
  components(@Request() req, @Query(new RequiredQueryValidationPipe()) tec: TeamEnvCompQueryParams): Observable<MessageEvent> {
    const { org } = req;

    return from(this.sseSvc.compStream).pipe(
      map((comp: Component) => {
        if (!comp.organization || comp.organization.id !== org.id) {
          return;
        }

        return {
          data: comp,
          type: 'Component'
        }
      })
    );
  }

  @Sse("environments")
  environments(@Request() req, @Query(new RequiredQueryValidationPipe()) te: TeamEnvQueryParams): Observable<MessageEvent> {
    const { org } = req;

    return from(this.sseSvc.envStream).pipe(
      map((env: Environment) => {
        if (!env.organization || env.organization.id !== org.id) {
          return;
        }

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

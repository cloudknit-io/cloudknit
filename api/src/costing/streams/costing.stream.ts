import { Controller, Param, Request, Sse } from '@nestjs/common'
import { from, Observable, Observer } from 'rxjs'
import { map } from 'rxjs/operators'
import { Component } from 'src/typeorm/component.entity'
import { ComponentService } from '../services/component.service'
import { Mapper } from '../utilities/mapper'
import { MessageEvent } from '@nestjs/common';

@Controller({
  path: 'stream',
  version: '1'
})
export class CostingStream {
  constructor(private readonly componentService: ComponentService) {}

  @Sse('all')
  componentCost(@Request() req): Observable<MessageEvent> {
    this.componentService.getAll(req.org);

    return from(this.componentService.stream).pipe(
      map((components: Component[]) => ({
        data: Mapper.getStreamData(components),
      })),
    )
  }

  @Sse('team/:teamId')
  teamCost(@Request() req, @Param('teamId') teamId: string): Observable<MessageEvent> {
    this.componentService.getTeamCost(req.org, teamId);

    return from(this.componentService.stream).pipe(
      map((components: Component[]) => ({
        data: Mapper.getStreamData(components),
      })),
    )
  }

  @Sse('environment/:teamId/:environmentId')
  environmentCost(
    @Request() req,
    @Param('teamId') teamId: string,
    @Param('environmentId') environmentId: string,
  ): Observable<MessageEvent> {
    this.componentService.getEnvironmentCost(req.org, teamId, environmentId)
    return from(this.componentService.stream).pipe(
      map((components: Component[]) => ({
        data: Mapper.getStreamData(components),
      })),
    )
  }

  @Sse('notify')
  notify(): Observable<MessageEvent> {
    return this.componentService.notifyStream.asObservable();
  }
}

export interface CostingStreamDto {
  team: {
    teamId: string
    cost: number
  }
  environment: {
    environmentId: string
    cost: number
  }
  component: {
    componentId: string
    cost: number
  }
}

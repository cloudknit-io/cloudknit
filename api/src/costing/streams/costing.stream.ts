import { Controller, Param, Sse } from '@nestjs/common'
import { from, Observable, Observer } from 'rxjs'
import { filter, map } from 'rxjs/operators'
import { Component } from 'src/typeorm/costing/entities/Component'
import { ComponentService } from '../services/component.service'
import { Mapper } from '../utilities/mapper'
import { MessageEvent } from '@nestjs/common';

@Controller({
  path: 'costing/stream/api/v1',
})
export class CostingStream {
  constructor(private readonly componentService: ComponentService) {}

  @Sse('all')
  componentCost(): Observable<MessageEvent> {
    this.componentService.getAll()
    return from(this.componentService.stream).pipe(
      map((components: Component[]) => ({
        data: Mapper.getStreamData(components),
      })),
    )
  }

  @Sse('team/:teamId')
  teamCost(@Param('teamId') teamId: string): Observable<MessageEvent> {
    this.componentService.getTeamCost(teamId)
    return from(this.componentService.stream).pipe(
      map((components: Component[]) => ({
        data: Mapper.getStreamData(components),
      })),
    )
  }

  @Sse('environment/:teamId/:environmentId')
  environmentCost(
    @Param('teamId') teamId: string,
    @Param('environmentId') environmentId: string,
  ): Observable<MessageEvent> {
    this.componentService.getEnvironmentCost(teamId, environmentId)
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

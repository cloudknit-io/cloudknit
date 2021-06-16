import { Controller, Param, Sse } from '@nestjs/common'
import { from, Observable } from 'rxjs'
import { map } from 'rxjs/operators'
import { Component } from 'src/typeorm/costing/entities/Component'
import { ComponentService } from '../services/component.service'
import { Mapper } from '../utilities/mapper'

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
    this.componentService.getTeamCost(teamId);
    return from(this.componentService.stream).pipe(
      map((components: Component[]) => ({
        data: Mapper.getStreamData(components),
      })),
    )
  }

  @Sse('environment/:teamId/:environmentId')
  environmentCost(@Param('teamId') teamId: string, @Param('environmentId') environmentId: string): Observable<MessageEvent> {
    this.componentService.getEnvironmentCost(teamId, environmentId);
    return from(this.componentService.stream).pipe(
      map((components: Component[]) => ({
        data: Mapper.getStreamData(components),
      })),
    )
  }
}

export interface MessageEvent {
  data: string | object
  id?: string
  type?: string
  retry?: number
}

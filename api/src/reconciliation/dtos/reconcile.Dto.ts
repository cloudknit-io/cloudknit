import { ApiProperty } from '@nestjs/swagger'

export class EvnironmentReconcileDto {
  @ApiProperty()
  reconcileId?: string;
  @ApiProperty()
  name: string
  @ApiProperty()
  teamName: string
  @ApiProperty()
  status: string
  @ApiProperty()
  startDateTime: string
  @ApiProperty()
  endDateTime?: string
  @ApiProperty({
    type: () => [ComponentReconcileDto],
  })
  componentReconciles?: ComponentReconcileDto[]
}

export class ComponentReconcileDto {
  @ApiProperty()
  id: string;
  @ApiProperty()
  name: string
  @ApiProperty()
  status: string
  @ApiProperty()
  startDateTime: string
  @ApiProperty()
  endDateTime?: string
}

import { ApiProperty } from '@nestjs/swagger'
import { CostResource } from './Resource.dto';

export class ComponentType {
  @ApiProperty()
  componentName: string;
  @ApiProperty()
  cost: number;
  @ApiProperty()
  resources: CostResource[]
  @ApiProperty()
  isDestroyed: boolean
  @ApiProperty()
  status: string;
  @ApiProperty()
  duration: number;
}

export class CostingDto {
  @ApiProperty()
  teamName: string
  @ApiProperty()
  environmentName: string
  @ApiProperty({
    type: ComponentType
  })
  component: ComponentType
}

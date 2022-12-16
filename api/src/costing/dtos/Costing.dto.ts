import { ApiProperty } from '@nestjs/swagger'
import { Resource } from 'src/typeorm/resources/Resource.entity';

export class ComponentType {
  @ApiProperty()
  componentName: string;
  @ApiProperty()
  cost: number;
  @ApiProperty()
  resources: Resource[]
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

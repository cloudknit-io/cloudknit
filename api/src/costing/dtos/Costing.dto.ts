import { ApiProperty } from '@nestjs/swagger'
import { Resource } from 'src/typeorm/resources/Resource.entity';

export class ComponentType {
  @ApiProperty()
  componentName: string;
  @ApiProperty()
  cost: number;
  @ApiProperty()
  resources: Resource[]
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




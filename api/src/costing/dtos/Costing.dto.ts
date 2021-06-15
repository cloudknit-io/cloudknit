import { ApiProperty } from '@nestjs/swagger'

export class ComponentType {
  @ApiProperty()
  componentName: string;
  @ApiProperty()
  cost: number;
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




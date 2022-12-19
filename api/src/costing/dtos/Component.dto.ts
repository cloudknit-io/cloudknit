import { ApiProperty } from '@nestjs/swagger';
import { EnvironmentDto } from './Environment.dto';
import { CostResource } from './Resource.dto';

export class ComponentDto {

  @ApiProperty({
    name: 'Component Name'
  })
  componentName: string = '';

  @ApiProperty({
    name: 'Cost'
  })
  estimatedCost: number = -1;

  @ApiProperty()
  id: string;

  @ApiProperty()
  environment?: EnvironmentDto;

  @ApiProperty()
  status: string;

  @ApiProperty()
  lastReconcileDatetime: string;

  @ApiProperty()
  duration: number;

  @ApiProperty()
  isDestroyed?: boolean

  @ApiProperty()
  costResources?: CostResource[];
}

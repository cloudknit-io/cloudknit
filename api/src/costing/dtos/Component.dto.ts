import { ApiProperty } from '@nestjs/swagger';
import { EnvironmentDto } from './Environment.dto';
import { CostResource } from './Resource.dto';

export class ComponentDto {

  @ApiProperty({
    name: 'Component Name'
  })
  name: string = '';

  @ApiProperty({
    name: 'Cost'
  })
  estimatedCost: number = -1;

  @ApiProperty()
  id: number;

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

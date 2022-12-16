import { ApiProperty } from '@nestjs/swagger';
import { EnvironmentDto } from './Environment.dto';

export class ComponentDto {

  @ApiProperty({
    name: 'Component Name'
  })
  componentName: string = '';

  @ApiProperty({
    name: 'Cost'
  })
  cost: number = -1;

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

  // TODO : Make ResourceDto
  @ApiProperty()
  resources?: object[];
}

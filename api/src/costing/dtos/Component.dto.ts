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
}

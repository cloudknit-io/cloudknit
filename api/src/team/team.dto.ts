import { ApiProperty, ApiQuery } from '@nestjs/swagger';
import { IsIn, IsOptional, IsString } from 'class-validator';

export class TeamQueryParams {
  @ApiProperty({ required: false, type: 'boolean' })
  @IsString()
  @IsIn(['true', 'false'])
  @IsOptional()
  withCost: string = 'false';

  @ApiProperty({ required: false, type: 'boolean' })
  @IsString()
  @IsIn(['true', 'false'])
  @IsOptional()
  withEnvironments: string = 'false';

  @ApiProperty({ required: false, type: 'boolean' })
  @IsString()
  @IsIn(['true', 'false'])
  @IsOptional()
  withComponents: string = 'false';
}

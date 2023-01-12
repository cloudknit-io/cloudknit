import { ApiProperty, ApiQuery } from "@nestjs/swagger";
import { IsIn, IsOptional, IsString } from "class-validator";

export class TeamQueryParams {
  @ApiProperty({required: false})
  @IsString()
  @IsIn(['true', 'false'])
  @IsOptional()
  withCost: string = 'false';

  @ApiProperty({required: false})
  @IsString()
  @IsIn(['true', 'false'])
  @IsOptional()
  withEnvironments: string = 'false';
}

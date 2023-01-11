import { IsBoolean, IsIn, IsNotEmpty, IsOptional, IsString } from "class-validator";

export class TeamQueryParams {
  @IsString()
  @IsIn(['true', 'false'])
  @IsOptional()
  withCost: string = 'false';

  @IsString()
  @IsIn(['true', 'false'])
  @IsOptional()
  withEnvironments: string = 'false';
}

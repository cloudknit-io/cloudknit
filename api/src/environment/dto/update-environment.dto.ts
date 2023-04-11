import { ApiProperty, ApiPropertyOptional, PartialType } from '@nestjs/swagger';
import { Transform } from 'class-transformer';
import { IsBoolean } from 'class-validator';
import { EnvironmentReconcile } from 'src/typeorm';
import { CreateEnvironmentDto } from './create-environment.dto';

export class UpdateEnvironmentDto extends PartialType(CreateEnvironmentDto) {
  isDeleted?: boolean;
  latestEnvRecon?: EnvironmentReconcile;
  lastReconcileDatetime?: string;
}


export class PatchEnvQueryParams {
  @ApiProperty()
  @IsBoolean()
  @Transform(({ value} ) => value === 'true')
  reconcile: boolean;
}

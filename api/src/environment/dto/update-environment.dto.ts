import { ApiProperty, ApiPropertyOptional, PartialType } from '@nestjs/swagger';
import { Transform } from 'class-transformer';
import { IsBoolean, IsOptional } from 'class-validator';
import { EnvironmentReconcile } from 'src/typeorm';
import { CreateEnvironmentDto } from './create-environment.dto';

export class UpdateEnvironmentDto extends PartialType(CreateEnvironmentDto) {
  isDeleted?: boolean;
  latestEnvRecon?: EnvironmentReconcile;
  lastReconcileDatetime?: string;
  @ApiPropertyOptional()
  @IsOptional()
  @IsBoolean()
  isReconcile?: boolean;
}

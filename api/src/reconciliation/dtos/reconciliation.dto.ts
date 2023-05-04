import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';
import { Type } from 'class-transformer';
import {
  IsArray,
  IsBoolean,
  IsDateString,
  IsNotEmpty,
  IsNumber,
  IsOptional,
  IsString,
  ValidateNested,
} from 'class-validator';
import { CostResource } from 'src/component/dto/update-component.dto';
import { EnvSpecComponentDto } from 'src/environment/dto/env-spec.dto';

export class CreateEnvironmentReconciliationBaseDto {
  @ApiProperty({ type: EnvSpecComponentDto, isArray: true })
  @IsNotEmpty()
  @ValidateNested({ each: true })
  @Type(() => EnvSpecComponentDto)
  components: EnvSpecComponentDto[];

  @ApiProperty()
  @IsArray()
  errorMessage?: string[];
}

export class CreateEnvironmentReconciliationDto extends CreateEnvironmentReconciliationBaseDto {
  @ApiProperty()
  @IsString()
  @IsNotEmpty()
  name: string;

  @ApiProperty()
  @IsString()
  @IsNotEmpty()
  teamName: string;

  @ApiProperty()
  @IsDateString()
  @IsNotEmpty()
  startDateTime: string;
}

export class CreateErrorEnvironmentReconciliationDto extends CreateEnvironmentReconciliationBaseDto {
  @ApiProperty()
  @IsString()
  @IsNotEmpty()
  name: string;

  @ApiProperty()
  @IsString()
  @IsNotEmpty()
  teamName: string;

  @ApiProperty()
  @IsDateString()
  @IsNotEmpty()
  startDateTime: string;

  @ApiProperty()
  @IsDateString()
  @IsNotEmpty()
  endDateTime: string;

  @ApiProperty()
  @IsString()
  @IsNotEmpty()
  status: string;

  @ApiProperty()
  @IsNumber()
  @IsOptional()
  estimatedCost: number;
}

export class UpdateEnvironmentReconciliationDto {
  @ApiProperty()
  @IsString()
  @IsNotEmpty()
  status: string;

  @ApiProperty()
  @IsDateString()
  @IsOptional()
  endDateTime?: string;

  @ApiProperty()
  @IsNumber()
  @IsOptional()
  estimatedCost?: number;
}

export class CreateComponentReconciliationDto {
  @ApiProperty()
  @IsNumber()
  @IsNotEmpty()
  envReconcileId: number;

  @ApiProperty()
  @IsString()
  @IsNotEmpty()
  name: string;

  @ApiPropertyOptional()
  @IsDateString()
  startDateTime?: string;

  @ApiPropertyOptional()
  @IsString()
  status?: string;
}

export class UpdateComponentReconciliationDto {
  @ApiProperty()
  @IsString()
  @IsOptional()
  status?: string;

  @ApiProperty()
  @IsDateString()
  @IsOptional()
  endDateTime?: string;

  @ApiProperty()
  @IsString()
  @IsOptional()
  approvedBy?: string;

  @ApiPropertyOptional()
  @IsOptional()
  @IsNumber({
    maxDecimalPlaces: 0,
  })
  duration?: number;

  @ApiPropertyOptional()
  @IsOptional()
  @IsString()
  lastWorkflowRunId?: string;

  @ApiPropertyOptional()
  @IsOptional()
  @IsNumber({
    maxDecimalPlaces: 5,
  })
  estimatedCost?: number;

  @ApiPropertyOptional({ type: [CostResource] })
  @IsOptional()
  @ValidateNested()
  @Type(() => CostResource)
  costResources?: CostResource[];

  @ApiPropertyOptional()
  @IsOptional()
  @IsBoolean()
  isDestroyed?: boolean;

  @ApiPropertyOptional()
  @IsOptional()
  @IsBoolean()
  isSkipped?: boolean;
}

export class CreatedEnvironmentReconcile {
  @ApiProperty()
  @IsNumber()
  @IsNotEmpty()
  reconcileId: number;
}

export class CreatedComponentReconcile {
  @ApiProperty()
  @IsNumber()
  @IsNotEmpty()
  reconcileId: number;
}

export class RespGetEnvReconStatus {
  @ApiProperty()
  @IsString()
  @IsNotEmpty()
  status: string;
}

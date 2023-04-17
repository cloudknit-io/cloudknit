import { ApiProperty } from '@nestjs/swagger';
import { Type } from 'class-transformer';
import {
  IsArray,
  IsDateString,
  IsNotEmpty,
  IsNumber,
  IsOptional,
  IsString,
  ValidateNested,
} from 'class-validator';
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

export class UpdateEnvironmentReconciliationDto {
  @ApiProperty()
  @IsString()
  @IsNotEmpty()
  status: string;

  @ApiProperty()
  @IsDateString()
  @IsOptional()
  endDateTime?: string;
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

  @ApiProperty()
  @IsDateString()
  @IsNotEmpty()
  startDateTime: string;
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

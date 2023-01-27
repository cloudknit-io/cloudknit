import { ApiProperty } from '@nestjs/swagger';
import {
  IsDateString,
  IsNotEmpty,
  IsNumber,
  IsOptional,
  IsString,
} from 'class-validator';

export class CreateEnvironmentReconciliationDto {
  @ApiProperty()
  @IsString()
  @IsNotEmpty()
  name: string;

  @ApiProperty()
  @IsString()
  @IsNotEmpty()
  teamName: string;

  @ApiProperty()
  @IsString()
  @IsNotEmpty()
  gitSha: string;

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

export class RespGetEnvReconStatus {
  @ApiProperty()
  @IsString()
  @IsNotEmpty()
  status: string;
}

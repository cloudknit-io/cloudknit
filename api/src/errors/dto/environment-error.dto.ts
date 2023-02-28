import { ApiProperty } from '@nestjs/swagger';
import { IsArray, IsNotEmpty, IsNumber, IsString } from 'class-validator';
import { CreateEnvironmentDto } from 'src/environment/dto/create-environment.dto';
import { EnvSpecDto } from 'src/environment/dto/env-spec.dto';

export enum ErrorType {
  VALIDATION_ERROR = 1
}

export class EnvironmentErrorDto extends CreateEnvironmentDto {
  errorType: ErrorType;
  errorMessage: string[];
}

export class EnvironmentErrorSpecDto extends EnvSpecDto {
  @ApiProperty()
  @IsNotEmpty()
  @IsNumber()
  errorType: ErrorType;

  @ApiProperty()
  @IsNotEmpty()
  @IsArray()
  errorMessage: string[];

  @ApiProperty()
  @IsNotEmpty()
  @IsString()
  status: string;
}

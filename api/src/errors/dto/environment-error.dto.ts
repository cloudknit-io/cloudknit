import { ApiProperty } from '@nestjs/swagger';
import { IsArray, IsNotEmpty, IsString } from 'class-validator';
import { CreateEnvironmentDto } from 'src/environment/dto/create-environment.dto';
import { EnvSpecDto } from 'src/environment/dto/env-spec.dto';

export class EnvironmentErrorDto extends CreateEnvironmentDto {
  errorMessage: string[];
}

export class EnvironmentErrorSpecDto {
  @ApiProperty()
  @IsNotEmpty()
  @IsString()
  envName: string;

  @ApiProperty()
  @IsNotEmpty()
  @IsArray()
  errorMessage: string[];

  @ApiProperty()
  @IsNotEmpty()
  @IsString()
  status: string;
}

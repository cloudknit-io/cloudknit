import { ApiProperty } from '@nestjs/swagger';
import { IsNotEmpty, IsString } from 'class-validator';
import { EnvironmentReconcile } from 'src/typeorm';

export interface EnvironmentReconcileWrap extends EnvironmentReconcile {
  duration: number;
}

export class GetEnvReconStatusQueryParams {
  @ApiProperty()
  @IsString()
  @IsNotEmpty()
  teamName: string;

  @ApiProperty()
  @IsString()
  @IsNotEmpty()
  envName: string;
}

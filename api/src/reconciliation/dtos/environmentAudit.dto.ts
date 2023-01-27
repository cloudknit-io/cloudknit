import { ApiProperty } from '@nestjs/swagger';
import { IsString } from 'class-validator';
import { EnvironmentReconcile } from 'src/typeorm';

export interface EnvironmentReconcileWrap extends EnvironmentReconcile {
  duration: number;
}

export class EnvShaParams {
  @ApiProperty({ required: true, type: 'string' })
  @IsString()
  sha: string;
}
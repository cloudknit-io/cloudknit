import { ApiProperty } from '@nestjs/swagger';
import { IsNotEmpty, IsString } from 'class-validator';
import { ComponentReconcile } from 'src/typeorm';

export interface ComponentReconcileWrap extends ComponentReconcile {
  duration: number;
}

export class ApprovedByDto {
  @ApiProperty()
  @IsString()
  @IsNotEmpty()
  email: string;
}

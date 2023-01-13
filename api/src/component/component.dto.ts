import { ApiProperty } from '@nestjs/swagger';
import { IsIn, IsOptional, IsString } from 'class-validator';
import { Component } from 'src/typeorm';

export class ComponentQueryParams {
  @ApiProperty({ required: false, type: 'boolean' })
  @IsString()
  @IsIn(['true', 'false'])
  @IsOptional()
  withLastAuditStatus: string = 'false';
}

export class ComponentWrap extends Component {
  lastAuditStatus?: string;
}

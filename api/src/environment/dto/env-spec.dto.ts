import { ApiProperty } from '@nestjs/swagger';
import { Type } from 'class-transformer';
import {
  IsNotEmpty,
  IsOptional,
  IsString,
  ValidateNested,
} from 'class-validator';

export class EnvSpecComponentDto {
  @ApiProperty()
  @IsNotEmpty()
  @IsString()
  name: string;

  @ApiProperty()
  @IsNotEmpty()
  @IsString()
  type: string;

  @ApiProperty()
  @IsOptional()
  dependsOn: string[];
}

export class EnvSpecDto {
  @ApiProperty()
  @IsNotEmpty()
  @IsString()
  envName: string;

  @ApiProperty({ type: EnvSpecComponentDto, isArray: true })
  @IsNotEmpty()
  @ValidateNested({ each: true })
  @Type(() => EnvSpecComponentDto)
  components: EnvSpecComponentDto[];
}

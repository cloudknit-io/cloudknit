import { Type } from 'class-transformer';
import {
  IsNotEmpty,
  IsOptional,
  IsString,
  ValidateNested,
} from 'class-validator';

export class EnvSpecDto {
  @IsNotEmpty()
  @IsString()
  envName: string;

  @IsNotEmpty()
  @ValidateNested({ each: true })
  @Type(() => EnvSpecComponentDto)
  components: EnvSpecComponentDto[];
}

export class EnvSpecComponentDto {
  @IsNotEmpty()
  @IsString()
  name: string;

  @IsNotEmpty()
  @IsString()
  type: string;

  @IsOptional()
  dependsOn: string[];
}

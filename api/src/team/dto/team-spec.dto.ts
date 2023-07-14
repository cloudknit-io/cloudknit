import { ApiProperty } from '@nestjs/swagger';
import { Type } from 'class-transformer';
import { IsBoolean, IsNotEmpty, IsString, ValidateNested } from 'class-validator';

export class TeamConfigRepoDto {
  @ApiProperty({ required: true })
  @IsString()
  @IsNotEmpty()
  source: string;

  @ApiProperty({ required: true })
  @IsString()
  @IsNotEmpty()
  path: string;
}

export class TeamSpecDto {
  @ApiProperty({ required: true })
  @IsNotEmpty()
  @IsString()
  teamName: string;

  @ApiProperty({ required: false, type: 'boolean' })
  teardownProtection?: boolean;

  @ApiProperty({ required: true })
  @IsNotEmpty()
  @ValidateNested()
  @Type(() => TeamConfigRepoDto)
  configRepo: TeamConfigRepoDto;
}

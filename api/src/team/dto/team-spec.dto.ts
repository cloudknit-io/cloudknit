import { Type } from "class-transformer"
import { IsNotEmpty, IsString, ValidateNested } from "class-validator"

export class TeamConfigRepoDto {
  @IsString()
  @IsNotEmpty()
  source: string

  @IsString()
  @IsNotEmpty()
  path: string
}

export class TeamSpecDto {
  @IsNotEmpty()
  teamName: string

  @IsNotEmpty()
  @ValidateNested()
  @Type(() => TeamConfigRepoDto)
  configRepo: TeamConfigRepoDto
}

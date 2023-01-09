import { IsString } from "class-validator";


export class TeamEnvCompQueryParams {
  @IsString()
  teamName: string

  @IsString()
  envName: string

  @IsString()
  compName: string
}

export class TeamEnvQueryParams {
  @IsString()
  teamName: string

  @IsString()
  envName: string
}

import { ValidationPipe } from "@nestjs/common";
import { IsString, ValidatorOptions } from "class-validator";

export class RequiredQueryValidationPipe extends ValidationPipe {
  protected validatorOptions: ValidatorOptions;
  constructor() {
    super();
    
    this.validatorOptions = {
      forbidUnknownValues: true,
      skipMissingProperties: false
    };
  }
}

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

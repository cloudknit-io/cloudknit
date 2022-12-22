import { ValidationPipe } from "@nestjs/common";
import { ValidatorOptions } from "class-validator";

export class RequiredQueryValidationPipe extends ValidationPipe {
  protected validatorOptions: ValidatorOptions;
  // ValidationPipe({forbidUnknownValues: true, skipMissingProperties: false})
  constructor() {
    super();
    
    this.validatorOptions = {
      forbidUnknownValues: true,
      skipMissingProperties: false
    };
  }
}

export class TeamEnvCompQueryParams {
  teamName: string
  envName: string
  compName: string
}

export class TeamEnvQueryParams {
  teamName: string
  envName: string
}

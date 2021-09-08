import { Body, Controller, Get, Param, Post } from "@nestjs/common";
import { AwsSecretDto } from "./dtos/aws-secret.dto";
import { SecretsService } from "./secrets.service";

@Controller("secrets")
export class SecretsController {
  constructor(private readonly secretsService: SecretsService) {}

  @Post("update/aws-secret")
  public async updateAwsSecret(
    @Body() awsSecrets: AwsSecretDto
  ) {
    return await this.secretsService.putSsmSecrets(awsSecrets);
  }

  @Post("exists/aws-secret")
  public async secretsExist(@Body() pathNames: string[]) {
    return await this.secretsService.ssmSecretsExists(pathNames);
  }
}

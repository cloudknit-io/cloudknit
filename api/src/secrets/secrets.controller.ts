import { Body, Controller, Delete, Get, Param, Post, Query } from "@nestjs/common";
import { AwsSecretDto } from "./dtos/aws-secret.dto";
import { SecretsService } from "./secrets.service";

@Controller("secrets")
export class SecretsController {
  constructor(private readonly secretsService: SecretsService) {}

  @Post("update/aws-secret")
  public async updateAwsSecret(@Body() req: any) {
    const { awsSecrets } = req;
    return await this.secretsService.putSsmSecrets(awsSecrets);
  }

  @Post("exists/aws-secret")
  public async secretsExist(@Body() req: any) {
    const { pathNames } = req;
    return await this.secretsService.ssmSecretsExists(pathNames);
  }

  @Post("get/ssm-secrets")
  public async getSSMSecrets(@Body() req: any) {
    const { path } = req;
    return await this.secretsService.getSsmSecretsByPath(path);
  }

  @Post("get/environments")
  public async getEnvironments(@Body() req: any) {
    const { path } = req;
    return await this.secretsService.getEnvironments(path);
  }

  @Delete("delete/ssm-secret")
  public async deleteSSMParameter(@Query('path') path: string) {
    return await this.secretsService.deleteSSMSecret(path);
  }
}

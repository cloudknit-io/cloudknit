import { Body, Controller, Get, Param, Post } from "@nestjs/common";
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
    console.log(req);
    const { path, recursive } = req;
    return await this.secretsService.getSsmSecretsByPath(
      path,
      recursive === "true"
    );
  }
}

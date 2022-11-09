import { Body, Controller, Delete, Get, Logger, Param, Post, Query, Request } from "@nestjs/common";
import { SecretsService } from "./secrets.service";

@Controller({
  version: '1'
})
export class SecretsController {
  private readonly logger = new Logger(SecretsController.name);

  constructor(private readonly secretsService: SecretsService) {}

  @Post("update/aws-secret")
  public async updateAwsSecret(@Request() req, @Body() body: any) {
    const { awsSecrets } = body;
    return await this.secretsService.putSsmSecrets(req.org, awsSecrets);
  }

  @Post("exists/aws-secret")
  public async secretsExist(@Request() req, @Body() body: any) {
    const { pathNames } = body;
    return await this.secretsService.ssmSecretsExists(req.org, pathNames);
  }

  @Post("get/ssm-secret")
  public async getSSMSecret(@Request() req, @Body() body: any) {
    const { path } = body;
    const value = await this.secretsService.getSsmSecret(req.org, path);

    return { data: value };
  }

  @Post("get/ssm-secrets")
  public async getSSMSecrets(@Request() req, @Body() body: any) {
    const { path } = body;
    return await this.secretsService.getSsmSecretsByPath(req.org, path);
  }

  @Post("environments")
  public async getEnvironments(@Request() req, @Body() body: any) {
    const { path } = body;
    return await this.secretsService.getEnvironments(req.org, path);
  }

  @Delete("delete/ssm-secret")
  public async deleteSSMParameter(@Request() req, @Query('path') path: string) {
    this.logger.debug("Deleting secret", { orgName: req.org.name, path });
    return await this.secretsService.deleteSSMSecret(req.org, path);
  }
}

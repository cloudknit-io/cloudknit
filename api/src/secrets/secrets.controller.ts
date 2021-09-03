import { Controller, Get, Param } from "@nestjs/common";
import { SecretsService } from "./secrets.service";

@Controller("secrets")
export class SecretsController {
  constructor(private readonly secretsService: SecretsService) {}

  @Get("update/aws-secret/:accessKeyId/:secretAccessKey")
  public async updateAwsSecret(
    @Param("accessKeyId") accessKeyId: string,
    @Param("secretAccessKey") secretAccessKey: string
  ) {
    return await this.secretsService.createOrUpdateSecret(
      accessKeyId,
      secretAccessKey
    );
  }

  @Get("exists/aws-secret")
  public async secretExist() {
      return await this.secretsService.secretExist();
  }
}

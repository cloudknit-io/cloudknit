import { Body, Controller, Get, Param, Patch, Post, Req, Request } from "@nestjs/common";
import { OrganizationService } from "./organization.service";

@Controller({
  version: '1'
})
export class OrganizationController {

  constructor(
      private readonly orgService: OrganizationService
  ){}

  @Post('oath/credentials') 
  public async saveOAuthCredentials(@Body() payload) {
      return await this.orgService.saveOAuthCredentials(payload);
  }

  @Patch('oath/credentials') 
  public async patchOAuthCredentials(@Body() payload) {
      return await this.orgService.patchOrganizationData(payload);
  }

  @Post('github/credentials') 
  public async saveGitHubCredentials(@Body() payload) {
      return await this.orgService.saveGitHubCredentials(payload);
  }

  @Patch('github/credentials')
  public async patchGitHubCredentials(@Body() payload) {
    return await this.orgService.patchCRD(payload);
  }

  @Get()
  public async getOrg(@Request() req) {
    return await req.org;
  }
}

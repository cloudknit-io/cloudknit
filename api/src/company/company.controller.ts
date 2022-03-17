import { Body, Controller, Patch, Post } from "@nestjs/common";
import { CompanyService } from "./company.service";

@Controller("company")
export class CompanyController {

  constructor(
      private readonly initService: CompanyService
  ){}

  @Post('oath/credentials') 
  public async saveOAuthCredentials(@Body() payload) {
      return await this.initService.saveOAuthCredentials(payload);
  }

  @Patch('oath/credentials') 
  public async patchOAuthCredentials(@Body() payload) {
      return await this.initService.patchOrganisationData(payload);
  }
}

import { Controller, Get, Request } from "@nestjs/common";
import { OrganizationService } from "./organization.service";

@Controller({
  version: '1'
})
export class OrganizationController {

  constructor(
      private readonly orgService: OrganizationService
  ){}

  @Get()
  public async getOrg(@Request() req) {
    return await req.org;
  }
}

import { Body, Controller, Get, Param, Patch, Request } from "@nestjs/common";
import { PatchOrganizationDto } from "./organization.dto";
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

  @Patch()
  public async patchOrganization(@Body() payload: PatchOrganizationDto, @Request() req) {
    return await this.orgService.patchOrganization(req.org, payload);
  }
}

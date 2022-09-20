import { Body, Controller, Get, Param, Post, Request } from "@nestjs/common";
import { CreateOrganizationDto } from "./rootOrganization.dto";
import { RootOrganizationsService } from "./rootOrganization.service";

@Controller({
  version: '1'
})
export class RootOrganizationsController {

  constructor(
      private readonly orgService: RootOrganizationsService
  ){}

  // @Get()
  // public async getAll() {
  //   return await this.orgService.getOrganizations();
  // }

  @Post()
  public async create(@Body() body: CreateOrganizationDto) {
    return await this.orgService.create(body);
  }
}

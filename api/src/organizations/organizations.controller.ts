import { Body, Controller, Get, Param, Post, Request } from "@nestjs/common";
import { CreateOrganizationDto } from "./Organization.dto";
import { OrganizationsService } from "./organizations.service";

@Controller({
  version: '1'
})
export class OrganizationsController {

  constructor(
      private readonly orgService: OrganizationsService
  ){}

  @Get()
  public async getAll() {
    return await this.orgService.getOrganizations();
  }

  @Post()
  public async create(@Body() body: CreateOrganizationDto) {
    return await this.orgService.create(body);
  }
}

import { BadRequestException, Body, Controller, Get, Logger, Post, Query, Request } from "@nestjs/common";
import { CreateOrganizationDto } from "./root.organization.dto";
import { RootOrganizationsService } from "./root.organization.service";
import { handleSqlErrors } from 'src/utilities/errorHandler';

@Controller({
  version: '1'
})
export class RootOrganizationsController {
  private readonly logger = new Logger(RootOrganizationsController.name);

  constructor(
      private readonly orgService: RootOrganizationsService
  ){}

  private OrganizationNameRegex = /^[a-zA-Z]+[a-zA-Z0-9]*(-[a-zA-Z0-9]+)*$/;

  @Get()
  public async getAll(@Query('github-org-name') ghOrgName: string) {
    if (!ghOrgName) {
      throw new BadRequestException();
    }

    return await this.orgService.getOrgByGithubOrg(ghOrgName);
  }

  @Post()
  public async create(@Request() req, @Body() body: CreateOrganizationDto) {
    if (!body.name || !this.OrganizationNameRegex.test(body.name) || body.name.length > 63) {
      throw new BadRequestException("Organization name is invalid");
    }

    try {
      return await this.orgService.create(body);
    } catch (error) {
      handleSqlErrors(error, "organization already exists");

      this.logger.error({ message: 'error creating organization', body }, error.stack);

      throw error;
    }
  }
}

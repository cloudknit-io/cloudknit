import { BadRequestException, Body, Controller, Get, InternalServerErrorException, Logger, Param, Post, Query, Request } from "@nestjs/common";
import { CreateOrganizationDto } from "./rootOrganization.dto";
import { RootOrganizationsService } from "./rootOrganization.service";

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
  public async create(@Body() body: CreateOrganizationDto) {
    if (!this.OrganizationNameRegex.test(body.name) || body.name.length > 63) {
      throw new BadRequestException("Organization name is invalid");
    }

    try {
      return await this.orgService.create(body);
    } catch (error) {
      if (error.message.startsWith('Duplicate entry')) {
        throw new BadRequestException('Organization name already exists');
      }

      this.logger.error({ message: 'error creating organization', body }, error.stack);

      throw error;
    }
  }
}

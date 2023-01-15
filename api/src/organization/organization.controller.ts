import {
  BadRequestException,
  Body,
  Controller,
  Get,
  Logger,
  Param,
  Patch,
  Post,
  Query,
  Request,
} from '@nestjs/common';
import { APIRequest } from 'src/types';
import { handleSqlErrors } from 'src/utilities/errorHandler';
import {
  CreateOrganizationDto,
  PatchOrganizationDto,
} from './organization.dto';
import { OrganizationService } from './organization.service';

@Controller({
  version: '1',
})
export class OrganizationController {
  private readonly logger = new Logger(OrganizationController.name);

  constructor(private readonly orgService: OrganizationService) {}

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
    if (
      !body.name ||
      !this.OrganizationNameRegex.test(body.name) ||
      body.name.length > 63
    ) {
      throw new BadRequestException('Organization name is invalid');
    }

    try {
      return await this.orgService.create(body);
    } catch (error) {
      handleSqlErrors(error, 'organization already exists');

      this.logger.error(
        { message: 'error creating organization', body },
        error.stack
      );

      throw error;
    }
  }

  @Get('/:orgId')
  public async getOrg(
    @Request() req: APIRequest,
    @Param('orgId') forOpenApi: string
  ) {
    return req.org;
  }

  @Patch('/:orgId')
  public async patchOrganization(
    @Request() req: APIRequest,
    @Body() payload: PatchOrganizationDto,
    @Param('orgId') forOpenApi: string
  ) {
    return await this.orgService.patchOrganization(req.org, payload);
  }
}

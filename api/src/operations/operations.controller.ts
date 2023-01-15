import { Controller, Get, Request } from '@nestjs/common';
import { OrganizationService } from 'src/organization/organization.service';
import { OrgApiParam } from 'src/types';
import { OperationsService } from './operations.service';

@Controller({
  version: '1',
})
export class OperationsController {
  constructor(
    private readonly opsService: OperationsService,
    private readonly orgService: OrganizationService
  ) {}

  @Get('/is-provisioned')
  @OrgApiParam()
  public async check(@Request() req) {
    const isProvisioned = await this.opsService.isOrgProvisioned(req.org);

    if (isProvisioned === true) {
      const org = await this.orgService.patchOrganization(req.org, {
        provisioned: true,
      });
      return org;
    }

    return req.org;
  }
}

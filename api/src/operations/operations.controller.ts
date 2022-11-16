import { Controller, Get, Request } from '@nestjs/common';
import { OperationsService } from './operations.service';

@Controller({
  version: '1'
})
export class OperationsController {
  constructor(
    private readonly opsService: OperationsService
  ){}

  @Get("/is-provisioned")
  public async check(@Request() req) {
    const isProvisioned = await this.opsService.isOrgProvisioned(req.org);

    return { provisioned: isProvisioned };
  }
}

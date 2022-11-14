import { HttpException, HttpStatus, Injectable, Logger, NestMiddleware } from '@nestjs/common';
import { InjectRepository } from "@nestjs/typeorm";
import { Repository } from "typeorm";
import { Organization } from "src/typeorm";
import { Response, NextFunction } from 'express';
import { APIRequest } from 'src/types';

@Injectable()
export class OrganizationMiddleware implements NestMiddleware {
  private readonly logger = new Logger(OrganizationMiddleware.name);

  constructor(
    @InjectRepository(Organization) private readonly orgRepo: Repository<Organization>
  ) {}

  async use(req: APIRequest, res: Response, next: NextFunction) {
    let org, id;
    const orgId = req.params.orgId;

    try {
      id = parseInt(req.params.orgId, 10);
    } catch (e) {}
    
    if (isNaN(id)) {
      try {
        org = await this.orgRepo.findOne({
          where: { name: orgId }
        });
      } catch (e) {
        this.logger.error('could not get org by name', e.message, orgId)
      }
    } else {
      try {
        org = await this.orgRepo.findOne({
          where: { id }
        });
      } catch (e) {
        this.logger.error('could not get org by number', e.message, orgId);
      }
    }

    if (!org) {
      this.logger.error(`Bad orgId: ${orgId}`);
      throw new HttpException('Forbidden', HttpStatus.FORBIDDEN);
    }

    req.org = org;

    next();
  }
}

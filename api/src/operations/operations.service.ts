import { Injectable, Logger } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { S3Handler } from 'src/utilities/s3Handler';
import { Organization } from 'src/typeorm';
import { Repository } from 'typeorm';
import { get } from 'src/config';

@Injectable()
export class OperationsService {
  private readonly logger = new Logger(OperationsService.name);
  private readonly s3h = S3Handler.instance();
  private readonly config = get();
  private readonly ckEnvironment = this.config.environment;

  constructor(
    @InjectRepository(Organization)
    private OrganizationRepo: Repository<Organization>
  ) {}

  async isOrgProvisioned(org: Organization) {
    // gets provisioned object in s3
    const resp = await this.s3h.getObjectList(
      `cloudknit-${this.ckEnvironment}-system`,
      `provisioned-orgs/${org.name}`
    );

    // we should loop through all entries in case there are similarly named companies
    //
    // Assume we have two companies. One is named Apple and the other is named
    // AppleTechnologies. If the current company is Apple then the above query (prefix = Apple)
    // could return both entries since the underlying S3 library does fuzzy matching.
    for (const entry of resp.entries()) {
      const parts = entry[1].Key.split('/');

      this.logger.log('Bucket entry', { parts });

      if (parts[1] === org.name) {
        return true;
      }
    }

    return false;
  }
}

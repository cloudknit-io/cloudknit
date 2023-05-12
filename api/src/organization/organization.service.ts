import { BadRequestException, Injectable, Logger } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { Organization, User } from 'src/typeorm';
import {
  CreateOrganizationDto,
  PatchOrganizationDto,
} from './organization.dto';
import { patchCompany } from 'src/k8s/patch-company';
import { getGithubOrgFromRepoUrl } from './utilities';
import { UsersService } from 'src/users/users.service';
import { SubmitProvisionOrg } from 'src/argowf/api';
import { get } from 'src/config';

@Injectable()
export class OrganizationService {
  private readonly logger = new Logger(OrganizationService.name);

  constructor(
    @InjectRepository(Organization) private orgRepo: Repository<Organization>,
    @InjectRepository(User) private userRepo: Repository<User>,
    private readonly usersService: UsersService
  ) {}

  async getOrganizations() {
    return this.orgRepo.find();
  }

  async getOrgByGithubOrg(ghOrgName: string) {
    return this.orgRepo.findOne({
      where: {
        githubOrgName: ghOrgName,
      },
    });
  }

  async create(newOrg: CreateOrganizationDto) {
    let user;

    if (newOrg.termsAgreedUserId) {
      user = await this.usersService.getUserById(newOrg.termsAgreedUserId);

      if (!user) {
        throw new BadRequestException('could not find user');
      }
    }

    if (newOrg.githubRepo) {
      const orgName = getGithubOrgFromRepoUrl(newOrg.githubRepo);

      if (!orgName) {
        throw new BadRequestException('bad github repo url');
      }

      newOrg.githubOrgName = orgName;
    }

    const org = await this.orgRepo.save({
      name: newOrg.name.toLowerCase(),
      githubRepo: newOrg.githubRepo,
      githubOrgName: newOrg.githubOrgName,
      termsAgreedUserId: user ? user.id : null,
    });

    this.logger.log({ message: 'created organization', org, user });

    if (user) {
      user.organizations = [...user.organizations, org];
      user = await this.userRepo.save(user);
    }

    if (get().argo.wf.skipProvision) {
      this.logger.log('Skipping provision workflow');
      return org;
    }

    try {
      await SubmitProvisionOrg({ orgName: org.name, orgId: org.id });
    } catch (error) {
      this.logger.error({
        message: `could not submit provision-org workflow for ${org.name}`,
        error: error.message,
      });
    }

    return org;
  }

  async getOrganization(id: number) {
    return this.orgRepo.findOne({
      where: {
        id,
      },
    });
  }

  async patchOrganization(org: Organization, payload: PatchOrganizationDto) {
    const reconcileProps = ['githubRepo', 'provisioned'];
    let changed = false;
    let updates = new Organization();

    for (const prop of reconcileProps) {
      if (payload.hasOwnProperty(prop)) {
        changed = true;
        this.logger.log(`patching org ${prop} with value ${payload[prop]}`);
        updates[prop] = payload[prop];
      }
    }

    if (!changed) {
      throw new BadRequestException('no values to update');
    }

    if (payload.githubRepo) {
      const orgName = getGithubOrgFromRepoUrl(payload.githubRepo);

      if (!orgName) {
        throw new BadRequestException('bad github repo url');
      }

      updates.githubOrgName = orgName;
    }

    await this.orgRepo
      .createQueryBuilder()
      .update(Organization)
      .set(updates)
      .where('id = :id', { id: org.id })
      .execute();

    if (payload.githubRepo) {
      try {
        await patchCompany(org, payload.githubRepo);
      } catch (error) {
        this.logger.log('Company CR was not updated', { org });
      }
    }

    this.logger.log({ message: 'updated org', org, updates });

    return await this.orgRepo.findOneBy({
      id: org.id,
    });
  }

  async getEmptyOrg() {
    const orgs = await this.getOrganizations();
    return orgs.find(org => org.users.length === 0); 
  }
}

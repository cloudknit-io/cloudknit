import { BadRequestException, Injectable, Logger } from "@nestjs/common";
import { InjectRepository } from "@nestjs/typeorm";
import { Repository } from "typeorm";
import { Organization, User } from "src/typeorm";
import { get } from "src/config";
import { CreateOrganizationDto } from "./root.organization.dto";
import { UsersService } from "src/users/users.service";
import { SubmitProvisionOrg } from "src/argowf/api";
import { getGithubOrgFromRepoUrl } from "src/organization/utilities";

@Injectable()
export class RootOrganizationsService {
  private readonly logger = new Logger(RootOrganizationsService.name);

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
        githubOrgName: ghOrgName
      }
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
      termsAgreedUserId: user ? user.id : null
    });

    this.logger.log({ message: 'created organization', org, user })

    if (user) {
      user.organizations = [ ...user.organizations, org ];
      user = await this.userRepo.save(user);
    }

    if (get().argo.wf.skipProvision) {
      this.logger.log('Skipping provision workflow');
      return org;
    }

    try {
      await SubmitProvisionOrg({ orgName: org.name, orgId: org.id });
    } catch (error) {
      this.logger.error({message: `could not submit provision-org workflow for ${org.name}`, error: error.message });
    }

    return org;
  }
}

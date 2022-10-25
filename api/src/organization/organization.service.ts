import { BadRequestException, Injectable } from "@nestjs/common";
import { InjectRepository } from "@nestjs/typeorm";
import { Repository } from "typeorm";
import { Organization } from "src/typeorm";
import { PatchOrganizationDto } from "./organization.dto";

@Injectable()
export class OrganizationService {

  constructor(@InjectRepository(Organization) private orgRepo: Repository<Organization>) { }

  async getOrganization(id: number) {
    return this.orgRepo.findOne({
      where: {
        id
      }
    });
  }

  async patchOrganization(org: Organization, payload: PatchOrganizationDto) {
    const { githubRepo } = payload;
    if (!githubRepo) {
      throw new BadRequestException('payload does not have github repo');
    }

    org.githubRepo = githubRepo;

    return await this.orgRepo.save({
      ...org
    });
  }
}

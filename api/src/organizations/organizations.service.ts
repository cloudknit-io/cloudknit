import { Injectable, Logger } from "@nestjs/common";
import { InjectRepository } from "@nestjs/typeorm";
import { Repository } from "typeorm";
import { Organization } from "src/typeorm";
import { CreateOrganizationDto } from "./Organization.dto";

@Injectable()
export class OrganizationsService {
  private readonly logger = new Logger(OrganizationsService.name);

  constructor(@InjectRepository(Organization) private orgRepo: Repository<Organization>) {}

  async getOrganizations() {
    return this.orgRepo.find();
  }

  async create(org: CreateOrganizationDto) {
    return this.orgRepo.save({
      name: org.name,
      clientId: org.clientId,
      clientSecret: org.clientSecret,
      githubRepo: org.githubRepo,
      githubPath: org.githubPath,
      githubSource: org.githubSource
    })
  }
}

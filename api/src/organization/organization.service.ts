import { Injectable } from "@nestjs/common";
import { InjectRepository } from "@nestjs/typeorm";
import { Repository } from "typeorm";
import { Organization } from "src/typeorm";

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
}

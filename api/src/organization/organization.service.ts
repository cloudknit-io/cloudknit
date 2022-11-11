import { BadRequestException, Injectable, Logger } from "@nestjs/common";
import { InjectRepository } from "@nestjs/typeorm";
import { Repository } from "typeorm";
import { Organization } from "src/typeorm";
import { PatchOrganizationDto } from "./organization.dto";
import { patchCompany } from "src/k8s/patch-company";

@Injectable()
export class OrganizationService {
  private readonly logger = new Logger(OrganizationService.name);

  constructor(@InjectRepository(Organization) private orgRepo: Repository<Organization>) { }

  async getOrganization(id: number) {
    return this.orgRepo.findOne({
      where: {
        id
      }
    });
  }

  async patchOrganization(org: Organization, payload: PatchOrganizationDto) {
    const reconcileProps = ['githubRepo', 'provisioned'];
    let changed = false;
    let updates = {};

    for (const prop of reconcileProps) {
      if (payload.hasOwnProperty(prop)) {
        changed = true;
        this.logger.log(`patching org ${prop} with value ${payload[prop]}`);
        updates[prop] = payload[prop];
      }
    }

    if (!changed) {
      throw new BadRequestException('organization was not updated');
    }

    await this.orgRepo.createQueryBuilder()
      .update(Organization)
      .set(updates)
      .where("id = :id", { id: org.id })
      .execute();

    if (payload.githubRepo) {
      try {
        await patchCompany(org, payload.githubRepo);
      } catch (error) {
        this.logger.log('Company CR was not updated', { org })
      }
    }

    this.logger.log({message: 'updated org', org, updates });

    return await this.orgRepo.findOneBy({
      id: org.id
    });
  }
}

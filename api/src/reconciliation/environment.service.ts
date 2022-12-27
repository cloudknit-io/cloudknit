import { Injectable, Logger } from "@nestjs/common";
import { InjectRepository } from "@nestjs/typeorm";
import { get } from "src/config";
import { Organization } from "src/typeorm";
import { Environment } from "src/typeorm/environment.entity";
import { Repository } from "typeorm";
import { EnvironmentDto } from "./dtos/environment.dto";
import { SSEService } from "./sse.service";

@Injectable()
export class EnvironmentService {
  private readonly config = get();
  private readonly ckEnvironment = this.config.environment;
  private readonly logger = new Logger(EnvironmentService.name);

  constructor(
    @InjectRepository(Environment)
    private readonly environmentRepository: Repository<Environment>,
    private readonly sseSvc: SSEService
  ) { }

  async putEnvironment(org: Organization, environment: EnvironmentDto) {
    const existing = await this.getEnvironment(org, environment.name, environment.teamName);

    if (!existing) {
      let env = new Environment();
      env.name = environment.name;
      env.teamName = environment.teamName;
      env.organization = org;
      env.duration = environment.duration;

      return await this.environmentRepository.save(env);
    }

    existing.duration = environment.duration;
    const entry = await this.environmentRepository.save(existing);
    entry.organization = org;
    this.sseSvc.sendEnvironment(entry);

    return entry;
  }

  async getEnvironment(org: Organization, envName: string, teamName: string) {
    return await this.environmentRepository.findOne({
      where: {
        name: envName,
        teamName,
        organization: {
          id : org.id
        }
      }
    });
  }
}

import { Injectable, Logger } from "@nestjs/common";
import { InjectRepository } from "@nestjs/typeorm";
import { SSEService } from "src/reconciliation/sse.service";
import { TeamService } from "src/team/team.service";
import { Environment, Organization, Team } from "src/typeorm";
import { Repository } from "typeorm";
import { CreateEnvironmentDto } from "./dto/create-environment.dto";
import { UpdateEnvironmentDto } from "./dto/update-environment.dto";

@Injectable()
export class EnvironmentService {
  private readonly logger = new Logger(EnvironmentService.name);

  constructor(
    @InjectRepository(Environment)
    private readonly envRepo: Repository<Environment>,
    private readonly teamSvc: TeamService,
    private readonly sseSvc: SSEService
  ) { }

  async create(env: CreateEnvironmentDto) {
    return this.envRepo.save(env);
  }

  async putEnvironment(org: Organization, environment: UpdateEnvironmentDto) {
    const team = await this.teamSvc.findByName(org, environment.teamName);
    const existing = await this.findByName(org, environment.name, team);

    if (!existing) {
      let env = new Environment();
      env.name = environment.name;
      env.team = team;
      env.organization = org;
      env.duration = environment.duration;

      return await this.envRepo.save(env);
    }

    existing.duration = environment.duration;
    const entry = await this.envRepo.save(existing);
    entry.organization = org;
    this.sseSvc.sendEnvironment(entry);

    return entry;
  }

  async findById(org: Organization, id: number, team?: Team, relations?: {}) {
    const where = {
      id,
      organization: {
        id : org.id
      },
      team: null
    };

    if (team) {
      where.team = {
        id: team.id
      }
    }

    return await this.envRepo.findOne({ where, relations });
  }

  async findByName(org: Organization, name: string, team?: Team, relations?: {}) {
    const where = {
      name,
      organization: {
        id : org.id
      },
      team: null
    };

    if (team) {
      where.team = {
        id: team.id
      }
    }

    return await this.envRepo.findOne({ where, relations });
  }
}

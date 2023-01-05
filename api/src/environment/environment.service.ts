import { Injectable, Logger } from "@nestjs/common";
import { InjectRepository } from "@nestjs/typeorm";
import { SSEService } from "src/reconciliation/sse.service";
import { Environment, Organization, Team } from "src/typeorm";
import { Equal, Repository } from "typeorm";
import { UpdateEnvironmentDto } from "./dto/update-environment.dto";

@Injectable()
export class EnvironmentService {
  private readonly logger = new Logger(EnvironmentService.name);

  constructor(
    @InjectRepository(Environment)
    private readonly envRepo: Repository<Environment>,
    private readonly sseSvc: SSEService
  ) { }

  async update(org: Organization, id: number, updateEnvDto: UpdateEnvironmentDto): Promise<Environment> {
    const env = await this.findById(org, id);

    this.envRepo.merge(env, updateEnvDto);

    const updatedEnv = await this.envRepo.save(env);
    this.sseSvc.sendEnvironment(updatedEnv);

    return updatedEnv;
  }

  async updateByName(org: Organization, team: Team, name: string, updateEnvDto: UpdateEnvironmentDto): Promise<Environment> {
    const env = await this.findByName(org, team, name);

    this.envRepo.merge(env, updateEnvDto);

    const updatedEnv = await this.envRepo.save(env);
    this.sseSvc.sendEnvironment(updatedEnv);

    return updatedEnv;
  }

  async findById(org: Organization, id: number, withTeam: boolean = false): Promise<Environment> {
    return this.envRepo.findOne({
      where: {
        id: Equal(id),
        organization: {
          id: Equal(org.id)
        },
      },
      relations: {
        team: withTeam
      }
    });
  }

  async findByName(org: Organization, team: Team, name: string, relations?: {}) {
    const where = {
      name: Equal(name),
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

  async remove(org: Organization, id: number): Promise<Environment> {
    const env = await this.findById(org, id);

    env.isDeleted = true;

    const updatedEnv = await this.envRepo.save(env);
    this.sseSvc.sendEnvironment(updatedEnv);

    return updatedEnv;
  }
}

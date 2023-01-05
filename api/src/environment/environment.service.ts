import { Injectable } from "@nestjs/common";
import { InjectRepository } from "@nestjs/typeorm";
import { Environment, Organization, Team } from "src/typeorm";
import { Equal, Repository } from "typeorm";
import { UpdateEnvironmentDto } from "./dto/update-environment.dto";

@Injectable()
export class EnvironmentService {
  constructor(
    @InjectRepository(Environment)
    private readonly envRepo: Repository<Environment>,
  ) { }

  async update(org: Organization, id: number, updateEnvDto: UpdateEnvironmentDto): Promise<Environment> {
    const env = await this.findById(org, id);

    this.envRepo.merge(env, updateEnvDto);
    env.organization = org;

    return this.envRepo.save(env);
  }

  async updateByName(org: Organization, team: Team, name: string, updateEnvDto: UpdateEnvironmentDto): Promise<Environment> {
    const env = await this.findByName(org, team, name);

    this.envRepo.merge(env, updateEnvDto);
    env.organization = org;

    return this.envRepo.save(env);
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
    env.organization = org;

    return this.envRepo.save(env);
  }
}

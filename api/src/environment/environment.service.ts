import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Environment, Organization, Team } from 'src/typeorm';
import { Equal, Repository } from 'typeorm';
import { CreateEnvironmentDto } from './dto/create-environment.dto';
import { UpdateEnvironmentDto } from './dto/update-environment.dto';

@Injectable()
export class EnvironmentService {
  constructor(
    @InjectRepository(Environment)
    private readonly envRepo: Repository<Environment>
  ) {}

  async create(
    org: Organization,
    team: Team,
    createEnvDto: CreateEnvironmentDto
  ) {
    return this.envRepo.save({
      organization: org,
      team,
      ...createEnvDto,
    });
  }

  async findAll(org: Organization, team: Team) {
    return this.envRepo.find({
      where: {
        team: {
          id: team.id,
        },
        organization: {
          id: org.id,
        },
      },
    });
  }

  async updateById(
    org: Organization,
    id: number,
    updateEnvDto: UpdateEnvironmentDto
  ): Promise<Environment> {
    const env = await this.findById(org, id);
    this.envRepo.merge(env, updateEnvDto);
    return this.envRepo.save(env);
  }

  async updateByName(
    org: Organization,
    team: Team,
    name: string,
    updateEnvDto: UpdateEnvironmentDto
  ): Promise<Environment> {
    const env = await this.findByName(org, team, name);
    this.envRepo.merge(env, updateEnvDto);
    return this.envRepo.save(env);
  }

  async findById(
    org: Organization,
    id: number,
    withTeam: boolean = false,
    withComps: boolean = false
  ): Promise<Environment> {
    return this.envRepo.findOne({
      where: {
        id: Equal(id),
        organization: {
          id: Equal(org.id),
        },
      },
      relations: {
        team: withTeam,
        components: withComps,
      },
    });
  }

  async findByName(
    org: Organization,
    team: Team,
    name: string,
    withTeam: boolean = false,
    withComps: boolean = false
  ) {
    return this.envRepo.findOne({
      where: {
        name: Equal(name),
        organization: {
          id: org.id,
        },
        team: {
          id: team.id,
        },
      },
      relations: {
        team: withTeam,
        components: withComps,
      },
    });
  }

  async remove(org: Organization, id: number): Promise<Environment> {
    const env = await this.findById(org, id);

    env.isDeleted = true;

    return this.envRepo.save(env);
  }

  async updateCost(org: Organization, env: Environment) {
    env = await this.findById(org, env.id, false, true);

    let estimatedCost = 0.0;

    for (const comp of env.components) {
      if (comp.estimatedCost > 0) {
        estimatedCost += comp.estimatedCost;
      }
    }

    await this.updateById(org, env.id, { estimatedCost });
  }
}

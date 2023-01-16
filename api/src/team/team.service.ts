import { Injectable, Logger } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Organization, Team } from 'src/typeorm';
import { Equal, FindOptionsRelations, Repository } from 'typeorm';
import { CreateTeamDto } from './dto/create-team.dto';
import { UpdateTeamDto } from './dto/update-team.dto';

@Injectable()
export class TeamService {
  private readonly logger = new Logger(TeamService.name);

  constructor(
    @InjectRepository(Team)
    private teamRepo: Repository<Team>
  ) {}

  async create(createTeamDto: CreateTeamDto) {
    return this.teamRepo.save(createTeamDto);
  }

  async findAll(
    org: Organization,
    withEnv: boolean = false,
    withComponents: boolean = false
  ) {
    console.log(withComponents);
    let relation: FindOptionsRelations<Team> = {
      environments: withEnv,
    };

    if (withComponents) {
      relation = {
        environments: {
          components: true,
        },
      };
    }
    return this.teamRepo.find({
      where: {
        organization: {
          id: org.id,
        },
      },
      relations: relation,
    });
  }

  async findByName(org: Organization, name: string): Promise<Team> {
    return this.teamRepo.findOne({
      where: {
        name: Equal(name),
        organization: {
          id: Equal(org.id),
        },
      },
    });
  }

  async findById(org: Organization, id: number): Promise<Team> {
    return this.teamRepo.findOne({
      where: {
        id: Equal(id),
        organization: {
          id: Equal(org.id),
        },
      },
    });
  }

  async update(
    org: Organization,
    id: number,
    updateTeamDto: UpdateTeamDto
  ): Promise<Team> {
    const team = await this.findById(org, id);

    this.teamRepo.merge(team, updateTeamDto);

    return this.teamRepo.save(team);
  }

  async remove(org: Organization, id: number): Promise<Team> {
    const team = await this.findById(org, id);

    team.isDeleted = true;

    return this.teamRepo.save(team);
  }
}

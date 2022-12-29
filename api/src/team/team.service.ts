import { Injectable, Logger } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Organization, Team } from 'src/typeorm';
import { Repository } from 'typeorm';
import { UpdateTeamDto } from './dto/update-team.dto';

@Injectable()
export class TeamService {
  private readonly logger = new Logger(TeamService.name);

  constructor(
    @InjectRepository(Team)
    private teamRepo: Repository<Team>
  ) {}

  async findByName(org: Organization, name: string): Promise<Team> {
    return this.teamRepo.findOne({
      where: {
        name,
        organization: {
          id: org.id
        }
      }
    })
  }

  async findById(org: Organization, id: number): Promise<Team> {
    return this.teamRepo.findOne({
      where: {
        id,
        organization: {
          id: org.id
        }
      }
    })
  }

  async update(org: Organization, id: number, updateTeamDto: UpdateTeamDto): Promise<Team> {
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

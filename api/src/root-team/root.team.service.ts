import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { CreateTeamDto } from 'src/team/dto/create-team.dto';
import { Organization, Team } from 'src/typeorm';
import { Repository } from 'typeorm';

@Injectable()
export class RootTeamService {
  constructor(
    @InjectRepository(Team)
    private teamRepo: Repository<Team>,
  ) {}
  
  async create(createTeamDto: CreateTeamDto) {
    return this.teamRepo.save(createTeamDto);
  }

  async findAll(org: Organization, withEnv: boolean = false) {
    return this.teamRepo.find({
      where: {
        organization: {
          id: org.id
        }
      },
      relations: {
        environments: withEnv
      }
    })
  }
}

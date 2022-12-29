import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { CreateEnvironmentDto } from 'src/environment/dto/create-environment.dto';
import { Environment, Organization, Team } from 'src/typeorm';
import { Repository } from 'typeorm';

@Injectable()
export class RootEnvironmentService {
  constructor(
    @InjectRepository(Environment)
    private envRepo: Repository<Environment>,
  ) {}

  async create(createEnvDto: CreateEnvironmentDto) {
    return this.envRepo.save(createEnvDto);
  }

  async findAll(org: Organization, team: Team) {
    return this.envRepo.find({
      where: {
        team: {
          id: team.id
        },
        organization: {
          id: org.id
        }
      }
    })
  }
}

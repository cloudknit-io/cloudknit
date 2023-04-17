import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { UpdateEnvironmentDto } from 'src/environment/dto/update-environment.dto';
import { Environment, Organization, Team } from 'src/typeorm';
import { Repository } from 'typeorm';

@Injectable()
export class ErrorsService {
  constructor(
    @InjectRepository(Environment)
    private readonly envRepo: Repository<Environment>
  ) {}

  async updateEnv(env: Environment, envUpdate: UpdateEnvironmentDto) {
    await this.envRepo.merge(env, envUpdate);
    return this.envRepo.save(env);
  }
}

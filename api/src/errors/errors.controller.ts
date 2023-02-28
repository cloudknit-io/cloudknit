import { Body, Controller, InternalServerErrorException, Logger, Post, Request } from '@nestjs/common';
import { ApiTags } from '@nestjs/swagger';
import { CreateEnvironmentDto } from 'src/environment/dto/create-environment.dto';
import { EnvironmentService } from 'src/environment/environment.service';
import { Environment, Organization, Team } from 'src/typeorm';
import { APIRequest, TeamApiParam } from 'src/types';
import { handleSqlErrors } from 'src/utilities/errorHandler';
import { EnvironmentErrorSpecDto } from './dto/environment-error.dto';

@Controller({
  version: '1',
})
@ApiTags('errors')
export class ErrorsController {
  private readonly logger = new Logger(ErrorsController.name);
  constructor(private readonly envSvc: EnvironmentService) {}

  @Post()
  @TeamApiParam()
  async saveOrUpdate(@Request() req: APIRequest, @Body() body: EnvironmentErrorSpecDto) {
    const { org, team } = req;
    let env = await this.envSvc.findByName(org, team, body.envName);
    if (!env) {
      env = await this.createEnv(org, team, {
        name: body.envName,
        dag: [],
      });
    }

    return await this.envSvc.updateById(org, env.id, {
      dag: [],
      name: env.name,
      errorType: body.errorType,
      errorMessage: body.errorMessage,
      status: body.status
    });
  }

  async createEnv(
    org: Organization,
    team: Team,
    createEnv: CreateEnvironmentDto
  ): Promise<Environment> {
    let env: Environment;

    try {
      env = await this.envSvc.create(org, team, createEnv);
      this.logger.log({ message: `created new environment`, env });
    } catch (err) {
      handleSqlErrors(err, 'environment already exists');

      this.logger.error({
        message: 'could not create environment',
        createEnv,
        err,
      });
      throw new InternalServerErrorException('could not create environment');
    }

    return env;
  }
}

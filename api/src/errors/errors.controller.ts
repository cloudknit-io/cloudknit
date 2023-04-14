import {
  Body,
  Controller,
  InternalServerErrorException,
  Logger,
  Post,
  Request
} from '@nestjs/common';
import { ApiTags } from '@nestjs/swagger';
import { EnvironmentService } from 'src/environment/environment.service';
import { ReconciliationService } from 'src/reconciliation/reconciliation.service';
import { APIRequest, TeamApiParam } from 'src/types';
import { handleSqlErrors } from 'src/utilities/errorHandler';
import { EnvironmentErrorSpecDto } from './dto/environment-error.dto';

@Controller({
  version: '1',
})
@ApiTags('errors')
export class ErrorsController {
  private readonly logger = new Logger(ErrorsController.name);
  constructor(
    private readonly envSvc: EnvironmentService,
    private readonly reconSvc: ReconciliationService,
  ) {}

  @Post()
  @TeamApiParam()
  async saveOrUpdate(
    @Request() req: APIRequest,
    @Body() body: EnvironmentErrorSpecDto
  ) {
    try {
      const { org, team } = req;
      console.log('body: ', body);
      let env = await this.envSvc.findByName(org, team, body.envName);
      if (!env) {
        env = await this.envSvc.create(org, team, {
          name: body.envName,
          dag: [],
        });
      }
      if (JSON.stringify(body.errorMessage) == JSON.stringify(env.latestEnvRecon?.errorMessage)) {
        return;
      }

      const envRecon = await this.reconSvc.createEnvRecon(org, team, env, {
        name: env.name,
        startDateTime: new Date().toISOString(),
        errorMessage: body.errorMessage,
        components: env.dag,
        teamName: team.name,
      });

      await this.reconSvc.updateEnvRecon(envRecon, {
        status: 'Failed',
        endDateTime: new Date().toISOString(),
      });

      envRecon.environment = null;

      await this.envSvc.mergeAndSaveEnv(org, env, {
        latestEnvRecon: envRecon
      });

    } catch (err) {
      handleSqlErrors(err, 'environment already exists');

      this.logger.error({
        message: 'could not create environment',
        body,
        err,
      });
      throw new InternalServerErrorException('could not create environment');
    }
  }
}

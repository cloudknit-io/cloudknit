import {
  Body,
  Controller,
  InternalServerErrorException,
  Logger,
  Post,
  Request,
} from '@nestjs/common';
import { ApiTags } from '@nestjs/swagger';
import { EnvironmentService } from 'src/environment/environment.service';
import { APIRequest, TeamApiParam } from 'src/types';
import { handleSqlErrors } from 'src/utilities/errorHandler';
import { EnvironmentErrorSpecDto } from './dto/environment-error.dto';
import { ErrorsService } from './errors.service';

@Controller({
  version: '1',
})
@ApiTags('errors')
export class ErrorsController {
  private readonly logger = new Logger(ErrorsController.name);
  constructor(
    private readonly envSvc: EnvironmentService,
    private readonly errorSvc: ErrorsService
  ) {}

  @Post()
  @TeamApiParam()
  async saveOrUpdate(
    @Request() req: APIRequest,
    @Body() body: EnvironmentErrorSpecDto
  ) {
    try {
      const { org, team } = req;
      let env = await this.envSvc.findByName(org, team, body.envName);
      if (!env) {
        env = await this.envSvc.create(org, team, {
          name: body.envName,
          dag: [],
        });
      }

      if (!body.errorMessage) {
        return this.errorSvc.processValidRecon(org, team, env);
      }

      return this.errorSvc.processInvalidRecon(
        org,
        team,
        env,
        body.errorMessage
      );
    } catch (err) {
      handleSqlErrors(err, 'There was an error processing error request');

      this.logger.error({
        message: 'could not create error env reconcile entry',
        body,
        err,
      });
      throw new InternalServerErrorException('could not create error env reconcile entry');
    }
  }
}

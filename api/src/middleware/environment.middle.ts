import { BadRequestException, HttpException, HttpStatus, Injectable, Logger, NestMiddleware } from '@nestjs/common';
import { Response, NextFunction } from 'express';
import { EnvironmentService } from 'src/environment/environment.service';
import { APIRequest } from 'src/types';

@Injectable()
export class EnvironmentMiddleware implements NestMiddleware {
  private readonly logger = new Logger(EnvironmentMiddleware.name);

  constructor(private readonly envSvc: EnvironmentService) {}

  async use(req: APIRequest, res: Response, next: NextFunction) {
    const { org, team } = req;
    const envId = req.params.environmentId;

    let env, id;

    try {
      id = parseInt(req.params.environmentId, 10);
    } catch (e) {}

    if (isNaN(id)) {
      try {
        env = await this.envSvc.findByName(org, team, envId);
      } catch (e) {
        this.logger.error({
          message: 'could not get environment by name',
          envId,
          error: e.message,
        });
      }
    } else {
      try {
        env = await this.envSvc.findById(org, id, true);
      } catch (e) {
        this.logger.error({
          message: 'could not get environment by number',
          envId,
          error: e.message,
        });
      }
    }

    if (!env) {
      this.logger.error({ message: 'bad environmentId', envId });
      throw new BadRequestException('environment not found');
    }

    req.env = env;

    next();
  }
}

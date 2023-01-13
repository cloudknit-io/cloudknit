import { HttpException, HttpStatus, Injectable, Logger, NestMiddleware } from '@nestjs/common';
import { Response, NextFunction } from 'express';
import { APIRequest } from 'src/types';
import { TeamService } from 'src/team/team.service';

@Injectable()
export class TeamMiddleware implements NestMiddleware {
  private readonly logger = new Logger(TeamMiddleware.name);

  constructor(
    private readonly teamSvc: TeamService
  ) {}

  async use(req: APIRequest, res: Response, next: NextFunction) {
    const org = req.org;
    const teamId = req.params.teamId;
    
    let team, id;

    try {
      id = parseInt(teamId, 10);
    } catch (e) {}
    
    if (isNaN(id)) {
      try {
        team = await this.teamSvc.findByName(org, teamId)
      } catch (e) {
        this.logger.error({message: 'could not get team by name', teamId, error: e.message})
      }
    } else {
      try {
        team = await this.teamSvc.findById(org, id);
      } catch (e) {
        this.logger.error({message: 'could not get team by number', teamId, error: e.message})
      }
    }

    if (!team) {
      this.logger.error({ message: 'bad teamId', teamId});
      throw new HttpException('Forbidden', HttpStatus.FORBIDDEN);
    }

    req.team = team;

    next();
  }
}

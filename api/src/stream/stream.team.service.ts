import { Inject, Injectable, Logger } from '@nestjs/common';
import { Team } from 'src/typeorm';
import {
  Connection,
  EntitySubscriberInterface,
  InsertEvent,
  RemoveEvent,
  UpdateEvent,
} from 'typeorm';
import { StreamService } from './stream.service';

@Injectable()
export class StreamTeamService implements EntitySubscriberInterface<Team> {
  private readonly logger = new Logger(StreamTeamService.name);

  constructor(
    @Inject(Connection) conn: Connection,
    private readonly sseSvc: StreamService
  ) {
    conn.subscribers.push(this);
  }

  listenTo(): string | Function {
    return Team;
  }

  afterInsert(event: InsertEvent<Team>) {
    this.validateAndSend(event.entity, 'afterInsert');
  }

  afterUpdate(event: UpdateEvent<Team>): void | Promise<Team> {
    this.validateAndSend(event.entity as Team, 'afterUpdate');
  }

  afterRemove(event: RemoveEvent<Team>): void | Promise<any> {
    this.validateAndSend(event.entity, 'afterRemove');
  }

  validateAndSend(team: Team, operation: string) {
    if (team.organization || team.orgId) {
      this.sseSvc.sendTeam(team);
      return;
    }

    this.logger.error({
      message: 'team stream object has no organization',
      env: team,
      operation,
    });
  }
}

import { Module } from '@nestjs/common';
import { StreamService } from './stream.service';
import { StreamController } from './stream.controller';
import { StreamEnvironmentService } from './stream.environment.service';
import { StreamComponentService } from './stream.component.service';
import { StreamEnvironmentReconcileService } from './stream.env-reconcile.service';
import { StreamComponentReconcileService } from './stream.comp-reconcile.service';
import { StreamTeamService } from './stream.team.service';
import { EnvironmentService } from 'src/environment/environment.service';
import { TeamService } from 'src/team/team.service';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Environment, Team } from 'src/typeorm';

@Module({
  imports: [TypeOrmModule.forFeature([Environment, Team])],
  controllers: [StreamController],
  providers: [
    StreamService,
    StreamEnvironmentService,
    StreamComponentService,
    StreamEnvironmentReconcileService,
    StreamComponentReconcileService,
    StreamTeamService,
    EnvironmentService,
    TeamService,
  ],
})
export class StreamModule {}

import { Module } from '@nestjs/common';
import { RootTeamService } from './root.team.service';
import { RootTeamController } from './root.team.controller';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Team } from 'src/typeorm';
import { TeamService } from 'src/team/team.service';

@Module({
  imports: [
    TypeOrmModule.forFeature([
      Team
    ])
  ],
  controllers: [RootTeamController],
  providers: [
    RootTeamService,
    TeamService
  ]
})
export class RootTeamModule {}

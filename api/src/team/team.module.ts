import {
  MiddlewareConsumer,
  Module,
  NestModule,
  RequestMethod,
} from '@nestjs/common';
import { TeamService } from './team.service';
import { TeamController } from './team.controller';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Component, Environment, Team } from 'src/typeorm';
import { EnvironmentService } from 'src/environment/environment.service';
import { TeamMiddleware } from 'src/middleware/team.middle';

@Module({
  imports: [TypeOrmModule.forFeature([Team, Component, Environment])],
  controllers: [TeamController],
  providers: [TeamService, EnvironmentService],
})
export class TeamModule implements NestModule {
  configure(consumer: MiddlewareConsumer) {
    consumer.apply(TeamMiddleware).forRoutes({
      path: '*/teams/:teamId*',
      method: RequestMethod.ALL,
    });
  }
}

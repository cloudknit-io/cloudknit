import {
  MiddlewareConsumer,
  Module,
  NestModule,
  RequestMethod
} from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { EnvironmentService } from 'src/environment/environment.service';
import { TeamMiddleware } from 'src/middleware/team.middle';
import { Component, Environment, Team } from 'src/typeorm';
import { TeamController } from './team.controller';
import { TeamService } from './team.service';

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

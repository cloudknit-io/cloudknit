import { MiddlewareConsumer, Module, NestModule, RequestMethod } from '@nestjs/common';
import { TeamService } from './team.service';
import { TeamController } from './team.controller';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Component, Environment, Team } from 'src/typeorm';
import { ComponentService } from 'src/costing/services/component.service';
import { EnvironmentService } from 'src/environment/environment.service';
import { SSEService } from 'src/reconciliation/sse.service';
import { TeamMiddleware } from 'src/middleware/team.middle';

@Module({
  imports: [
    TypeOrmModule.forFeature([
      Team,
      Component,
      Environment,
    ])
  ],
  controllers: [TeamController],
  providers: [
    TeamService,
    ComponentService,
    EnvironmentService,
    SSEService
  ]
})
export class TeamModule implements NestModule {
  configure(consumer: MiddlewareConsumer) {
    consumer.apply(TeamMiddleware).forRoutes({
      path: '*/teams/:teamId*',
      method: RequestMethod.ALL
    });
  }
}
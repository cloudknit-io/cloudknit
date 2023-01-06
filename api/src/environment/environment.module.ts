import { MiddlewareConsumer, Module, NestModule, RequestMethod } from '@nestjs/common';
import { EnvironmentService } from './environment.service';
import { EnvironmentController } from './environment.controller';
import { TeamService } from 'src/team/team.service';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Component, ComponentReconcile, Environment, EnvironmentReconcile, Team } from 'src/typeorm';
import { EnvironmentMiddleware } from 'src/middleware/environment.middle';
import { ReconciliationService } from 'src/reconciliation/reconciliation.service';

@Module({
  imports: [
    TypeOrmModule.forFeature([
      Environment,
      Team,
      Component,
      EnvironmentReconcile,
      ComponentReconcile
    ])
  ],
  controllers: [EnvironmentController],
  providers: [
    EnvironmentService,
    TeamService,
    ReconciliationService
  ]
})
export class EnvironmentModule implements NestModule {
  configure(consumer: MiddlewareConsumer) {
    consumer.apply(EnvironmentMiddleware).forRoutes({
      path: '*/environments/:environmentId*',
      method: RequestMethod.ALL
    });
  }
}

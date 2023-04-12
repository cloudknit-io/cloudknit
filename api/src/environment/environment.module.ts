import {
  MiddlewareConsumer,
  Module,
  NestModule,
  RequestMethod,
} from '@nestjs/common';
import { EnvironmentService } from './environment.service';
import { EnvironmentController } from './environment.controller';
import { TeamService } from 'src/team/team.service';
import { TypeOrmModule } from '@nestjs/typeorm';
import {
  Component,
  ComponentReconcile,
  Environment,
  EnvironmentReconcile,
  Team,
} from 'src/typeorm';
import { EnvironmentMiddleware } from 'src/middleware/environment.middle';
import { ReconciliationService } from 'src/reconciliation/reconciliation.service';
import { ComponentService } from 'src/component/component.service';
import { SystemService } from 'src/system/system.service';

@Module({
  imports: [
    TypeOrmModule.forFeature([
      Environment,
      Team,
      Component,
      EnvironmentReconcile,
      ComponentReconcile,
    ]),
  ],
  controllers: [EnvironmentController],
  providers: [
    ComponentService,
    EnvironmentService,
    TeamService,
    ReconciliationService,
    SystemService
  ],
})
export class EnvironmentModule implements NestModule {
  configure(consumer: MiddlewareConsumer) {
    consumer.apply(EnvironmentMiddleware).forRoutes({
      path: '*/environments/:environmentId*',
      method: RequestMethod.ALL,
    });
  }
}

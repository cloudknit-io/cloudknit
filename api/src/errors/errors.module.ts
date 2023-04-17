import {
    MiddlewareConsumer,
    Module,
    NestModule,
    RequestMethod,
  } from '@nestjs/common';
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
import { EnvironmentService } from 'src/environment/environment.service';
import { ErrorsController } from './errors.controller';
  
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
    controllers: [ErrorsController],
    providers: [
      ComponentService,
      EnvironmentService,
      TeamService,
      ReconciliationService,
    ],
  })
  export class ErrorsModule implements NestModule {
    configure(consumer: MiddlewareConsumer) {
      consumer.apply(EnvironmentMiddleware).forRoutes({
        path: '*/errors/:environmentId*',
        method: RequestMethod.ALL,
      });
    }
  }
  
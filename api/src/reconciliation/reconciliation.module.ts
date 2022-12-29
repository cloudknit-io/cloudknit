import { Module } from '@nestjs/common'
import { TypeOrmModule } from '@nestjs/typeorm';
import { ComponentService } from 'src/costing/services/component.service';
import { EnvironmentService } from 'src/environment/environment.service';
import { RootEnvironmentService } from 'src/root-environment/root.environment.service';
import { TeamService } from 'src/team/team.service';
import { Component, ComponentReconcile, Environment, EnvironmentReconcile, Team } from 'src/typeorm';
import { ReconciliationController } from './reconciliation.controller';
import { ReconciliationService } from './reconciliation.service';
import { SSEService } from './sse.service';


@Module({
  imports: [
    TypeOrmModule.forFeature([
      EnvironmentReconcile,
      ComponentReconcile,
      Environment,
      Component,
      Team
    ])
  ],
  controllers: [
    ReconciliationController
  ],
  providers: [
    ReconciliationService,
    RootEnvironmentService,
    EnvironmentService,
    SSEService,
    TeamService,
    ComponentService
  ],
})
export class ReconciliationModule {
}

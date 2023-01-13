import { Module } from '@nestjs/common';
import { ComponentService } from './component.service';
import { ComponentController } from './component.controller';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Component, ComponentReconcile, Environment, EnvironmentReconcile, Team } from 'src/typeorm';
import { ReconciliationService } from 'src/reconciliation/reconciliation.service';
import { TeamService } from 'src/team/team.service';
import { EnvironmentService } from 'src/environment/environment.service';

@Module({
  imports: [
    TypeOrmModule.forFeature([
      Component,
      Environment,
      EnvironmentReconcile,
      ComponentReconcile,
      Team
    ])
  ],
  controllers: [ComponentController],
  providers: [
    ComponentService,
    ReconciliationService,
    EnvironmentService,
    TeamService
  ]
})
export class ComponentModule {}

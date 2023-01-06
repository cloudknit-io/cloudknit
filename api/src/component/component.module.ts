import { Module } from '@nestjs/common';
import { ComponentService } from './component.service';
import { ComponentController } from './component.controller';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Component, ComponentReconcile, EnvironmentReconcile, Team } from 'src/typeorm';
import { ReconciliationService } from 'src/reconciliation/reconciliation.service';
import { TeamService } from 'src/team/team.service';

@Module({
  imports: [
    TypeOrmModule.forFeature([
      Component,
      EnvironmentReconcile,
      ComponentReconcile,
      Team
    ])
  ],
  controllers: [ComponentController],
  providers: [
    ComponentService,
    ReconciliationService,
    TeamService
  ]
})
export class ComponentModule {}

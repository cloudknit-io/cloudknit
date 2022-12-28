import { Module } from '@nestjs/common';
import { EnvironmentService } from './environment.service';
import { EnvironmentController } from './environment.controller';
import { TeamService } from 'src/team/team.service';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Component, Environment, Team } from 'src/typeorm';
import { SSEService } from 'src/reconciliation/sse.service';
import { ComponentService } from 'src/costing/services/component.service';

@Module({
  imports: [
    TypeOrmModule.forFeature([
      Environment,
      Team,
      Component
    ])
  ],
  controllers: [EnvironmentController],
  providers: [
    EnvironmentService,
    TeamService,
    SSEService,
    ComponentService
  ]
})
export class EnvironmentModule {}

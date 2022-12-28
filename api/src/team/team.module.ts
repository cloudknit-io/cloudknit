import { Module } from '@nestjs/common';
import { TeamService } from './team.service';
import { TeamController } from './team.controller';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Component, Environment, Team } from 'src/typeorm';
import { ComponentService } from 'src/costing/services/component.service';
import { EnvironmentService } from 'src/environment/environment.service';
import { SSEService } from 'src/reconciliation/sse.service';

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
export class TeamModule {}

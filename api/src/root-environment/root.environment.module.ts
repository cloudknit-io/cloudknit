import { Module } from '@nestjs/common';
import { RootEnvironmentService } from './root.environment.service';
import { RootEnvironmentController } from './root.environment.controller';
import { Component, Environment, Team } from 'src/typeorm';
import { TypeOrmModule } from '@nestjs/typeorm';
import { ComponentService } from 'src/costing/services/component.service';
import { EnvironmentService } from 'src/environment/environment.service';
import { TeamService } from 'src/team/team.service';

@Module({
  imports: [
    TypeOrmModule.forFeature([
      Environment,
      Component,
      Team
    ])
  ],
  controllers: [RootEnvironmentController],
  providers: [
    ComponentService,
    EnvironmentService,
    RootEnvironmentService,
    TeamService
  ]
})
export class RootEnvironmentModule {}

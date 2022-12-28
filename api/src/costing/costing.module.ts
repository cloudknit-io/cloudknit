import { Module } from '@nestjs/common'
import { TypeOrmModule } from '@nestjs/typeorm';
import { EnvironmentService } from 'src/environment/environment.service';
import { SSEService } from 'src/reconciliation/sse.service';
import { TeamService } from 'src/team/team.service';
import { Team } from 'src/typeorm';
import { Component } from 'src/typeorm/component.entity';
import { Environment } from 'src/typeorm/environment.entity';
import { CostingController } from './costing.controller';
import { ComponentService } from './services/component.service';


@Module({
  imports: [
    TypeOrmModule.forFeature([
      Component,
      Environment,
      Team
    ])
  ],
  controllers: [
      CostingController,
  ],
  providers: [
    ComponentService,
    EnvironmentService,
    TeamService,
    SSEService
  ],
})
export class CostingModule {
}

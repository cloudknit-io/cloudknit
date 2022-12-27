import { Module } from '@nestjs/common'
import { TypeOrmModule } from '@nestjs/typeorm';
import { ComponentReconcile } from 'src/typeorm/component-reconcile.entity';
import { Component } from 'src/typeorm/component.entity';
import { EnvironmentReconcile } from 'src/typeorm/environment-reconcile.entity';
import { Environment } from 'src/typeorm/environment.entity';
import { EnvironmentService } from './environment.service';
import { ReconciliationController } from './reconciliation.controller';
import { ReconciliationService } from './reconciliation.service';
import { SSEService } from './sse.service';


@Module({
  imports: [
    TypeOrmModule.forFeature([
      EnvironmentReconcile,
      ComponentReconcile,
      Environment,
      Component
    ])
  ],
  controllers: [
    ReconciliationController
  ],
  providers: [
    ReconciliationService,
    EnvironmentService,
    SSEService
  ],
})
export class ReconciliationModule {
}

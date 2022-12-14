import { Module } from '@nestjs/common'
import { TypeOrmModule } from '@nestjs/typeorm';
import { Component } from 'src/typeorm/component.entity';
import { Environment } from 'src/typeorm/reconciliation/environment.entity';
import { resourceEntities } from 'src/typeorm/resources';
import { CostingController } from './costing.controller';
import { ComponentService } from './services/component.service';
import { CostingStream } from './streams/costing.stream';


@Module({
  imports: [
    TypeOrmModule.forFeature([
      Component,
      ...resourceEntities,
      Environment
    ])
  ],
  controllers: [
      CostingController,
      CostingStream
  ],
  providers: [ComponentService],
})
export class CostingModule {
}

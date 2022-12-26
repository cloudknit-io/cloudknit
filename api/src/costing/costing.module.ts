import { Module } from '@nestjs/common'
import { TypeOrmModule } from '@nestjs/typeorm';
import { Component } from 'src/typeorm/component.entity';
import { Environment } from 'src/typeorm/reconciliation/environment.entity';
import { CostingController } from './costing.controller';
import { ComponentService } from './services/component.service';
import { CostingStream } from './streams/costing.stream';


@Module({
  imports: [
    TypeOrmModule.forFeature([
      Component,
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

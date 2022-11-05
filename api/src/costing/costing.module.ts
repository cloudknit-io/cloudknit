import { Module } from '@nestjs/common'
import { TypeOrmModule } from '@nestjs/typeorm';
import { costingEntities } from 'src/typeorm/costing';
import { Environment } from 'src/typeorm/reconciliation/environment.entity';
import { resourceEntities } from 'src/typeorm/resources';
import { CostingController } from './costing.controller';
import { ComponentService } from './services/component.service';
import { CostingStream } from './streams/costing.stream';


@Module({
  imports: [
    TypeOrmModule.forFeature([
      ...costingEntities,
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

import { Module } from '@nestjs/common'
import { TypeOrmModule } from '@nestjs/typeorm';
import { costingEntities } from 'src/typeorm/costing';
import { resourceEntities } from 'src/typeorm/resources';
import { CostingController } from './costing.controller';
import { ComponentService } from './services/component.service';
import { CostingStream } from './streams/costing.stream';


@Module({
  imports: [
    TypeOrmModule.forFeature([
      ...costingEntities,
      ...resourceEntities
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

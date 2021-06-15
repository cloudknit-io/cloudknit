import { Module } from '@nestjs/common'
import { TypeOrmModule } from '@nestjs/typeorm';
import { costingEntities } from 'src/typeorm/costing';
import { CostingController } from './costing.controller';
import { ComponentService } from './services/component.service';


@Module({
  imports: [
    TypeOrmModule.forFeature([
      ...costingEntities
    ])
  ],
  controllers: [
      CostingController
  ],
  providers: [ComponentService],
})
export class CostingModule {
}

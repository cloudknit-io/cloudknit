import { Module } from '@nestjs/common'
import { TypeOrmModule } from '@nestjs/typeorm';
import { reconcileEntities } from 'src/typeorm/reconciliation';
import { ReconciliationController } from './reconciliation.controller';
import { ReconciliationService } from './services/reconciliation.service';


@Module({
  imports: [
    TypeOrmModule.forFeature([
      ...reconcileEntities
    ])
  ],
  controllers: [
    ReconciliationController
  ],
  providers: [
    ReconciliationService
  ],
})
export class ReconciliationModule {
}

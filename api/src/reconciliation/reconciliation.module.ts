import { Module } from '@nestjs/common'
import { TypeOrmModule } from '@nestjs/typeorm';
import { reconcileEntities } from 'src/typeorm/reconciliation';
import { EnvironmentService } from './environment.service';
import { ReconciliationController } from './reconciliation.controller';
import { ReconciliationService } from './reconciliation.service';
import { SSEService } from './sse.service';


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
    ReconciliationService,
    EnvironmentService,
    SSEService
  ],
})
export class ReconciliationModule {
}

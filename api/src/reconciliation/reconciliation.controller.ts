import { Body, Controller, Post } from '@nestjs/common'
import { EvnironmentReconcileDto } from './dtos/reconcile.Dto';
import { ReconciliationService } from './services/reconciliation.service'


@Controller({
  path: 'reconciliation/api/v1',
})
export class ReconciliationController {
  constructor(private readonly reconciliationService: ReconciliationService) {}

  @Post('environment/save')
  async saveEnvironment(@Body() runData: EvnironmentReconcileDto) {
    return await this.reconciliationService.saveOrUpdateEnvironment(runData);
  }

  @Post('component/save')
  async saveComponent(@Body() runData: EvnironmentReconcileDto) {
    return await this.reconciliationService.saveOrUpdateComponent(runData);
  }

}

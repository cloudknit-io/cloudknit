import { Body, Controller, Get, Param, Post, Sse } from '@nestjs/common'
import { Observable, Observer } from 'rxjs';
import { Mapper } from 'src/costing/utilities/mapper';
import { ComponentReconcile } from 'src/typeorm/reconciliation/component-reconcile.entity';
import { ComponentAudit } from './dtos/componentAudit.dto';
import { EnvironmentAudit } from './dtos/environmentAudit.dto';
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

  @Get('component/:id')
  async getComponents(@Param('id') id: string): Promise<ComponentAudit[]> {
    return await this.reconciliationService.getComponentAuditList(id);
  }

  @Get('environment/:id')
  async getEnvironments(@Param('id') id: string): Promise<EnvironmentAudit[]> {
    return await this.reconciliationService.getEnvironmentAuditList(id);
  }

  @Sse('components/notify/:id')
  notify(@Param('id') id: string): Observable<MessageEvent> {
    return new Observable((observer: Observer<MessageEvent>) => {
      this.reconciliationService.notifyStream.subscribe(
        async (component: ComponentReconcile) => {
          if (component.name !== id) {
            return;
          }
          const data: ComponentAudit[] = Mapper.getComponentAuditList([component]);
          observer.next({
            data: data[0],
          })
        },
      )
    });
  }
}

export interface MessageEvent {
  data: string | object
  id?: string
  type?: string
  retry?: number
}

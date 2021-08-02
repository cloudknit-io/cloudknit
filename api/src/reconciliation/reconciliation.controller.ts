import { Body, Controller, Get, Param, Post, Sse } from '@nestjs/common'
import { Observable, Observer } from 'rxjs';
import { Mapper } from 'src/costing/utilities/mapper';
import { S3Handler } from 'src/costing/utilities/s3Handler';
import { ComponentReconcile } from 'src/typeorm/reconciliation/component-reconcile.entity';
import { EnvironmentReconcile } from 'src/typeorm/reconciliation/environment-reconcile.entity';
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

  @Get('component/logs/:team/:environment/:component/:id')
  async getLogs(@Param('team') team: string, @Param('environment') environment: string, @Param('component') component: string, @Param('id') id: number) {
    return await this.reconciliationService.getLogs(team, environment, component, id);
  }

  @Get('component/latestLogs/:team/:environment/:component')
  async getLatestLogs(@Param('team') team: string, @Param('environment') environment: string, @Param('component') component: string) {
    return await this.reconciliationService.getLatestLogs(team, environment, component);
  }

  @Get('component/plan/logs/:team/:environment/:component/:id/:latest')
  async getPlanLogs(@Param('team') team: string, @Param('environment') environment: string, @Param('component') component: string, @Param('id') id: number, @Param('latest') latest: boolean) {
    return await this.reconciliationService.getPlanLogs(team, environment, component, id, latest);
  }

  @Get('component/apply/logs/:team/:environment/:component/:id/:latest')
  async getApplyLogs(@Param('team') team: string, @Param('environment') environment: string, @Param('component') component: string, @Param('id') id: number, @Param('latest') latest: boolean) {
    return await this.reconciliationService.getApplyLogs(team, environment, component, id, latest);
  }



  @Sse('components/notify/:id')
  notifyComponents(@Param('id') id: string): Observable<MessageEvent> {
    return new Observable((observer: Observer<MessageEvent>) => {
      this.reconciliationService.notifyStream.subscribe(
        async (component: ComponentReconcile) => {
          if (component.name !== id) {
            return;
          }
          const data: ComponentAudit[] = Mapper.getComponentAuditList([component]);
          observer.next({
            data: data,
          })
        },
      )
    });
  }

  @Sse('environments/notify/:id')
  notifyEnvironments(@Param('id') id: string): Observable<MessageEvent> {
    return new Observable((observer: Observer<MessageEvent>) => {
      this.reconciliationService.notifyStream.subscribe(
        async (environment: EnvironmentReconcile) => {
          if (environment.name !== id) {
            return;
          }
          const data: EnvironmentAudit[] = Mapper.getEnvironmentAuditList([environment]);
          observer.next({
            data: data,
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

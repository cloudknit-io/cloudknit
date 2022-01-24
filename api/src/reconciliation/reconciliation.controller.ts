import { Body, Controller, Get, Param, Patch, Post, Res, Sse, UploadedFile, UseInterceptors } from "@nestjs/common";
import { FileInterceptor } from "@nestjs/platform-express";
import { response } from "express";
import { from, Observable, Observer } from "rxjs";
import { map } from "rxjs/operators";
import { Mapper } from "src/costing/utilities/mapper";
import { ComponentReconcile } from "src/typeorm/reconciliation/component-reconcile.entity";
import { EnvironmentReconcile } from "src/typeorm/reconciliation/environment-reconcile.entity";
import { Notification } from "src/typeorm/reconciliation/notification.entity";
import { ComponentAudit } from "./dtos/componentAudit.dto";
import { EnvironmentAudit } from "./dtos/environmentAudit.dto";
import { NotificationDto } from "./dtos/notification.dto";
import { EvnironmentReconcileDto } from "./dtos/reconcile.Dto";
import { ReconciliationService } from "./services/reconciliation.service";

@Controller({
  path: "reconciliation/api/v1",
})
export class ReconciliationController {
  constructor(private readonly reconciliationService: ReconciliationService) {}

  @Post("environment/save")
  async saveEnvironment(@Body() runData: EvnironmentReconcileDto) {
    return await this.reconciliationService.saveOrUpdateEnvironment(runData);
  }

  @Post("notification/save")
  async saveNotification(@Body() notification: NotificationDto) {
    return await this.reconciliationService.saveNotification(notification);
  }

  @Post("component/save")
  async saveComponent(@Body() runData: EvnironmentReconcileDto) {
    return await this.reconciliationService.saveOrUpdateComponent(runData);
  }

  
  @Post("component/putObject")
  @UseInterceptors(FileInterceptor('file'))
  async putObject(@UploadedFile() file: Express.Multer.File, @Body() body: any) {
    console.log(file, body);
    return await this.reconciliationService.putObject(body.customerId, body.path, file);
  }

  @Post("component/downloadObject")
  async downloadObject(@Res() response, @Body() body: any) {
    console.log(body);
    const stream = await this.reconciliationService.downloadObject(body.customerId, body.path);
    stream.pipe(response);
  }

  @Patch("component/update")
  async patchComponent(@Body() runData: any) {
    return await this.reconciliationService.saveOrUpdateComponent(runData);
  }

  @Get("component/:id")
  async getComponents(@Param("id") id: string): Promise<ComponentAudit[]> {
    return await this.reconciliationService.getComponentAuditList(id);
  }

  @Get("notification/seen/:id")
  async setSeenStatusForNotification(@Param("id") id: string): Promise<void> {
    return await this.reconciliationService.setSeenStatusForNotification(
      parseInt(id)
    );
  }

  @Get("environment/:id")
  async getEnvironments(@Param("id") id: string): Promise<EnvironmentAudit[]> {
    return await this.reconciliationService.getEnvironmentAuditList(id);
  }

  @Get("component/logs/:companyId/:team/:environment/:component/:id")
  async getLogs(
    @Param("companyId") companyId: string,
    @Param("team") team: string,
    @Param("environment") environment: string,
    @Param("component") component: string,
    @Param("id") id: number
  ) {
    return await this.reconciliationService.getLogs(
      companyId,
      team,
      environment,
      component,
      id
    );
  }

  @Get("component/latestLogs/:companyId/:team/:environment/:component")
  async getLatestLogs(
    @Param("companyId") companyId: string,
    @Param("team") team: string,
    @Param("environment") environment: string,
    @Param("component") component: string
  ) {
    return await this.reconciliationService.getLatestLogs(
      companyId,
      team,
      environment,
      component
    );
  }

  @Get("component/state-file/:companyId/:team/:environment/:component")
  async getStateFile(
    @Param("companyId") companyId: string,
    @Param("team") team: string,
    @Param("environment") environment: string,
    @Param("component") component: string
  ) {
    return await this.reconciliationService.getStateFile(
      companyId,
      team,
      environment,
      component
    );
  }

  @Get(
    "component/plan/logs/:companyId/:team/:environment/:component/:id/:latest"
  )
  async getPlanLogs(
    @Param("companyId") companyId: string,
    @Param("team") team: string,
    @Param("environment") environment: string,
    @Param("component") component: string,
    @Param("id") id: number,
    @Param("latest") latest: string
  ) {
    return await this.reconciliationService.getPlanLogs(
      companyId,
      team,
      environment,
      component,
      id,
      latest === "true"
    );
  }

  @Get(
    "component/apply/logs/:companyId/:team/:environment/:component/:id/:latest"
  )
  async getApplyLogs(
    @Param("companyId") companyId: string,
    @Param("team") team: string,
    @Param("environment") environment: string,
    @Param("component") component: string,
    @Param("id") id: number,
    @Param("latest") latest: string
  ) {
    return await this.reconciliationService.getApplyLogs(
      companyId,
      team,
      environment,
      component,
      id,
      latest === "true"
    );
  }

  @Get("notifications/get/:companyId/:teamName")
  async getNotifications(
    @Param("companyId") companyId: string,
    @Param("teamName") teamName: string
  ) {
    await this.reconciliationService.getNotification(companyId, teamName);
    return await this.reconciliationService.getAllNotification(
      companyId,
      teamName
    );
  }

  @Sse("notifications/:companyId/:teamName")
  sendNotification(
    @Param("companyId") companyId: string,
    @Param("teamName") teamName: string
  ): Observable<MessageEvent> {
    return new Observable((observer: Observer<MessageEvent>) => {
      this.reconciliationService.notificationStream.subscribe(
        async (notification: Notification) => {
          if (
            notification.company_id !== companyId ||
            notification.team_name !== teamName
          ) {
            return;
          }
          observer.next({
            data: notification,
          });
        }
      );
    });
  }

  @Sse("components/notify/:id")
  notifyComponents(@Param("id") id: string): Observable<MessageEvent> {
    return from(this.reconciliationService.notifyStream).pipe(
      map((component: ComponentReconcile) => {
        if (component.name !== id) {
          return { data: [] };
        }
        const data: ComponentAudit[] = Mapper.getComponentAuditList([
          component,
        ]);
        return { data };
      })
    );
  }

  @Sse("environments/notify/:id")
  notifyEnvironments(@Param("id") id: string): Observable<MessageEvent> {
    return from(this.reconciliationService.notifyStream).pipe(
      map((environment: EnvironmentReconcile) => {
        if (environment.name !== id) {
          return { data: [] };
        }
        const data: EnvironmentAudit[] = Mapper.getEnvironmentAuditList([
          environment,
        ]);
        return { data };
      })
    );
  }
}

export interface MessageEvent {
  data: string | object;
  id?: string;
  type?: string;
  retry?: number;
}

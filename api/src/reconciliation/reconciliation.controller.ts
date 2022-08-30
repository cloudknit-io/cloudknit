import {
  Body,
  Controller,
  Get,
  Param,
  Patch,
  Post,
  Put,
  Req,
  Request,
  Res,
  Sse,
  UploadedFile,
  UseInterceptors,
} from "@nestjs/common";
import { FileInterceptor } from "@nestjs/platform-express";
import { from, Observable } from "rxjs";
import { map } from "rxjs/operators";
import { Mapper } from "src/costing/utilities/mapper";
import { ComponentReconcile } from "src/typeorm/reconciliation/component-reconcile.entity";
import { EnvironmentReconcile } from "src/typeorm/reconciliation/environment-reconcile.entity";
import { ComponentDto } from "./dtos/component.dto";
import { ComponentAudit } from "./dtos/componentAudit.dto";
import { EnvironmentDto } from "./dtos/environment.dto";
import { EnvironmentAudit } from "./dtos/environmentAudit.dto";
import { EvnironmentReconcileDto } from "./dtos/reconcile.Dto";
import { ReconciliationService } from "./reconciliation.service";

@Controller({
  version: '1'
})
export class ReconciliationController {
  constructor(private readonly reconciliationService: ReconciliationService) {}

  @Put("environments")
  async putEnvironment(@Request() req, @Body() body: EnvironmentDto) {
    return await this.reconciliationService.putEnvironment(req.org, body);
  }

  @Put("components")
  async putComponent(@Request() req, @Body() body: ComponentDto) {
    return await this.reconciliationService.putComponent(req.org, body);
  }

  @Get("environments/:id")
  async getEnvironment(@Request() req, @Param("id") id: string) {
    return await this.reconciliationService.getEnvironment(req.org, id);
  }

  @Get("components/:id")
  async getComponent(@Request() req, @Param("id") id: string) {
    return await this.reconciliationService.getComponent(req.org, id);
  }

  @Patch("approved-by/:id")
  async patchApprovedBy(@Param("id") id: string, @Req() req) {
    return await this.reconciliationService.patchApprovedBy(req.org, req.headers['x-user-email'] || '', id);
  }

  @Get("approved-by/:id/:rid") 
  async getApprovedBy(@Request() req, @Param("id") id: string, @Param("rid") rid: string) {
    return await this.reconciliationService.getApprovedBy(req.org, id, rid);
  }

  @Post("environment/save")
  async saveEnvironment(@Request() req, @Body() runData: EvnironmentReconcileDto) {
    return await this.reconciliationService.saveOrUpdateEnvironment(req.org, runData);
  }

  @Post("component/save")
  async saveComponent(@Request() req, @Body() runData: EvnironmentReconcileDto) {
    return await this.reconciliationService.saveOrUpdateComponent(req.org, runData);
  }

  @Post("component/putObject")
  @UseInterceptors(FileInterceptor("file"))
  async putObject(
    @UploadedFile() file: Express.Multer.File,
    @Body() body: any
  ) {
    // QUESTION : Is this still in use? It does not exist in BFF
    return {
      message: await this.reconciliationService.putObject(
        body.customerId,
        body.path,
        file
      ),
    };
  }

  @Post("component/downloadObject")
  async downloadObject(@Res() response, @Body() body: any) {
    // QUESTION : Is this still in use? It does not exist in BFF
    const stream = await this.reconciliationService.downloadObject(
      body.customerId,
      body.path
    );
    stream.pipe(response);
  }

  @Patch("component/update")
  async patchComponent(@Request() req, @Body() runData: any) {
    // QUESTION : Is this still in use? It does not exist in BFF
    return await this.reconciliationService.saveOrUpdateComponent(req.org, runData);
  }

  @Get("component/:id")
  async getComponents(@Request() req, @Param("id") id: string): Promise<ComponentAudit[]> {
    return await this.reconciliationService.getComponentAuditList(req.org, id);
  }

  @Get("environment/:id")
  async getEnvironments(@Request() req, @Param("id") id: string): Promise<EnvironmentAudit[]> {
    return await this.reconciliationService.getEnvironmentAuditList(req.org, id);
  }

  @Get("component/logs/:team/:environment/:component/:id")
  async getLogs(
    @Request() req,
    @Param("team") team: string,
    @Param("environment") environment: string,
    @Param("component") component: string,
    @Param("id") id: number
  ) {
    return await this.reconciliationService.getLogs(
      req.org,
      team,
      environment,
      component,
      id
    );
  }

  @Get("component/latestLogs/:team/:environment/:component")
  async getLatestLogs(
    @Request() req,
    @Param("team") team: string,
    @Param("environment") environment: string,
    @Param("component") component: string
  ) {
    return await this.reconciliationService.getLatestLogs(
      req.org,
      team,
      environment,
      component
    );
  }

  @Get("component/state-file/:team/:environment/:component")
  async getStateFile(
    @Request() req,
    @Param("team") team: string,
    @Param("environment") environment: string,
    @Param("component") component: string
  ) {
    return await this.reconciliationService.getStateFile(
      req.org,
      team,
      environment,
      component
    );
  }

  @Get(
    "component/plan/logs/:team/:environment/:component/:id/:latest"
  )
  async getPlanLogs(
    @Request() req,
    @Param("team") team: string,
    @Param("environment") environment: string,
    @Param("component") component: string,
    @Param("id") id: number,
    @Param("latest") latest: string
  ) {
    return await this.reconciliationService.getPlanLogs(
      req.org,
      team,
      environment,
      component,
      id,
      latest === "true"
    );
  }

  @Get(
    "component/apply/logs/:team/:environment/:component/:id/:latest"
  )
  async getApplyLogs(
    @Request() req,
    @Param("team") team: string,
    @Param("environment") environment: string,
    @Param("component") component: string,
    @Param("id") id: number,
    @Param("latest") latest: string
  ) {
    return await this.reconciliationService.getApplyLogs(
      req.org,
      team,
      environment,
      component,
      id,
      latest === "true"
    );
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

  @Sse("applications/:id")
  notifyApplications(@Param("id") id: string): Observable<MessageEvent> {
    return from(this.reconciliationService.applicationStream).pipe(
      map((application: any) => ({ data: application }))
    );
  }
}

export interface MessageEvent {
  data: string | object;
  id?: string;
  type?: string;
  retry?: number;
}

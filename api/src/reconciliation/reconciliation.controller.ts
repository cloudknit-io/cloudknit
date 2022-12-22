import {
  BadRequestException,
  Body,
  Controller,
  Get,
  Param,
  Patch,
  Post,
  Query,
  Req,
  Request,
  Sse,
} from "@nestjs/common";
import { from, Observable } from "rxjs";
import { map } from "rxjs/operators";
import { Mapper } from "src/costing/utilities/mapper";
import { ComponentReconcile } from "src/typeorm/reconciliation/component-reconcile.entity";
import { EnvironmentReconcile } from "src/typeorm/reconciliation/environment-reconcile.entity";
import { APIRequest } from "src/types";
import { ApprovedByDto, ComponentAudit } from "./dtos/componentAudit.dto";
import { EnvironmentAudit } from "./dtos/environmentAudit.dto";
import { EvnironmentReconcileDto } from "./dtos/reconcile.Dto";
import { EnvironmentService } from "./environment.service";
import { ReconciliationService } from "./reconciliation.service";
import { SSEService } from "./sse.service";

@Controller({
  version: '1'
})
export class ReconciliationController {
  constructor(
    private readonly reconciliationService: ReconciliationService,
    private readonly envSvc: EnvironmentService,
    private readonly sseSvc: SSEService
    ) {}

  @Get("environments")
  async getEnvironment(@Request() req, @Query("envName") envName: string, @Query("teamName") teamName: string) {
    const env = await this.envSvc.getEnvironment(req.org, envName, teamName);

    if (!env) {
      throw new BadRequestException('could not find environment');
    }

    return env;
  }

  @Get("components/:id")
  async getComponent(@Request() req, @Param("id") id: string) {
    const comp = await this.reconciliationService.getComponent(req.org, id);

    if (!comp) {
      throw new BadRequestException('could not find component');
    }

    return comp;
  }

  @Patch("approved-by")
  async patchApprovedBy(@Req() req: APIRequest, @Body() body: ApprovedByDto) {
    return await this.reconciliationService.patchApprovedBy(req.org, body);
  }

  @Get("approved-by")
  async getApprovedBy(@Request() req, @Query("compName") compName: string, @Query("envName") envName: string, @Query("teamName") teamName: string, @Query("rid") rid: number) {
    const body: ApprovedByDto = {
      compName,
      envName,
      teamName,
      rid,
    }
    return await this.reconciliationService.getApprovedBy(req.org, body);
  }

  @Post("environment/save")
  async saveEnvironment(@Request() req, @Body() runData: EvnironmentReconcileDto) {
    return await this.reconciliationService.saveOrUpdateEnvironment(req.org, runData);
  }

  @Post("component/save")
  async saveComponent(@Request() req, @Body() runData: EvnironmentReconcileDto) {
    return await this.reconciliationService.saveOrUpdateComponent(req.org, runData);
  }

  @Get("audit/components")
  async getComponents(@Request() req, @Query("compName") compName: string, @Query("envName") envName: string, @Query("teamName") teamName: string): Promise<ComponentAudit[]> {
    return await this.reconciliationService.getComponentAuditList(req.org, compName, envName, teamName);
  }

  @Get("audit/environments")
  async getEnvironments(@Request() req, @Query("envName") envName: string, @Query("teamName") teamName: string): Promise<EnvironmentAudit[]> {
    return await this.reconciliationService.getEnvironmentAuditList(req.org, envName, teamName);
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
    return from(this.sseSvc.notifyStream).pipe(
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
    return from(this.sseSvc.notifyStream).pipe(
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
    return from(this.sseSvc.applicationStream).pipe(
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

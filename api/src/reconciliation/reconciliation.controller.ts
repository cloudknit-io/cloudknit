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
import { ComponentReconcile } from "src/typeorm/component-reconcile.entity";
import { EnvironmentReconcile } from "src/typeorm/environment-reconcile.entity";
import { APIRequest } from "src/types";
import { ApprovedByDto, ComponentAudit } from "./dtos/componentAudit.dto";
import { EnvironmentAudit } from "./dtos/environmentAudit.dto";
import { EvnironmentReconcileDto } from "./dtos/reconcile.Dto";
import { EnvironmentService } from "./environment.service";
import { ReconciliationService } from "./reconciliation.service";
import { SSEService } from "./sse.service";
import { RequiredQueryValidationPipe, TeamEnvCompQueryParams, TeamEnvQueryParams } from "./validationPipes";


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
  async getEnvironment(@Request() req, @Query(new RequiredQueryValidationPipe()) tec: TeamEnvQueryParams) {
    const env = await this.envSvc.getEnvironment(req.org, tec.envName, tec.teamName);

    if (!env) {
      throw new BadRequestException('could not find environment');
    }

    return env;
  }

  @Get("components")
  async getComponent(@Request() req, @Query(new RequiredQueryValidationPipe()) tec: TeamEnvCompQueryParams) {
    const comp = await this.reconciliationService.getComponent(req.org, tec.compName, tec.envName, tec.teamName);

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
  async getApprovedBy(@Request() req, @Query(new RequiredQueryValidationPipe()) tec: TeamEnvCompQueryParams, @Query("rid") rid: number) {
    const body: ApprovedByDto = {
      compName: tec.compName,
      envName: tec.envName,
      teamName: tec.teamName,
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
  async getComponents(@Request() req, @Query(new RequiredQueryValidationPipe()) tec: TeamEnvCompQueryParams): Promise<ComponentAudit[]> {
    return await this.reconciliationService.getComponentAuditList(req.org, tec.compName, tec.envName, tec.teamName);
  }

  @Get("audit/environments")
  async getEnvironments(@Request() req, @Query(new RequiredQueryValidationPipe()) te: TeamEnvQueryParams): Promise<EnvironmentAudit[]> {
    return await this.reconciliationService.getEnvironmentAuditList(req.org, te.envName, te.teamName);
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

  @Sse("components/notify")
  notifyComponents(@Request() req, @Query(new RequiredQueryValidationPipe()) tec: TeamEnvCompQueryParams): Observable<MessageEvent> {
    return from(this.sseSvc.notifyStream).pipe(
      map((component: ComponentReconcile) => {
        if (component.name !== tec.compName || tec.envName !== component.environmentReconcile.name || tec.teamName !== component.environmentReconcile.team_name || component.organization.id !== req.org.id) {
          return { data: [] };
        }
        const data: ComponentAudit[] = Mapper.getComponentAuditList([
          component,
        ]);
        return { data };
      })
    );
  }

  @Sse("environments/notify")
  notifyEnvironments(@Request() req, @Query(new RequiredQueryValidationPipe()) te: TeamEnvQueryParams): Observable<MessageEvent> {
    return from(this.sseSvc.notifyStream).pipe(
      map((environment: EnvironmentReconcile) => {
        if (environment.name !== te.envName || te.teamName != environment.team_name || req.org.id !== environment.organization.id) {
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

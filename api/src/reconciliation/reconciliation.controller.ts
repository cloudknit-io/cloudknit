import {
  BadRequestException,
  Body,
  Controller,
  Get,
  Logger,
  Param,
  Patch,
  Query,
  Req,
  Request,
  Sse,
} from "@nestjs/common";
import { from, Observable } from "rxjs";
import { map } from "rxjs/operators";
import { ComponentService } from "src/costing/services/component.service";
import { Mapper } from "src/costing/utilities/mapper";
import { EnvironmentService } from "src/environment/environment.service";
import { RootEnvironmentService } from "src/root-environment/root.environment.service";
import { TeamService } from "src/team/team.service";
import { ComponentReconcile, EnvironmentReconcile } from "src/typeorm";
import { APIRequest } from "src/types";
import { ApprovedByDto, ComponentAudit } from "./dtos/componentAudit.dto";
import { EnvironmentAudit } from "./dtos/environmentAudit.dto";
import { EvnironmentReconcileDto } from "./dtos/reconcile.Dto";
import { ReconciliationService } from "./reconciliation.service";
import { SSEService } from "./sse.service";
import { RequiredQueryValidationPipe, TeamEnvCompQueryParams, TeamEnvQueryParams } from "./validationPipes";


@Controller({
  version: '1'
})
export class ReconciliationController {
  private readonly logger = new Logger(ReconciliationController.name);

  constructor(
    private readonly reconciliationService: ReconciliationService,
    private readonly rootEnvService: RootEnvironmentService,
    private readonly envSvc: EnvironmentService,
    private readonly sseSvc: SSEService,
    private readonly teamSvc: TeamService,
    private readonly compSvc: ComponentService
    ) {}

  @Get("environments")
  async getEnvironment(@Request() req, @Query(new RequiredQueryValidationPipe()) tec: TeamEnvQueryParams) {
    const {org, team } = req;
    const env = await this.envSvc.findByName(req.org, team, tec.envName);

    if (!env) {
      throw new BadRequestException('could not find environment');
    }

    return env;
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

  @Get("audit/components")
  async getComponents(@Request() req, @Query(new RequiredQueryValidationPipe()) tec: TeamEnvCompQueryParams): Promise<ComponentAudit[]> {
    const {org, team} = req;
    return await this.reconciliationService.getComponentAuditList(org, team, tec.compName, tec.envName);
  }

  @Get("audit/environments")
  async getEnvironments(@Request() req, @Query(new RequiredQueryValidationPipe()) te: TeamEnvQueryParams): Promise<EnvironmentAudit[]> {
    const {org, team} = req;
    return await this.reconciliationService.getEnvironmentAuditList(org, team, te.envName);
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

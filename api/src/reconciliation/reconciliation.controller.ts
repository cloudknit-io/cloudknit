import {
  BadRequestException,
  Body,
  ConsoleLogger,
  Controller,
  Get,
  InternalServerErrorException,
  Logger,
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
import { ComponentService } from "src/component/component.service";
import { Mapper } from "src/costing/utilities/mapper";
import { EnvironmentService } from "src/environment/environment.service";
import { TeamService } from "src/team/team.service";
import { ComponentReconcile, EnvironmentReconcile } from "src/typeorm";
import { APIRequest } from "src/types";
import { handleSqlErrors } from "src/utilities/errorHandler";
import { ApprovedByDto, ComponentAudit } from "./dtos/componentAudit.dto";
import { EnvironmentAudit } from "./dtos/environmentAudit.dto";
import { CreateComponentReconciliationDto, CreateEnvironmentReconciliationDto, UpdateComponentReconciliationDto, UpdateEnvironmentReconciliationDto } from "./dtos/reconciliation.dto";
import { ReconciliationService } from "./reconciliation.service";
import { RequiredQueryValidationPipe, TeamEnvCompQueryParams, TeamEnvQueryParams } from "./validationPipes";


@Controller({
  version: '1'
})
export class ReconciliationController {
  private readonly logger = new Logger(ReconciliationController.name);

  constructor(
    private readonly reconSvc: ReconciliationService,
    private readonly envSvc: EnvironmentService,
    private readonly teamSvc: TeamService,
    private readonly compSvc: ComponentService
    ) {}

  @Post('environment')
  async newEnvironmentReconciliation(@Req() req: APIRequest, @Body() body: CreateEnvironmentReconciliationDto) {
    const { org } = req;

    const team = await this.teamSvc.findByName(org, body.teamName);
    if (!team) {
      this.logger.error({ message: 'could not find team in newEnvironmentReconciliation', body});
      throw new BadRequestException('could not find team');
    }

    const env = await this.envSvc.findByName(org, team, body.name);
    if (!env) {
      this.logger.error({ message: 'could not find environment in newEnvironmentReconciliation', body});
      throw new BadRequestException('could not find environment');
    }

    let envReconEntry: EnvironmentReconcile;

    try {
      envReconEntry = await this.reconSvc.createEnvRecon(org, team, env, body);
    } catch(err) {
      handleSqlErrors(err);

      this.logger.error({ message: 'could not create environment recon', body, err });
      throw new InternalServerErrorException('could not create environment reconciliation');
    }

    try {
      const skippedEntries = await this.reconSvc.getSkippedEnvironments(org, team, env, [envReconEntry.reconcileId]);
      await this.reconSvc.bulkUpdateEnvironmentEntries(skippedEntries, 'skipped_reconcile');
    } catch (err) {
      this.logger.error('could not update skipped environment reconciles', err);
      throw new InternalServerErrorException();
    }

    return envReconEntry.reconcileId;
  }

  @Post('environment/:reconcileId')
  async updateEnvironmentReconciliation(@Req() req: APIRequest, @Param('reconcileId') reconcileId: number, @Body() body: UpdateEnvironmentReconciliationDto) {
    const { org } = req;

    const existingEntry = await this.reconSvc.getEnvReconByReconcileId(org, reconcileId, true);
    if (!existingEntry) {
      this.logger.error({message: 'could not find environment reconcile entry', reconcileId })
      throw new BadRequestException('could not find environment reconcile entry');
    }
    
    const env = existingEntry.environment;
    let envRecon: EnvironmentReconcile;

    try {
      envRecon = await this.reconSvc.updateEnvRecon(existingEntry, body);
    } catch (err) {
      handleSqlErrors(err);

      this.logger.error({ message: 'could not update env recon', reconcileId, existingEntry, body });
      throw new InternalServerErrorException('could not update environment reconcile');
    }

    let duration = -1;

    if (envRecon.endDateTime) {
      const ed = new Date(envRecon.endDateTime).getTime();
      const sd = new Date(envRecon.startDateTime).getTime();
      duration = ed - sd;
    }

    try {
      await this.envSvc.update(org, env.id, {
        duration,
        status: body.status
      });
    } catch (err) {
      handleSqlErrors(err);

      this.logger.error({ message: 'could not update environment with env recon', envRecon, body, duration, err});
      throw new InternalServerErrorException('could not update environment');
    }

    return envRecon;
  }

  @Post('component')
  async newComponentReconciliation(@Req() req: APIRequest, @Body() body: CreateComponentReconciliationDto): Promise<number> {
    const { org } = req;

    const envRecon = await this.reconSvc.getEnvReconByReconcileId(org, body.envReconcileId, true);
    if (!envRecon) {
      this.logger.error({ message: 'could not find environment-reconcile in newComponentReconciliation', body });
      throw new BadRequestException('could not find environment-reconcile');
    }

    let compRecon: ComponentReconcile;

    try {
      compRecon = await this.reconSvc.createCompRecon(org, envRecon, body);
    } catch (err) {
      handleSqlErrors(err);
      
      this.logger.error({ message: 'could not save component-reconcile in newComponentReconciliation', body });
      throw new BadRequestException('could not save component-reconcile');
    }

    try {
      const skippedEntries = await this.reconSvc.getSkippedComponents(org, envRecon.environment, [compRecon.reconcileId])
      await this.reconSvc.bulkUpdateComponentEntries(skippedEntries, 'skipped_reconcile');
    } catch (err) {
      this.logger.error('could not update skipped component reconciles', err);
      throw new InternalServerErrorException();
    }

    return compRecon.reconcileId;
  }

  @Post('component/:reconcileId')
  async updateComponentReconciliation(@Req() req: APIRequest, @Param('reconcileId') compReconcileId: number, @Body() body: UpdateComponentReconciliationDto) {
    const { org } = req;

    const compRecon: ComponentReconcile = await this.reconSvc.getCompReconById(org, compReconcileId, true);
    if (!compRecon) {
      this.logger.error({ message: 'could not find component-reconcile in updateComponentReconciliation', body });
      throw new BadRequestException('could not find component-reconcile');
    }
    
    const envRecon = compRecon.environmentReconcile;

    const env = await this.envSvc.findById(org, envRecon.environment.id);
    if (!env) {
      this.logger.error({ message: 'could not find environment in updateComponentReconciliation', body });
      throw new BadRequestException('could not find environment');
    }
    
    const comp = await this.compSvc.findByName(org, env, compRecon.name);
    if (!comp) {
      this.logger.error({ message: 'could not find component in updateComponentReconciliation', body });
      throw new BadRequestException('could not find component');
    }

    const updatedCompRecon = await this.reconSvc.updateCompRecon(compRecon, body);
    delete updatedCompRecon.environmentReconcile;
    this.logger.log({message: 'updated component reconcile entry', updatedCompRecon});
    
    let duration = comp.duration;
    if (updatedCompRecon.endDateTime) {
      const ed = new Date(body.endDateTime).getTime();
      const sd = new Date(compRecon.startDateTime).getTime();
      duration = ed - sd;
    }

    await this.compSvc.updateFromCompRecon(comp, {
      duration,
      status: updatedCompRecon.status
    });
    this.logger.log({ message: 'updated component', comp});

    return updatedCompRecon;
  }

  @Patch("approved-by")
  async patchApprovedBy(@Req() req: APIRequest, @Body() body: ApprovedByDto) {
    return await this.reconSvc.patchApprovedBy(req.org, body);
  }

  @Get("approved-by")
  async getApprovedBy(@Request() req: APIRequest, @Query(new RequiredQueryValidationPipe()) tec: TeamEnvCompQueryParams, @Query("rid") rid: number) {
    const { org } = req;

    if (rid > 0) {
      return this.reconSvc.getCompReconById(org, rid);
    }

    const team = await this.teamSvc.findByName(org, tec.teamName);
    const env = await this.envSvc.findByName(org, team, tec.envName);
    return await this.reconSvc.getCompReconByName(org, env, tec.compName);
  }

  @Get("audit/components")
  async getComponents(@Request() req, @Query(new RequiredQueryValidationPipe()) tec: TeamEnvCompQueryParams): Promise<ComponentAudit[]> {
    const {org, team} = req;
    return await this.reconSvc.getComponentAuditList(org, team, tec.compName, tec.envName);
  }

  @Get("audit/environments")
  async getEnvironments(@Request() req, @Query(new RequiredQueryValidationPipe()) te: TeamEnvQueryParams): Promise<EnvironmentAudit[]> {
    const {org, team} = req;
    return await this.reconSvc.getEnvironmentAuditList(org, team, te.envName);
  }

  @Get("component/logs/:team/:environment/:component/:id")
  async getLogs(
    @Request() req,
    @Param("team") team: string,
    @Param("environment") environment: string,
    @Param("component") component: string,
    @Param("id") id: number
  ) {
    return await this.reconSvc.getLogs(
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
    return await this.reconSvc.getStateFile(
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
    return {};
    // return await this.reconSvc.getPlanLogs(
    //   req.org,
    //   team,
    //   environment,
    //   component,
    //   id,
    //   latest === "true"
    // );
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
    return {};
    // return await this.reconSvc.getApplyLogs(
    //   req.org,
    //   team,
    //   environment,
    //   component,
    //   id,
    //   latest === "true"
    // );
  }
}

export interface MessageEvent {
  data: string | object;
  id?: string;
  type?: string;
  retry?: number;
}

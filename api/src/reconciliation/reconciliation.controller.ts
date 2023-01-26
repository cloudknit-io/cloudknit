import {
  BadRequestException,
  Body,
  Controller,
  Get,
  InternalServerErrorException,
  Logger,
  Param,
  Post,
  Req,
  Request,
} from '@nestjs/common';
import { ApiTags } from '@nestjs/swagger';
import { ApproveWorkflow as ResumeWorkflow } from 'src/argowf/api';
import { ComponentService } from 'src/component/component.service';
import { EnvironmentService } from 'src/environment/environment.service';
import { TeamService } from 'src/team/team.service';
import {
  Component,
  ComponentReconcile,
  Environment,
  EnvironmentReconcile,
  Organization,
  Team,
} from 'src/typeorm';
import { APIRequest, OrgApiParam } from 'src/types';
import { handleSqlErrors } from 'src/utilities/errorHandler';
import { ApprovedByDto } from './dtos/componentAudit.dto';
import {
  CreateComponentReconciliationDto,
  CreateEnvironmentReconciliationDto,
  UpdateComponentReconciliationDto,
  UpdateEnvironmentReconciliationDto,
} from './dtos/reconciliation.dto';
import { ReconciliationService } from './reconciliation.service';

@Controller({
  version: '1',
})
@ApiTags('reconciliation')
export class ReconciliationController {
  private readonly logger = new Logger(ReconciliationController.name);

  constructor(
    private readonly reconSvc: ReconciliationService,
    private readonly envSvc: EnvironmentService,
    private readonly teamSvc: TeamService,
    private readonly compSvc: ComponentService
  ) {}

  @Post('environment')
  @OrgApiParam()
  async newEnvironmentReconciliation(
    @Req() req: APIRequest,
    @Body() body: CreateEnvironmentReconciliationDto
  ) {
    const { org } = req;

    const team = await this.teamSvc.findByName(org, body.teamName);
    if (!team) {
      this.logger.error({
        message: 'could not find team in newEnvironmentReconciliation',
        body,
      });
      throw new BadRequestException('could not find team');
    }

    const env = await this.envSvc.findByName(org, team, body.name);
    if (!env) {
      this.logger.error({
        message: 'could not find environment in newEnvironmentReconciliation',
        body,
      });
      throw new BadRequestException('could not find environment');
    }

    let envReconEntry: EnvironmentReconcile;

    try {
      envReconEntry = await this.reconSvc.createEnvRecon(org, team, env, body);
    } catch (err) {
      handleSqlErrors(err);

      this.logger.error({
        message: 'could not create environment recon',
        body,
        err,
      });
      throw new InternalServerErrorException(
        'could not create environment reconciliation'
      );
    }

    try {
      const skippedEntries = await this.reconSvc.getSkippedEnvironments(
        org,
        team,
        env,
        [envReconEntry.reconcileId]
      );
      await this.reconSvc.bulkUpdateEnvironmentEntries(
        skippedEntries,
        'skipped_reconcile'
      );
    } catch (err) {
      this.logger.error('could not update skipped environment reconciles', err);
      throw new InternalServerErrorException();
    }

    return envReconEntry.reconcileId;
  }

  @Post('environment/:reconcileId')
  @OrgApiParam()
  async updateEnvironmentReconciliation(
    @Req() req: APIRequest,
    @Param('reconcileId') reconcileId: number,
    @Body() body: UpdateEnvironmentReconciliationDto
  ) {
    const { org } = req;

    const existingEntry = await this.reconSvc.getEnvReconByReconcileId(
      org,
      reconcileId,
      false
    );
    if (!existingEntry) {
      this.logger.error({
        message: 'could not find environment reconcile entry',
        reconcileId,
      });
      throw new BadRequestException(
        'could not find environment reconcile entry'
      );
    }

    let envRecon: EnvironmentReconcile;

    try {
      envRecon = await this.reconSvc.updateEnvRecon(existingEntry, body);
    } catch (err) {
      handleSqlErrors(err);

      this.logger.error({
        message: 'could not update env recon',
        reconcileId,
        existingEntry,
        body,
      });
      throw new InternalServerErrorException(
        'could not update environment reconcile'
      );
    }

    return envRecon;
  }

  @Post('component')
  @OrgApiParam()
  async newComponentReconciliation(
    @Req() req: APIRequest,
    @Body() body: CreateComponentReconciliationDto
  ): Promise<number> {
    const { org } = req;

    const envRecon = await this.reconSvc.getEnvReconByReconcileId(
      org,
      body.envReconcileId,
      true
    );
    if (!envRecon) {
      this.logger.error({
        message:
          'could not find environment-reconcile in newComponentReconciliation',
        body,
      });
      throw new BadRequestException('could not find environment-reconcile');
    }

    const comp = await this.compSvc.findByName(
      org,
      envRecon.environment,
      body.name
    );
    let compRecon: ComponentReconcile;

    try {
      compRecon = await this.reconSvc.createCompRecon(
        org,
        envRecon,
        comp,
        body
      );
    } catch (err) {
      handleSqlErrors(err);

      this.logger.error({
        message:
          'could not save component-reconcile in newComponentReconciliation',
        body,
      });
      throw new BadRequestException('could not save component-reconcile');
    }

    try {
      const skippedEntries = await this.reconSvc.getSkippedComponents(
        org,
        envRecon,
        comp,
        [compRecon.reconcileId]
      );
      await this.reconSvc.bulkUpdateComponentEntries(
        skippedEntries,
        'skipped_reconcile'
      );
    } catch (err) {
      this.logger.error('could not update skipped component reconciles', err);
      throw new InternalServerErrorException();
    }

    return compRecon.reconcileId;
  }

  @Post('component/:reconcileId')
  @OrgApiParam()
  async updateComponentReconciliation(
    @Req() req: APIRequest,
    @Param('reconcileId') compReconcileId: number,
    @Body() body: UpdateComponentReconciliationDto
  ) {
    const { org } = req;

    const compRecon: ComponentReconcile = await this.reconSvc.findCompReconById(
      org,
      compReconcileId,
      true
    );
    if (!compRecon) {
      this.logger.error({
        message:
          'could not find component-reconcile in updateComponentReconciliation',
        body,
      });
      throw new BadRequestException('could not find component-reconcile');
    }

    const updatedCompRecon = await this.reconSvc.updateCompRecon(
      compRecon,
      body
    );
    delete updatedCompRecon.environmentReconcile;
    this.logger.log({
      message: 'updated component reconcile entry',
      updatedCompRecon,
    });
    return updatedCompRecon;
  }

  /**
   * Resumes Argo Workflow run and sets approved by user
   * @param req APIRequest
   * @param compReconId Component reconcile ID to approve
   * @param body Email of user that issued approval
   */
  @Post('component/:compReconId/approve')
  @OrgApiParam()
  async approveWorkflow(
    @Req() req: APIRequest,
    @Param('compReconId') compReconId: number,
    @Body() body: ApprovedByDto
  ) {
    const { org } = req;

    const compRecon = await this.reconSvc.findCompReconById(org, compReconId);
    if (!compRecon) {
      throw new BadRequestException('could not find component reconcile');
    }

    const lastWorkflowRunId = compRecon.component.lastWorkflowRunId;

    try {
      // Resume Argo Workflow run
      await ResumeWorkflow(org, lastWorkflowRunId);
    } catch (err) {
      this.logger.error({
        message: 'could not approve workflow',
        compRecon,
        lastWorkflowRunId,
        err,
      });
      throw new InternalServerErrorException('could not approve workflow');
    }

    try {
      await this.compSvc.update(org, compRecon.component, {
        status: 'initializing_apply',
      });
    } catch (err) {
      handleSqlErrors(err);

      this.logger.error({
        message: 'could not update component status in approveWorkflow',
        compRecon,
        err,
      });
      throw new InternalServerErrorException('could not approve workflow');
    }

    // Set approved by on reconcile entry
    try {
      return this.reconSvc.updateCompRecon(compRecon, {
        approvedBy: body.email,
      });
    } catch (err) {
      handleSqlErrors(err);

      this.logger.error({
        message: 'could not set approved by in approveWorkflow',
        compRecon,
        body,
        err,
      });
      throw new InternalServerErrorException('could not approve workflow');
    }
  }

  @Get('component/logs/:team/:environment/:component/:id')
  @OrgApiParam()
  async getLogs(
    @Request() req,
    @Param('team') teamName: string,
    @Param('environment') envName: string,
    @Param('component') compName: string,
    @Param('id') id: number
  ) {
    const { org } = req;

    const comp = await this.compSvc.findByNameWithTeamName(
      org,
      teamName,
      envName,
      compName,
      true
    );
    if (!comp) {
      this.logger.error({
        message: 'could not find logs',
        teamName,
        envName,
        compName,
        id,
      });
      throw new BadRequestException('could not find logs');
    }

    return await this.reconSvc.getLogs(
      req.org,
      teamName,
      comp.environment,
      comp,
      id
    );
  }

  @Get('component/state-file/:team/:environment/:component')
  @OrgApiParam()
  async getStateFile(
    @Request() req,
    @Param('team') team: string,
    @Param('environment') environment: string,
    @Param('component') component: string
  ) {
    return await this.reconSvc.getStateFile(
      req.org,
      team,
      environment,
      component
    );
  }

  @Get('component/plan/logs/:team/:environment/:component/:id/:latest')
  @OrgApiParam()
  async getPlanLogs(
    @Request() req: APIRequest,
    @Param('team') teamName: string,
    @Param('environment') envName: string,
    @Param('component') compName: string,
    @Param('id') id: number,
    @Param('latest') latest: string
  ) {
    const { org } = req;

    const comp = await this.compSvc.findByNameWithTeamName(
      org,
      teamName,
      envName,
      compName,
      true
    );
    if (!comp) {
      this.logger.error({
        message: 'could not find plan logs',
        teamName,
        envName,
        compName,
        id,
        latest,
      });
      throw new BadRequestException('could not find logs');
    }

    return this.getTfLogs(
      org,
      teamName,
      comp.environment,
      comp,
      id,
      latest === 'true',
      'plan_output'
    );
  }

  @Get('component/apply/logs/:team/:environment/:component/:id/:latest')
  @OrgApiParam()
  async getApplyLogs(
    @Request() req,
    @Param('team') teamName: string,
    @Param('environment') envName: string,
    @Param('component') compName: string,
    @Param('id') id: number,
    @Param('latest') latest: string
  ) {
    const { org } = req;

    const comp = await this.compSvc.findByNameWithTeamName(
      org,
      teamName,
      envName,
      compName,
      true
    );
    if (!comp) {
      this.logger.error({
        message: 'could not find apply logs',
        teamName,
        envName,
        compName,
        id,
        latest,
      });
      throw new BadRequestException('could not find logs');
    }

    return this.getTfLogs(
      org,
      teamName,
      comp.environment,
      comp,
      id,
      latest === 'true',
      'apply_output'
    );
  }

  @OrgApiParam()
  async getTfLogs(
    org: Organization,
    team: string,
    env: Environment,
    comp: Component,
    id: number,
    latest: boolean,
    logType: string
  ) {
    let logs;

    if (latest) {
      const compRecon = await this.reconSvc.getLatestCompReconcile(org, comp);
      logs = await this.reconSvc.getLatestLogs(org, team, env, comp, compRecon);
    } else {
      logs = await this.reconSvc.getLogs(org, team, env, comp, id);
    }

    if (Array.isArray(logs)) {
      return logs.filter((e) => e.key.includes(logType));
    }

    return logs;
  }

  @Get('/token/:token/environment/:id')
  @OrgApiParam()
  async getAuditStatus(
    @Req() req: APIRequest,
    @Param('token') token: string,
    @Param('id') id: string,
  ) {
    const { org } = req;

    const env = await this.envSvc.findByTokenAndSHA(org, token, id);
    if (!env) {
      this.logger.error({
        message: 'could not find environment while getting env audit status',
      });
      throw new BadRequestException('could not find environment');
    }

    return this.reconSvc.getEnvReconStatusByEnv(org, env);
  }
}

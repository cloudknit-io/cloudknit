import {
  BadRequestException,
  Controller,
  Get,
  Logger,
  Param,
  Query,
  Request
} from '@nestjs/common';
import { OnEvent } from '@nestjs/event-emitter';
import { ApiTags } from '@nestjs/swagger';
import { EnvironmentService } from 'src/environment/environment.service';
import { ReconciliationService } from 'src/reconciliation/reconciliation.service';
import { Component, Environment, Organization } from 'src/typeorm';
import {
  APIRequest,
  ComponentReconcileEntityUpdateEvent,
  EnvironmentApiParam,
  InternalEventType,
} from 'src/types';
import { ComponentQueryParams, ComponentWrap } from './component.dto';
import { ComponentService } from './component.service';

@Controller({
  version: '1',
})
@ApiTags('components')
export class ComponentController {
  private readonly logger = new Logger(ComponentController.name);

  constructor(
    private readonly compSvc: ComponentService,
    private readonly envSvc: EnvironmentService,
    private readonly reconSvc: ReconciliationService
  ) {}

  @Get()
  @EnvironmentApiParam()
  async findAll(
    @Request() req: APIRequest,
    @Query() params: ComponentQueryParams
  ): Promise<ComponentWrap[]> {
    const { org, env } = req;
    const withLastReconcile = params.withLastAuditStatus === 'true';

    if (!withLastReconcile) {
      return this.compSvc.findAll(org, env);
    }

    return this.compSvc.findAllWithLastCompRecon(org, env);
  }

  @Get('/:componentId')
  @EnvironmentApiParam()
  async findOne(
    @Request() req: APIRequest,
    @Param('componentId') id: string
  ): Promise<Component> {
    const { org, env } = req;

    return this.getCompFromRequest(org, env, id);
  }

  async getCompFromRequest(
    org: Organization,
    env: Environment,
    id: any
  ): Promise<Component> {
    let numId, comp;

    try {
      numId = parseInt(id, 10);
    } catch (e) {}

    if (isNaN(numId)) {
      comp = await this.compSvc.findByName(org, env, id);
    } else {
      comp = await this.compSvc.findById(org, numId);
    }

    if (!comp) {
      this.logger.error({ message: 'bad componentId', id });
      throw new BadRequestException('component not found');
    }

    return comp;
  }

  @Get('/:componentId/audit')
  @EnvironmentApiParam()
  async getAudits(
    @Request() req: APIRequest,
    @Param('componentId') id: string
  ) {
    const { org, env } = req;
    const comp = await this.getCompFromRequest(org, env, id);

    return this.reconSvc.getComponentAuditList(org, comp);
  }

  @Get('/:componentId/audit/latest')
  @EnvironmentApiParam()
  async getPreviousAudit(
    @Request() req: APIRequest,
    @Param('componentId') id: string
  ) {
    const { org, env } = req;
    const comp = await this.getCompFromRequest(org, env, id);

    return this.reconSvc.getLatestCompReconcile(org, comp);
  }

  @OnEvent(InternalEventType.ComponentReconcileEntityUpdate, { async: true })
  async compReconEnvUpdateListener(evt: ComponentReconcileEntityUpdateEvent) {
    const compRecon = evt.payload;
    
    let envRecon = compRecon.environmentReconcile;
    if (!envRecon) {
      envRecon = await this.reconSvc.getEnvReconByReconcileId(
        {
          id: compRecon.orgId,
        },
        compRecon.envReconId
      );
    }

    const date = new Date().toISOString();

    await this.compSvc.updateById(envRecon.organization, compRecon.compId, {
      lastReconcileDatetime: date
    });
  }
}

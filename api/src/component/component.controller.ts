import {
  BadRequestException,
  Body,
  Controller,
  Get,
  Logger,
  Param,
  Post,
  Put,
  Query,
  Request,
} from '@nestjs/common';
import { EnvironmentService } from 'src/environment/environment.service';
import { ReconciliationService } from 'src/reconciliation/reconciliation.service';
import { Component, Environment, Organization } from 'src/typeorm';
import { APIRequest, EnvironmentApiParam } from 'src/types';
import { ComponentQueryParams, ComponentWrap } from './component.dto';
import { ComponentService } from './component.service';
import { UpdateComponentDto } from './dto/update-component.dto';

@Controller({
  version: '1',
})
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

  @Put('/:componentId')
  @EnvironmentApiParam()
  async updateComponent(
    @Request() req: APIRequest,
    @Param('componentId') id: string,
    @Body() body: UpdateComponentDto
  ): Promise<Component> {
    const { org, env } = req;

    const comp = await this.getCompFromRequest(org, env, id);
    const updatedComp = await this.compSvc.update(org, comp, body);

    return updatedComp;
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
}

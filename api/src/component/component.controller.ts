import { BadRequestException, Body, Controller, Get, Logger, Param, Post, Put, Request } from '@nestjs/common';
import { EnvironmentService } from 'src/environment/environment.service';
import { ReconciliationService } from 'src/reconciliation/reconciliation.service';
import { Component, Environment, Organization } from 'src/typeorm';
import { APIRequest } from 'src/types';
import { ComponentService } from './component.service';
import { UpdateComponentDto } from './dto/update-component.dto';

@Controller({
  version: '1'
})
export class ComponentController {
  private readonly logger = new Logger(ComponentController.name); 

  constructor(
    private readonly compSvc: ComponentService,
    private readonly envSvc: EnvironmentService,
    private readonly reconSvc: ReconciliationService
  ) {}

  @Get()
  async findAll(@Request() req): Promise<Component[]> {
    const {org, env} = req;
    
    return this.compSvc.findAll(org, env);
  }

  @Get(':id')
  async findOne(@Request() req, @Param('id') id: string): Promise<Component> {
    const {org, env} = req;

    return this.getCompFromRequest(org, env, id);
  }

  @Put(':id')
  async updateComponent(@Request() req: APIRequest, @Param('id') id: string, @Body() body: UpdateComponentDto): Promise<Component> {
    const {org, env} = req;

    const comp = await this.getCompFromRequest(org, env, id);

    const updatedComp = await this.compSvc.update(comp, body);

    if (!isNaN(body.estimatedCost)) {
      try {
        this.envSvc.updateCost(org, env);
      } catch (err) {
        this.logger.error({ message: 'could not update environment cost', env })
      }
    }

    return updatedComp;
  }

  async getCompFromRequest(org: Organization, env: Environment, id: any): Promise<Component> {
    let numId, comp;

    try {
      numId = parseInt(id, 10);
    } catch (e) {}
    
    if (isNaN(numId)) {
      comp = await this.compSvc.findByName(org, env, id)
    } else {
      comp = await this.compSvc.findById(org, numId);
    }

    if (!comp) {
      this.logger.error({ message: 'bad componentId', id });
      throw new BadRequestException('component not found');
    }

    return comp;
  }

  @Get(':id/audit')
  async getAudits(@Request() req: APIRequest, @Param('id') id: string) {
    const { org, env } = req;
    const comp = await this.getCompFromRequest(org, env, id);

    return this.reconSvc.getComponentAuditList(org, comp);
  }
}

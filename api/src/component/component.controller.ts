import { Body, Controller, Get, Param, Post, Put, Request } from '@nestjs/common';
import { Component, Environment, Organization } from 'src/typeorm';
import { APIRequest } from 'src/types';
import { ComponentService } from './component.service';
import { UpdateComponentDto } from './dto/update-component.dto';

@Controller({
  version: '1'
})
export class ComponentController {
  constructor(private readonly compSvc: ComponentService) {}

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

    return this.compSvc.update(comp, body);
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

    return comp;
  }
}

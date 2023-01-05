import { Controller, Get, Param, Request } from '@nestjs/common';
import { Component } from 'src/typeorm';
import { ComponentService } from './component.service';

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

    let numId, comp;

    try {
      numId = parseInt(req.params.environmentId, 10);
    } catch (e) {}
    
    if (isNaN(numId)) {
      comp = await this.compSvc.findByName(org, env, id)
    } else {
      comp = await this.compSvc.findById(org, numId);
    }

    return comp;
  }
}

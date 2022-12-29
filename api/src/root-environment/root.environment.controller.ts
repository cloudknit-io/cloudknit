import { Controller, Get, Post, Body, Patch, Param, Delete, Request, BadRequestException, Logger, InternalServerErrorException } from '@nestjs/common';
import { ComponentService } from 'src/costing/services/component.service';
import { CreateEnvironmentDto } from 'src/environment/dto/create-environment.dto';
import { EnvSpecComponentDto, EnvSpecDto } from 'src/environment/dto/env-spec.dto';
import { EnvironmentService } from 'src/environment/environment.service';
import { SqlErrorCodes } from 'src/types';
import { RootEnvironmentService } from './root.environment.service';

@Controller({
  version: '1'
})
export class RootEnvironmentController {
  private readonly logger = new Logger(RootEnvironmentController.name);

  constructor(
    private readonly rootEnvSvc: RootEnvironmentService,
    private readonly envSvc: EnvironmentService,
    private readonly compSvc: ComponentService
  ) {}

  @Post('/spec')
  async spec(@Request() req, @Body() body: EnvSpecDto) {
    const { org, team } = req;
    
    let env = await this.envSvc.findByName(org, team, body.envName);

    if (!env) {
      env = await this.create(req, {
        name: body.envName,
        organization: org,
        team,
        dag: body.components
      });
      
      // create all new components
    } else {
      const dbComps = await this.compSvc.getAllForEnvironmentById(org, env);
      const existingComponents = body.components.filter(incoming => {
        return dbComps.find(dbComp => incoming.name === dbComp.name)
      });
      const newComponents = body.components.filter(incoming => {
        return !dbComps.find(dbComp => incoming.name === dbComp.name)
      });
      const dag: EnvSpecComponentDto[] = [...existingComponents, ...newComponents].map(val => {
        return {
          name: val.name,
          type: val.type,
          dependsOn: val.dependsOn
        };
      });

      env = await this.envSvc.update(org, env.id, {
        dag,
        name: env.name,
        duration: env.duration,
        isDeleted: env.isDeleted
      })

      // create new components
      // update existing
      // delete missing
    }

    return env;
  }

  @Post()
  async create(@Request() req, @Body() createEnv: CreateEnvironmentDto) {
    const {org, team} = req;

    createEnv.organization = org;
    createEnv.team = team;

    try {
      return await this.rootEnvSvc.create(createEnv);
    } catch (err) {
      if (err.code === SqlErrorCodes.DUP_ENTRY) {
        throw new BadRequestException('environment already exists');
      }

      if (err.code === SqlErrorCodes.NO_DEFAULT) {
        throw new BadRequestException(err.sqlMessage);
      }
      
      this.logger.error({ message: 'could not create environment', createEnv, err });
      throw new InternalServerErrorException('could not create environment');
    }
  }

  @Get()
  async findAll(@Request() req) {
    const {org, team} = req;

    return this.rootEnvSvc.findAll(org, team);
  }
}

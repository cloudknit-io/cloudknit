import { Controller, Get, Post, Body, Patch, Param, Delete, Request, BadRequestException, Logger, InternalServerErrorException } from '@nestjs/common';
import { ComponentService } from 'src/costing/services/component.service';
import { CreateEnvironmentDto } from 'src/environment/dto/create-environment.dto';
import { EnvSpecComponentDto, EnvSpecDto } from 'src/environment/dto/env-spec.dto';
import { EnvironmentService } from 'src/environment/environment.service';
import { Environment, Organization, Team } from 'src/typeorm';
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
      this.createEnv(org, team, {
        name: body.envName,
        dag: body.components
      });
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
  async new(@Request() req, @Body() createEnv: CreateEnvironmentDto) {
    const {org, team} = req;
    
    return this.createEnv(org, team, createEnv);
  }

  @Get()
  async findAll(@Request() req) {
    const {org, team} = req;

    return this.rootEnvSvc.findAll(org, team);
  }

  async createEnv(org: Organization, team: Team, createEnv: CreateEnvironmentDto): Promise<Environment> {
    let env: Environment;

    try {
      env = await this.rootEnvSvc.create(org, team, createEnv);
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

    if (!createEnv.dag || createEnv.dag.length == 0) {
      return env;
    }

    // create all new components
    try {
      await this.compSvc.batchCreate(org, env, createEnv.dag.map(comp => comp.name))
    } catch (err) {
      if (err.code === SqlErrorCodes.DUP_ENTRY) {
        throw new BadRequestException('component already exists');
      }

      if (err.code === SqlErrorCodes.NO_DEFAULT) {
        throw new BadRequestException(err.sqlMessage);
      }

      this.logger.error({ message: 'could not batch create components during environment creation', err});
      throw new InternalServerErrorException('could not create components');
    }

    return env;
  }
}

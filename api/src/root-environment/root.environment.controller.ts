import { Controller, Get, Post, Body, Patch, Param, Delete, Request, BadRequestException, Logger, InternalServerErrorException } from '@nestjs/common';
import { ComponentService } from 'src/costing/services/component.service';
import { CreateEnvironmentDto } from 'src/environment/dto/create-environment.dto';
import { EnvSpecDto } from 'src/environment/dto/env-spec.dto';
import { SqlErrorCodes } from 'src/types';
import { RootEnvironmentService } from './root.environment.service';

@Controller({
  version: '1'
})
export class RootEnvironmentController {
  private readonly logger = new Logger(RootEnvironmentController.name);

  constructor(
    private readonly rootEnvSvc: RootEnvironmentService,
    private readonly compSvc: ComponentService
  ) {}

  @Post('/spec')
  async spec(@Request() req, @Body() body: EnvSpecDto) {
    const { org, team } = req;
    let env = req.env;

    if (!env) {
      env = await this.rootEnvSvc.create({
        name: body.envName,
        team,
        org
      });

      if (!env) {
        this.logger.error({ message: 'could not create environment', env, specDto: body})
        throw new InternalServerErrorException('could not create environment');
      }
    }

    // get components or create
    const dbComps = await this.compSvc.getAllForEnvironment(org, env);

    const existingComponents = body.components.filter(incoming => {
      return dbComps.find(dbComp => incoming.name === dbComp.name)
    });

    const newComponents = body.components.filter(incoming => {
      return !dbComps.find(dbComp => incoming.name === dbComp.name)
    });
  }

  @Post()
  async create(@Request() req, @Body() createEnv: CreateEnvironmentDto) {
    const {org, team} = req;

    createEnv.org = org;
    createEnv.team = team;

    try {
      return await this.rootEnvSvc.create(createEnv);
    } catch (ex) {
      if (ex.code === SqlErrorCodes.DUP_ENTRY) {
        throw new BadRequestException('environment already exists');
      }
      
      this.logger.error({ message: 'could not create environment', createEnv, ex });
      throw new InternalServerErrorException('could not create environment');
    }
  }

  @Get()
  async findAll(@Request() req) {
    const {org, team} = req;

    return this.rootEnvSvc.findAll(org, team);
  }
}

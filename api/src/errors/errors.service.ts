import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { EnvironmentService } from 'src/environment/environment.service';
import { ReconciliationService } from 'src/reconciliation/reconciliation.service';
import { Environment, Organization, Team } from 'src/typeorm';

@Injectable()
export class ErrorsService {
  constructor(
    private readonly reconSvc: ReconciliationService,
    private readonly envSvc: EnvironmentService
  ) {}

  async processValidRecon(org: Organization, team: Team, env: Environment) {
    if (env.latestEnvRecon?.status !== 'validation_failed') {
      // all well no need to create a new entry
      return;
    }

    // This means the yaml is fixed
    // check if there is a prev non validation error status
    const recon = await this.reconSvc.getLatestNonValidationErrorRec(org, env);
    if (!recon) {
      // No need to create a new entry since it is a new recon, will be handled by knitter
      return;
    }

    const newRecon = await this.getNewEnvRecon(
      org,
      team,
      env,
      null,
      recon.status,
      recon.estimatedCost
    );
    // Create a new recon with this recon values

    newRecon.environment = null;

    return this.envSvc.mergeAndSaveEnv(org, env, {
      latestEnvRecon: newRecon,
    });
  }

  async processInvalidRecon(
    org: Organization,
    team: Team,
    env: Environment,
    errorMessage: string[]
  ) {
    // This means there is a validation_error in yaml

    // check if its the same and return
    if (
      JSON.stringify(errorMessage) ==
      JSON.stringify(env.latestEnvRecon?.errorMessage)
    ) {
      return;
    }

    // This is a new error so create a new entry

    const envRecon = await this.getNewEnvRecon(
      org,
      team,
      env,
      errorMessage,
      'validation_failed',
      0
    );
    envRecon.environment = null;

    return this.envSvc.mergeAndSaveEnv(org, env, {
      latestEnvRecon: envRecon,
    });
  }

  async getNewEnvRecon(
    org: Organization,
    team: Team,
    env: Environment,
    errorMessage: string[] | null = null,
    status: string = 'initializing',
    estimatedCost: number = 0
  ) {
    return this.reconSvc.createErrorEnvRecon(org, team, env, {
      components: env.dag,
      startDateTime: new Date().toISOString(),
      endDateTime: new Date().toISOString(),
      name: env.name,
      status,
      teamName: team.name,
      errorMessage,
      estimatedCost,
    });
  }
}

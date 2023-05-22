import { Controller, Post, Request } from '@nestjs/common';
import { ApiTags } from '@nestjs/swagger';
import { EnvironmentService } from 'src/environment/environment.service';
import { APIRequest, EnvironmentApiParam } from 'src/types';
import { GithubApiService } from './github-api.service';
import { get } from 'src/config';

@Controller({
  version: '1',
})
@ApiTags('github-api')
export class GithubApiController {
  constructor(
    private readonly envSvc: EnvironmentService,
    private readonly gitSvc: GithubApiService
  ) {}

  @Post('/:environmentId')
  @EnvironmentApiParam()
  async gitCommit(@Request() req: APIRequest) {
    const { org, team, env } = req;
    const environment = await this.envSvc.findById(org, env.id);
    if (environment) {
      return this.gitSvc.gitCommit(org, get().github.owner, get().github.repo, `environments/play/env.yaml`);
    }
  }
}

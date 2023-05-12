import { Routes } from '@nestjs/core';
import { AuthModule } from './auth/auth.module';
import { OrganizationModule } from './organization/organization.module';
import { ReconciliationModule } from './reconciliation/reconciliation.module';
import { SecretsModule } from './secrets/secrets.module';
import { UsersModule } from './users/users.module';
import { SystemModule } from './system/system.module';
import { OperationsModule } from './operations/operations.module';
import { TeamModule } from './team/team.module';
import { EnvironmentModule } from './environment/environment.module';
import { ComponentModule } from './component/component.module';
import { StreamModule } from './stream/stream.module';
import { ErrorsModule } from './errors/errors.module';
import { GithubApiModule } from './github-api/github-api.module';

export const appRoutes: Routes = [
  {
    path: '/users',
    module: UsersModule,
  },
  {
    path: '/system',
    module: SystemModule,
  },
  {
    path: '/orgs',
    module: OrganizationModule,
    children: [
      {
        path: '/:orgId/secrets',
        module: SecretsModule,
      },
      {
        path: '/:orgId/stream',
        module: StreamModule,
      },
      {
        path: '/:orgId/auth',
        module: AuthModule,
      },
      {
        path: '/:orgId/ops',
        module: OperationsModule,
      },
      {
        path: '/:orgId/reconciliation',
        module: ReconciliationModule,
      },
      {
        path: '/:orgId/teams',
        module: TeamModule,
        children: [
          {
            path: '/:teamId/environments',
            module: EnvironmentModule,
            children: [
              {
                path: '/:environmentId/components',
                module: ComponentModule,
              },
            ],
          },
          {
            path: '/:teamId/github',
            module: GithubApiModule,
          },
          {
            path: '/:teamId/errors',
            module: ErrorsModule,
          },
        ],
      },
    ],
  },
];

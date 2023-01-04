import { Routes } from "@nestjs/core";
import { AuthModule } from "./auth/auth.module";
import { CostingModule } from "./costing/costing.module";
import { OrganizationModule } from "./organization/organization.module";
import { RootOrganizationsModule } from "./root-organization/root.organization.module";
import { ReconciliationModule } from "./reconciliation/reconciliation.module";
import { SecretsModule } from "./secrets/secrets.module";
import { UsersModule } from "./users/users.module";
import { SystemModule } from "./system/system.module";
import { OperationsModule } from "./operations/operations.module";
import { TeamModule } from "./team/team.module";
import { RootTeamModule } from "./root-team/root.team.module";
import { EnvironmentModule } from "./environment/environment.module";
import { RootEnvironmentModule } from "./root-environment/root.environment.module";
import { ComponentModule } from "./component/component.module";

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
    module: RootOrganizationsModule,
    children: [
      {
        path: '/:orgId',
        module: OrganizationModule,
        children: [
          {
            path: "secrets",
            module: SecretsModule
          },
          {
            path: "auth",
            module: AuthModule
          },
          {
            path: "ops",
            module: OperationsModule
          },
          {
            path: 'reconciliation',
            module: ReconciliationModule
          },
          {
            path: 'teams',
            module: RootTeamModule,
            children: [
              {
                path: ':teamId',
                module: TeamModule,
                children: [
                  {
                    path: 'environments',
                    module: RootEnvironmentModule,
                    children: [
                      {
                        path: ':environmentId',
                        module: EnvironmentModule,
                        children: [
                          {
                            path: 'components',
                            module: ComponentModule
                          }
                        ]
                      },
                      {
                        path: 'costing',
                        module: CostingModule
                      }
                    ]
                  },
                ]
              }
            ]
          }
        ]
      }
    ]
  }
]

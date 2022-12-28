import { Routes } from "@nestjs/core";
import { AuthModule } from "./auth/auth.module";
import { CostingModule } from "./costing/costing.module";
import { OrganizationModule } from "./organization/organization.module";
import { RootOrganizationsModule } from "./rootOrganization/rootOrganization.module";
import { ReconciliationModule } from "./reconciliation/reconciliation.module";
import { SecretsModule } from "./secrets/secrets.module";
import { UsersModule } from "./users/users.module";
import { SystemModule } from "./system/system.module";
import { OperationsModule } from "./operations/operations.module";
import { TeamModule } from "./team/team.module";
import { RootTeamModule } from "./root-team/root-team.module";

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
            path: 'teams',
            module: RootTeamModule,
            children: [
              {
                path: '/:teamId',
                module: TeamModule,
                children: [
                  {
                    path: 'costing',
                    module: CostingModule
                  },
                  {
                    path: 'reconciliation',
                    module: ReconciliationModule
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

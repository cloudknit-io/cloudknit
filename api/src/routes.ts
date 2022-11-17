import { Routes } from "@nestjs/core";
import { AuthModule } from "./auth/auth.module";
import { CostingModule } from "./costing/costing.module";
import { OrganizationModule } from "./organization/organization.module";
import { RootOrganizationsModule } from "./rootOrganization/rootOrganization.module";
import { ReconciliationModule } from "./reconciliation/reconciliation.module";
import { SecretsModule } from "./secrets/secrets.module";
import { UsersModule } from "./users/users.module";
import { OperationsModule } from "./operations/operations.module";

export const orgRoutes: Routes = [
  {
    path: '/users',
    module: UsersModule,
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
            path: 'costing',
            module: CostingModule
          },
          {
            path: 'reconciliation',
            module: ReconciliationModule
          },
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
        ]
      }
    ]
  }
]

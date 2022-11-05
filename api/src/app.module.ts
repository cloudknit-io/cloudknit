import { Module } from "@nestjs/common";
import { ConfigModule } from "@nestjs/config";
import { RouterModule } from "@nestjs/core";
import { TypeOrmModule, TypeOrmModuleOptions } from "@nestjs/typeorm";
import { AuthModule } from "./auth/auth.module";
import { CostingModule } from "./costing/costing.module";
import { OrganizationModule } from "./organization/organization.module";
import { RootOrganizationsModule } from "./rootOrganization/rootOrganization.module";
import { ReconciliationModule } from "./reconciliation/reconciliation.module";
import { orgRoutes } from "./routes";
import { SecretsModule } from "./secrets/secrets.module";
import { entities } from "./typeorm";
import { UsersModule } from "./users/users.module";
import { get } from "./config";

const config = get();

const typeOrmModuleOptions: TypeOrmModuleOptions = {
  type: "mysql",
  host: config.TypeORM.host,
  port: config.TypeORM.port,
  username: config.TypeORM.username,
  password: config.TypeORM.password,
  database: config.TypeORM.database,
  entities: entities,
  migrations: [],
  synchronize: true,
};

@Module({
  imports: [
    ConfigModule.forRoot({
      envFilePath: ".env.dev",
    }),
    RouterModule.register(orgRoutes),
    TypeOrmModule.forRoot(typeOrmModuleOptions),
    UsersModule,
    RootOrganizationsModule,
    OrganizationModule,
    CostingModule,
    ReconciliationModule,
    SecretsModule,
    AuthModule,
  ],
  controllers: [],
  providers: [],
})
export class AppModule {}

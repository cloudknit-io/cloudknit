import { MiddlewareConsumer, Module, NestModule } from "@nestjs/common";
import { ConfigModule } from "@nestjs/config";
import { RouterModule } from "@nestjs/core";
import { TypeOrmModule, TypeOrmModuleOptions } from "@nestjs/typeorm";
import { AuthModule } from "./auth/auth.module";
import { CostingModule } from "./costing/costing.module";
import { OrganizationModule } from "./organization/organization.module";
import { RootOrganizationsModule } from "./rootOrganization/rootOrganization.module";
import { ReconciliationModule } from "./reconciliation/reconciliation.module";
import { appRoutes } from "./routes";
import { SecretsModule } from "./secrets/secrets.module";
import { entities } from "./typeorm";
import { UsersModule } from "./users/users.module";
import { SystemModule } from "./system/system.module";
import { get } from "./config";
import { OperationsModule } from './operations/operations.module';
import { AppLoggerMiddleware } from "./middleware/logger.middle";
import { TeamModule } from './team/team.module';
import { RootTeamModule } from './root-team/root-team.module';
import { EnvironmentModule } from './environment/environment.module';

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
    RouterModule.register(appRoutes),
    TypeOrmModule.forRoot(typeOrmModuleOptions),
    UsersModule,
    SystemModule,
    RootOrganizationsModule,
    OrganizationModule,
    CostingModule,
    ReconciliationModule,
    SecretsModule,
    AuthModule,
    OperationsModule,
    TeamModule,
    RootTeamModule,
    EnvironmentModule
  ],
  controllers: [],
  providers: [],
})
export class AppModule implements NestModule {
  configure(consumer: MiddlewareConsumer) {
    consumer.apply(AppLoggerMiddleware).forRoutes('*');
  }
}

import { Module } from "@nestjs/common";
import { ConfigModule } from "@nestjs/config";
import { TypeOrmModule, TypeOrmModuleOptions } from "@nestjs/typeorm";
import { AppController } from "./app.controller";
import { AppService } from "./app.service";
import { AuthModule } from "./auth/auth.module";
import { CostingModule } from "./costing/costing.module";
import { CompanyModule } from "./company/company.module";
import { ReconciliationModule } from "./reconciliation/reconciliation.module";
import { SecretsModule } from "./secrets/secrets.module";
import { entities } from "./typeorm";

const typeOrmModuleOptions: TypeOrmModuleOptions = {
  type: "mysql",
  host: process.env.TYPEORM_HOST,
  port: parseInt(process.env.TYPEORM_PORT),
  username: process.env.TYPEORM_USERNAME,
  password: process.env.TYPEORM_PASSWORD,
  database: process.env.TYPEORM_DATABASE,
  entities: entities,
  synchronize: true,
};

@Module({
  imports: [
    ConfigModule.forRoot({
      envFilePath: ".env.dev",
    }),
    AuthModule,
    TypeOrmModule.forRoot(typeOrmModuleOptions),
    CostingModule,
    ReconciliationModule,
    SecretsModule,
    AuthModule,
    CompanyModule,
  ],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule {}

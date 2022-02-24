import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Connection } from 'typeorm';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { AuthModule } from './auth/auth.module';
import { CostingModule } from './costing/costing.module';
import { ReconciliationModule } from './reconciliation/reconciliation.module';
import { SecretsModule } from './secrets/secrets.module';
import { entities } from './typeorm';


@Module({
  imports: [
    ConfigModule.forRoot({
      envFilePath: '.env.dev',
    }),
    AuthModule,
    TypeOrmModule.forRoot({
        type: "mysql",
        host: process.env.TYPEORM_HOST || "mysqldb",
        port: 3306,
        username: process.env.TYPEORM_USERNAME || "root",
        password: process.env.TYPEORM_PASSWORD || "password",
        database: process.env.TYPEORM_DATABASE || "nestjsrealworld",
        entities: entities,
        synchronize: true
    }),
    CostingModule,
    ReconciliationModule,
    SecretsModule,
    AuthModule
  ],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule {}

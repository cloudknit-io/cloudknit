import { MiddlewareConsumer, Module, NestModule } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { RouterModule } from '@nestjs/core';
import { TypeOrmModule, TypeOrmModuleOptions } from '@nestjs/typeorm';
import { AuthModule } from './auth/auth.module';
import { OrganizationModule } from './organization/organization.module';
import { ReconciliationModule } from './reconciliation/reconciliation.module';
import { appRoutes } from './routes';
import { SecretsModule } from './secrets/secrets.module';
import { UsersModule } from './users/users.module';
import { SystemModule } from './system/system.module';
import { OperationsModule } from './operations/operations.module';
import { AppLoggerMiddleware } from './middleware/logger.middle';
import { TeamModule } from './team/team.module';
import { EnvironmentModule } from './environment/environment.module';
import { ComponentModule } from './component/component.module';
import { StreamModule } from './stream/stream.module';
import { CachingService } from './caching/caching.service';
import { CachingModule } from './caching/caching.module';
import { dbConfig } from './typeorm';

@Module({
  imports: [
    ConfigModule.forRoot({
      envFilePath: '.env.dev',
    }),
    RouterModule.register(appRoutes),
    TypeOrmModule.forRoot(dbConfig as TypeOrmModuleOptions),
    UsersModule,
    SystemModule,
    OrganizationModule,
    ReconciliationModule,
    SecretsModule,
    AuthModule,
    OperationsModule,
    TeamModule,
    EnvironmentModule,
    ComponentModule,
    StreamModule,
    CachingModule,
  ],
  controllers: [],
  providers: [CachingService],
})
export class AppModule implements NestModule {
  configure(consumer: MiddlewareConsumer) {
    consumer.apply(AppLoggerMiddleware).forRoutes('*');
  }
}

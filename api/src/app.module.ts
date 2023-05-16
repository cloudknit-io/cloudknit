import { Logger, MiddlewareConsumer, Module, NestModule } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { RouterModule } from '@nestjs/core';
import { EventEmitterModule } from '@nestjs/event-emitter';
import { TypeOrmModule, TypeOrmModuleOptions } from '@nestjs/typeorm';
import { Connection } from 'typeorm';
import { AuthModule } from './auth/auth.module';
import { CachingModule } from './caching/caching.module';
import { CachingService } from './caching/caching.service';
import { ComponentModule } from './component/component.module';
import { EnvironmentModule } from './environment/environment.module';
import { ErrorsModule } from './errors/errors.module';
import { AppLoggerMiddleware } from './middleware/logger.middle';
import { OperationsModule } from './operations/operations.module';
import { OrganizationModule } from './organization/organization.module';
import { ReconciliationModule } from './reconciliation/reconciliation.module';
import { appRoutes } from './routes';
import { SecretsModule } from './secrets/secrets.module';
import { StreamModule } from './stream/stream.module';
import { SystemModule } from './system/system.module';
import { TeamModule } from './team/team.module';
import { dbConfig } from './typeorm';
import { UsersModule } from './users/users.module';

@Module({
  imports: [
    ConfigModule.forRoot({
      envFilePath: '.env.dev',
    }),
    EventEmitterModule.forRoot({
      verboseMemoryLeak: true
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
    ErrorsModule,
  ],
  controllers: [],
  providers: [CachingService],
})
export class AppModule implements NestModule {
  constructor(connection: Connection) {
    this.synchronize(connection);
  }

  configure(consumer: MiddlewareConsumer) {
    consumer.apply(AppLoggerMiddleware).forRoutes('*');
  }

  synchronize(connection: Connection) {
    if (connection.entityMetadatas.length > 3) {
      return;
    }
    const logger = new Logger('Synchronize');
    logger.error(`Synchronize Error`, '0 tables found, synchronizing...');
    connection.synchronize(false);
  }
}

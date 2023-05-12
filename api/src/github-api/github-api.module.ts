import { MiddlewareConsumer, Module, RequestMethod } from '@nestjs/common';
import { GithubApiController } from './github-api.controller';
import { EnvironmentService } from 'src/environment/environment.service';
import { EnvironmentMiddleware } from 'src/middleware/environment.middle';
import { GithubApiService } from './github-api.service';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Environment } from 'src/typeorm/environment.entity';

@Module({
  imports: [TypeOrmModule.forFeature([Environment])],
  controllers: [GithubApiController],
  providers: [EnvironmentService, GithubApiService],
})
export class GithubApiModule {
  configure(consumer: MiddlewareConsumer) {
    consumer.apply(EnvironmentMiddleware).forRoutes({
      path: '*/github/:environmentId*',
      method: RequestMethod.ALL,
    });
  }
}

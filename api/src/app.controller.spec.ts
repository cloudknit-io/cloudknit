import { Test } from '@nestjs/testing';
import { AppService } from './app.service';
import * as request from 'supertest';
import { AppModule } from './app.module';
import { INestApplication } from '@nestjs/common';

describe('AppController', () => {
  let app: INestApplication;
  const appService: AppService = { getHello: () => 'Hello Everyone!' };

  beforeEach(async () => {
    const moduleRef = await Test.createTestingModule({
      imports: [AppModule],
    })
      .overrideProvider(AppService)
      .useValue(appService)
      .compile();

    app = moduleRef.createNestApplication();
    await app.init();
  });

  describe('root', () => {
    it('should return "Hello Everyone!"', () => {
      return request(app.getHttpServer())
        .get('/')
        .expect(200)
        .expect('Hello Everyone!');
    });
  });
});

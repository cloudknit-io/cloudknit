import { VersioningType } from '@nestjs/common';
import { NestFactory } from '@nestjs/core';
import { DocumentBuilder, SwaggerModule } from '@nestjs/swagger';
import { AppModule } from './app.module';
import { get, init } from './config';
import { WinstonLogger } from './logger';

async function bootstrap() {
  init(); // init api config

  const config = get();
  const app = await NestFactory.create(AppModule, {
    logger: new WinstonLogger()
  });

  app.enableVersioning({
    type: VersioningType.URI
  });

  const openApiDoc = new DocumentBuilder()
    .setTitle('CloudKnit API')
    .setDescription('Internal API to manage organizations, users, secrets, and interactions with Argo Cd and Argo WF')
    .setVersion('0.1.0')
    .build();
  
  const document = SwaggerModule.createDocument(app, openApiDoc);
  SwaggerModule.setup('api', app, document);

  await app.listen(config.port);
}

bootstrap();

import { Module } from '@nestjs/common';
import { SecretsController } from './secrets.controller';
import { SecretsService } from './secrets.service';

@Module({
  controllers: [SecretsController],
  providers: [SecretsService],
})
export class SecretsModule {}

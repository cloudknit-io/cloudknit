import { Module } from '@nestjs/common';
import { StreamService } from './stream.service';
import { StreamController } from './stream.controller';
import { StreamEnvironmentService } from './stream.environment.service';
import { StreamComponentService } from './stream.component.service';

@Module({
  controllers: [StreamController],
  providers: [
    StreamService,
    StreamEnvironmentService,
    StreamComponentService
  ]
})
export class StreamModule {}

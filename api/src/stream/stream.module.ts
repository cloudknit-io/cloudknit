import { Module } from '@nestjs/common';
import { StreamService } from './stream.service';
import { StreamController } from './stream.controller';
import { StreamEnvironmentService } from './stream.environment.service';
import { StreamComponentService } from './stream.component.service';
import { StreamEnvironmentReconcileService } from './stream.env-reconcile.service';
import { StreamComponentReconcileService } from './stream.comp-reconcile.service';

@Module({
  controllers: [StreamController],
  providers: [
    StreamService,
    StreamEnvironmentService,
    StreamComponentService,
    StreamEnvironmentReconcileService,
    StreamComponentReconcileService,
  ],
})
export class StreamModule {}

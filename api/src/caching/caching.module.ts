import { Module } from '@nestjs/common';
import { CachingService } from './caching.service';

@Module({
  imports: [
    // ClientsModule.register([
    //   {
    //     name: 'MATH_SERVICE',
    //     transport: Transport.REDIS,
    //     options: {
    //       host: 'localhost',
    //       port: 6379,
    //     }
    //   },
    // ]),
  ],
  controllers: [],
  providers: [
    CachingService,
  ]
})
export class CachingModule {}

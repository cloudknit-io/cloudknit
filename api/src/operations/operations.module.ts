import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Organization } from 'src/typeorm';
import { OperationsController } from './operations.controller';
import { OperationsService } from './operations.service';

@Module({
  imports: [TypeOrmModule.forFeature([Organization])],
  controllers: [OperationsController],
  providers: [OperationsService]
})
export class OperationsModule {}

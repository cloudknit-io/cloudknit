import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { OrganizationService } from 'src/organization/organization.service';
import { Organization } from 'src/typeorm';
import { OperationsController } from './operations.controller';
import { OperationsService } from './operations.service';

@Module({
  imports: [TypeOrmModule.forFeature([Organization])],
  controllers: [OperationsController],
  providers: [OperationsService, OrganizationService]
})
export class OperationsModule {}

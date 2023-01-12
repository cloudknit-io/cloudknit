import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { OrganizationService } from 'src/organization/organization.service';
import { Organization, User } from 'src/typeorm';
import { UsersService } from 'src/users/users.service';
import { OperationsController } from './operations.controller';
import { OperationsService } from './operations.service';

@Module({
  imports: [TypeOrmModule.forFeature([Organization, User])],
  controllers: [OperationsController],
  providers: [OperationsService, OrganizationService, UsersService]
})
export class OperationsModule {}

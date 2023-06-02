import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { User } from 'src/typeorm/User.entity';
import { UsersController } from './users.controller';
import { UsersService } from './users.service';
import { Organization } from 'src/typeorm';

@Module({
  imports: [TypeOrmModule.forFeature([User, Organization])],
  controllers: [UsersController],
  providers: [UsersService],
  exports: [UsersService],
})
export class UsersModule {}

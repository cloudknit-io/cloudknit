import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { AuthController } from './auth.controller';
import { AuthService } from './auth.service';
import { User } from 'src/typeorm';
import { Organization } from 'src/typeorm/Organization.entity';

@Module({
  imports: [
    TypeOrmModule.forFeature([
      User,
      Organization
    ])
  ],
  controllers: [AuthController],
  providers: [AuthService],
})
export class AuthModule {}

import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { User } from 'src/typeorm/entities/User';
import { AuthController } from './controller/auth/auth.controller';
import { AuthService } from './services/auth/auth.service';

@Module({
  imports: [
    TypeOrmModule.forFeature([
      User
    ])
  ],
  controllers: [AuthController],
  providers: [AuthService],
})
export class AuthModule {}

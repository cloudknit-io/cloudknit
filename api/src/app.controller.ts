import { Body, Controller, Get, Param, Post, Req } from '@nestjs/common';
import { AppService } from './app.service';
import { Request } from 'express';
import { User } from './auth/typeorm/entities/User';

@Controller()
export class AppController {
  constructor(private readonly appService: AppService) {}

  @Get()
  getHello(): string {
    return this.appService.getHello();
  }

  @Post('add-user')
  async addUser(@Body() user: User): Promise<User> {
    return await this.appService.addUser(user);
  }

  @Get('get-user/:id')
  async getUser(@Param('id') id: number): Promise<User> {
    return await this.appService.getUser(id);
  }
}

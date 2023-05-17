import {
  Body,
  Controller,
  Get,
  Logger,
  NotFoundException,
  Param,
  Post,
  Query,
} from '@nestjs/common';
import { ApiTags } from '@nestjs/swagger';
import { AuthController } from 'src/auth/auth.controller';
import { User } from 'src/typeorm/User.entity';
import { CreatePlaygroundUserDto, CreateUserDto } from './User.dto';
import { UsersService } from './users.service';

@Controller({
  version: '1',
})
@ApiTags('users')
export class UsersController {
  private readonly logger = new Logger(AuthController.name);

  constructor(private readonly userService: UsersService) {}

  @Get('/:username')
  public async getUser(@Param('username') username: string): Promise<User> {
    const user = await this.userService.getUser(username);

    if (!user) {
      throw new NotFoundException();
    }

    return user;
  }

  @Post()
  public async createUser(@Body() body: CreateUserDto): Promise<User> {
    const user = await this.userService.create(body);

    return user;
  }

  @Get('/playground/:ipv4')
  public async getPlaygroundUser(@Param('ipv4') ipv4: string): Promise<User> {
    const user = await this.userService.getPlaygroundUser(ipv4);
    if (!user) {
      throw new NotFoundException();
    }

    return user;
  }

  @Post('/playground')
  public async createPlaygroundUser(@Body() user: CreatePlaygroundUserDto): Promise<User> {
    return this.userService.createPlaygroundUser(user);
  }
}

import {
  Body,
  Controller,
  Get,
  NotFoundException,
  Param,
  Post
} from '@nestjs/common';
import { ApiTags } from '@nestjs/swagger';
import { User } from 'src/typeorm/User.entity';
import { CreateUserDto } from './User.dto';
import { UsersService } from './users.service';

@Controller({
  version: '1',
})
@ApiTags('users')
export class UsersController {

  constructor(private readonly userService: UsersService) {}

  @Get('/:username')
  public async getUser(@Param('username') username: string): Promise<User> {
    console.log(username);
    const user = await this.userService.getUser(username);

    if (!user) {
      throw new NotFoundException();
    }

    return user;
  }

  @Post()
  public async createUser(@Body() body: CreateUserDto): Promise<User> {
    console.log(body);
    const user = await this.userService.create(body);

    return user;
  }
}

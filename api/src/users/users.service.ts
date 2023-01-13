import { BadRequestException, Injectable, Logger } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { User } from 'src/typeorm';
import { Repository } from 'typeorm';
import { CreateUserDto } from './User.dto';

@Injectable()
export class UsersService {
  private readonly logger = new Logger(UsersService.name);

  constructor(@InjectRepository(User) private userRepo: Repository<User>) {}

  async getUser(username: string): Promise<User> {
    return this.userRepo.findOne({
      where: { username },
      relations: {
        organizations: true,
      },
    });
  }

  async getUserById(userId: number): Promise<User> {
    return this.userRepo.findOne({
      where: { id: userId },
      relations: {
        organizations: true,
      },
    });
  }

  async create(userDto: CreateUserDto) {
    // User should not exist
    const userExists = await this.getUser(userDto.username);

    if (userExists) {
      throw new BadRequestException({
        message: 'User with Github Id already exists',
      });
    }

    const newUser = new User();
    newUser.email = userDto.email;
    newUser.name = userDto.name;
    newUser.username = userDto.username;
    newUser.role = userDto.role;
    newUser.organizations = [];

    const user = await this.userRepo.save(newUser);

    this.logger.log('created user', { user: userDto });

    return user;
  }
}

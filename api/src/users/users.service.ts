import {
  BadRequestException,
  Injectable,
  Logger,
  NotFoundException,
} from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Organization, User } from 'src/typeorm';
import { UserRole } from 'src/types';
import { Repository } from 'typeorm';
import { CreateUserDto } from './User.dto';

@Injectable()
export class UsersService {
  private readonly logger = new Logger(UsersService.name);

  constructor(
    @InjectRepository(User) private userRepo: Repository<User>,
    @InjectRepository(Organization) private orgRepo: Repository<Organization>
  ) {}

  async getUser(username: string): Promise<User> {
    const user = await this.userRepo.findOne({
      where: { username },
      relations: {
        organizations: true,
      },
    });

    if (!user) {
      return null;
    }

    if (user.role === UserRole.GUEST) {
      return this.associateOrganization(user);
    }
    return user;
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

    console.log(user);

    if (user.role === UserRole.GUEST) {
      return this.associateOrganization(user);
    }

    this.logger.log('created user', { user: userDto });

    return user;
  }

  private async associateOrganization(user: User) {
    if (user.organizations.length === 0) {
      const org = await this.getOrganizationWithoutUserAssociation();
      const updatedUser = this.userRepo.merge(user, {
        organizations: [org]
      });
      user = await this.userRepo.save(updatedUser);
    }
    return user;
  }

  private async getOrganizationWithoutUserAssociation() {
    const orgs = await this.orgRepo.find({
      relations: {
        users: true,
      },
    });

    const org = orgs.find((org) => !org.users.some(user => user.role === UserRole.GUEST));

    if (!org) {
      throw new NotFoundException('No Organization is present at the moment.');
    }

    return org;
  }
}

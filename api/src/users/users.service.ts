import {
  BadRequestException,
  Injectable,
  Logger,
  NotFoundException,
} from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { randomUUID } from 'crypto';
import { Organization, User } from 'src/typeorm';
import { Repository } from 'typeorm';
import { CreatePlaygroundUserDto, CreateUserDto } from './User.dto';

@Injectable()
export class UsersService {
  private readonly logger = new Logger(UsersService.name);

  constructor(
    @InjectRepository(User) private userRepo: Repository<User>,
    @InjectRepository(Organization) private orgRepo: Repository<Organization>
  ) {}

  async getUser(username: string): Promise<User> {
    return this.userRepo.findOne({
      where: { username },
      relations: {
        organizations: true,
      },
    });
  }

  async getPlaygroundUser(ipv4: string): Promise<User> {
    let user = await this.userRepo.findOne({
      where: { ipv4 },
      relations: {
        organizations: true,
      },
    });

    if (user.organizations.length === 0) {
      const org = await this.getOrganizationWithoutUserAssociation();
      const updatedUser = this.userRepo.merge(user, {
        organizations: [org]
      });
      user = await this.userRepo.save(updatedUser);
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

    this.logger.log('created user', { user: userDto });

    return user;
  }

  public async createPlaygroundUser(user: CreatePlaygroundUserDto) {
    if (!user.ipv4) {
      throw new BadRequestException('request does not have a valid ip address');
    }

    const currentUser = await this.getPlaygroundUser(user.ipv4);

    if (currentUser) {
      throw new BadRequestException('User already exists');
    }

    // Get an org that is not associated to any user

    const org = await this.getOrganizationWithoutUserAssociation();
    

    const uuid = `guest-${randomUUID()}`;
    // Create user
    const newUser = new User();
    newUser.email = `${uuid}@cloudknit.io`;
    newUser.name = uuid;
    newUser.username = uuid;
    newUser.role = 'Guest';
    newUser.ipv4 = user.ipv4;
    newUser.organizations = [org];

    return this.userRepo.save(newUser);
  }

  private async getOrganizationWithoutUserAssociation() {
    const orgs = await this.orgRepo.find({
      relations: {
        users: true,
      },
    });

    const org = orgs.find((org) => org.users.length === 0);

    if (!org) {
      throw new NotFoundException('No Organization is present at the moment.');
    }

    return org;
  }
}

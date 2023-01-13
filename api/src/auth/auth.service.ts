import { BadRequestException, Injectable, Logger } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Organization, User } from 'src/typeorm';
import { CreateUserDto } from 'src/users/User.dto';
import { Repository } from 'typeorm';

@Injectable()
export class AuthService {
  private readonly logger = new Logger(AuthService.name);

  constructor(@InjectRepository(User) private readonly userRepo: Repository<User>) {}

  // public async getUser(username: string) {
  //   return this.userRepo.createQueryBuilder('user')
  //     .leftJoinAndSelect('user.organizations', 'organization')
  //     .where('organization.id = :orgId and user.username = :username', { orgId: org.id, username })
  //     .getOne();
  // }

  public async getOrgUser(org: Organization, username: string) {
    return this.userRepo
      .createQueryBuilder('user')
      .leftJoinAndSelect('user.organizations', 'organization')
      .where('organization.id = :orgId and user.username = :username', {
        orgId: org.id,
        username,
      })
      .getOne();
  }

  public async getOrgUserList(org: Organization) {
    return this.userRepo
      .createQueryBuilder('user')
      .leftJoinAndSelect('user.organizations', 'organization')
      .where('organization.id = :orgId', { orgId: org.id })
      .getMany();
  }

  public async createOrgUser(org: Organization, user: CreateUserDto) {
    // This is to get a user irrespective of the org he belongs to
    // Since, if an admin wants to add an existing user who is not a part of the currect org
    // will always give empty result and creating a new user would throw an error.
    const currentUser = await this.userRepo.findOne({
      where: {
        username: user.username,
      },
      relations: {
        organizations: true,
      },
    });

    if (currentUser) {
      // adds existing user to org
      for (let userOrg of currentUser.organizations) {
        if (userOrg.id == org.id) {
          throw new BadRequestException(`${currentUser.username} is already a member of ${org.name}`);
        }
      }

      currentUser.organizations = [...currentUser.organizations, org];

      this.logger.log(`adding user ${currentUser.username} to ${org.name}`);

      return this.userRepo.save(currentUser);
    }

    // Create user
    const newUser = new User();
    newUser.email = user.email;
    newUser.name = user.name;
    newUser.username = user.username;
    newUser.role = user.role;
    newUser.organizations = [org];

    return this.userRepo.save(newUser);
  }

  // Delete user only removes the User <-> Org associattion
  public async deleteUser(org: Organization, username: string) {
    const user = await this.userRepo.findOne({
      where: {
        username: username,
      },
      relations: {
        organizations: true,
      },
    });

    if (!user) {
      throw new BadRequestException('user does not exist');
    }

    const userOrgs = user.organizations;
    let newOrgs = [];

    for (let userOrg of userOrgs) {
      if (userOrg.id === org.id) {
        continue;
      }

      newOrgs.push(userOrg);
    }

    user.organizations = newOrgs;

    return this.userRepo.save(user);
  }
}

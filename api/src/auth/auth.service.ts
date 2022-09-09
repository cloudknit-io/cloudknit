import { BadRequestException, Injectable, Logger } from "@nestjs/common";
import { InjectRepository } from "@nestjs/typeorm";
import { Organization, User} from "src/typeorm";
import { CreateUserDto } from "src/users/User.dto";
import { Repository } from "typeorm";

@Injectable()
export class AuthService {
  private readonly logger = new Logger(AuthService.name);

  constructor(
    @InjectRepository(User) private readonly userRepo: Repository<User>
  ) {}

  public async getTermAgreementStatus(org: Organization, body: any) {
    return this.getOrgUser(org, body).then(
      (user: User) => user?.termAgreementStatus || false
    );
  }

  public async setTermAgreementStatus(org: Organization, username) {
    const user = await this.getOrgUser(org, username);

    user.termAgreementStatus = true;
    return await this.userRepo.save(user);
  }

  public async getOrgUser(org: Organization, username: string) {
    return this.userRepo.createQueryBuilder('user')
      .leftJoinAndSelect('user.organizations', 'organization')
      .where('organization.id = :orgId and user.username = :username', { orgId: org.id, username })
      .getOne();
  }

  public async getOrgUserList(org: Organization) {
    return this.userRepo.createQueryBuilder('user')
      .leftJoinAndSelect('user.organizations', 'organization')
      .where('organization.id = :orgId', {orgId: org.id})
      .getMany();
  }

  public async createOrgUser(org: Organization, user: CreateUserDto) {
    const currentUser = await this.getOrgUser(org, user.username);
    
    if (currentUser) {
      // adds existing user to org
      for (let userOrg of currentUser.organizations) {
        if (userOrg.id == org.id) {
          throw new BadRequestException({
            message: `${currentUser.username} is already a member of ${org.name}`,
          });
        }
      }

      currentUser.organizations = [...currentUser.organizations, org];
      
      this.logger.log(`adding user ${currentUser.username} to ${org.name}`);
      
      return this.userRepo.save(currentUser);
    }

    // Create user
    const newUser = new User();
    newUser.email = user.email;
    newUser.username = user.username;
    newUser.role = user.role;
    newUser.organizations = [org];
    
    return this.userRepo.save(newUser);
  }

  public async deleteUser(org: Organization, username: string) {
    const user = await this.getOrgUser(org, username);

    if (!user) {
      throw new BadRequestException({
        message: "user does not exist",
      });
    }

    const userOrgs = user.organizations;
    let newOrgs = []

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

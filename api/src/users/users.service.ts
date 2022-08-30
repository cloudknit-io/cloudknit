import { Injectable, Logger } from "@nestjs/common";
import { InjectRepository } from "@nestjs/typeorm";
import { User } from 'src/typeorm';
import { Repository } from "typeorm";

@Injectable()
export class UsersService {
  private readonly logger = new Logger(UsersService.name);

  constructor(@InjectRepository(User) private userRepo: Repository<User>) {}

  async getUser(username: string) {
    return this.userRepo.findOne({ 
      where : { username },
      relations: {
        organizations: true
      }
    });
  }

  // MAKE THE FOLLOWING ENDPOINTS IN THE USERS MODULE
  // - GET /users
  // - POST /users
  // - GET /users/:userId

  // async create(user: CreateUserDto) {
  //   // User should not exist
  //   const userExists = await this.getUser({ username });
    
  //   if (userExists) {
  //     throw new BadRequestException({
  //       message: "User with Github Id already exists",
  //     });
  //   }

  //   // Org should exist
  //   const organization = await this.getOrg(organizationId);
    
  //   if (!organization) {
  //     throw new BadRequestException({
  //       message: "Organization does not exist",
  //     });
  //   }

  //   // Create user
  //   const user = await this.userRepo.save({
  //     username,
  //     email,
  //     role: role || "User",
  //   });
    
  //   return this.orgUserRepo.save({
  //     user,
  //     organization
  //   })

  //   return this.userRepo.save(user)
  // }
}

import { BadRequestException, Injectable } from "@nestjs/common";
import { InjectRepository } from "@nestjs/typeorm";
import { User } from "src/typeorm/entities/User";
import { Repository } from "typeorm";

@Injectable()
export class AuthService {
  constructor(
    @InjectRepository(User) private readonly userRepo: Repository<User>
  ) {}

  public async getUser({ username }) {
    return this.userRepo.findOne({
      where: {
        username: username,
      },
    });
  }

  public async getTermAgreementStatus(body: any) {
    return this.getUser(body).then(
      (user: User) => user?.termAgreementStatus || false
    );
  }

  public async setTermAgreementStatus({ username }) {
    const user = await this.userRepo.findOne({
      where: {
        username: username,
      },
    });
    user.termAgreementStatus = true;
    return await this.userRepo.save(user);
  }

  public async getUserList(organizationId: string) {
    return this.userRepo.find({
      where: {
        company: organizationId,
      },
    });
  }

  public async addUser({ username, organizationId, email, role }) {
    const existing = await this.getUser({ username });
    if (existing) {
      throw new BadRequestException({
        message: "User with Github Id already exists!",
      });
    }
    return this.userRepo.save({
      username,
      company: organizationId,
      email,
      role: role || "User",
    });
  }

  public async deleteUser(username: string) {
    const existing = await this.getUser({ username });
    if (!existing) {
      throw new BadRequestException({
        message: "User with Github Id does not exist!",
      });
    }
    return this.userRepo.delete({
      username: username,
    });
  }
}

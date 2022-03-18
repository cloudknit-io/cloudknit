import { Injectable } from "@nestjs/common";
import { InjectRepository } from "@nestjs/typeorm";
import { User } from "src/typeorm/entities/User";
import { Repository } from "typeorm";

@Injectable()
export class AuthService {
  constructor(
    @InjectRepository(User) private readonly userRepo: Repository<User>
  ) {}

  public async getTermAgreementStatus({username, company}) {
    const user = await this.userRepo.findOne({
      where: {
        username: username,
        company: company,
      },
    });
    return user?.termAgreementStatus || false;
  }

  public async setTermAgreementStatus({username, company, email, agreedByEmail, agreedByUsername}) {
    return await this.userRepo.save({
      company: company,
      termAgreementStatus: true,
      username: username,
      email: email,
      agreedByEmail: agreedByEmail,
      agreedByUsername: agreedByUsername
    });
  }
}

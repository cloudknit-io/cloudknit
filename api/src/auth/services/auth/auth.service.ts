import { Injectable } from "@nestjs/common";
import { InjectRepository } from "@nestjs/typeorm";
import { User } from "src/typeorm/entities/User";
import { Repository } from "typeorm";

@Injectable()
export class AuthService {
  constructor(
    @InjectRepository(User) private readonly userRepo: Repository<User>
  ) {}

  public async getTermAgreementStatus(company: string) {
    const user = await this.userRepo.findOne({
      where: {
        company: company,
      },
    });
    return user.termAgreementStatus;
  }
}

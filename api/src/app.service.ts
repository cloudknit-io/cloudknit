import { Injectable } from '@nestjs/common';
import { Connection } from 'typeorm';
import { User } from './typeorm/entities/User';

@Injectable()
export class AppService {
  constructor(private readonly connection: Connection) {}

  getHello(): string {
    return 'Hello everyone!';
  }

  async addUser(user: User): Promise<User> {
    const userRepo = await this.connection.getRepository(User);
    return await userRepo.save(user);
  }

  async getUser(id: number): Promise<User> {
    const repo = await this.connection.getRepository(User);
    return await repo.findOne(id);
  }
}

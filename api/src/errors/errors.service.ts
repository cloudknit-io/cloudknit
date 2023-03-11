import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Environment } from 'src/typeorm';
import { Repository } from 'typeorm';

@Injectable()
export class ErrorsService {
  constructor(
    @InjectRepository(Environment)
    private readonly envRepo: Repository<Environment>
  ) {}
}
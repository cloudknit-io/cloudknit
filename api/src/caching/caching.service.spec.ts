import { Test, TestingModule } from '@nestjs/testing';
import { CachingService } from './caching.service';

describe('CachingService', () => {
  let service: CachingService;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [CachingService],
    }).compile();

    service = module.get<CachingService>(CachingService);
  });

  it('should be defined', () => {
    expect(service).toBeDefined();
  });
});

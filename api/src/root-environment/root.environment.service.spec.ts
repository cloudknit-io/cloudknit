import { Test, TestingModule } from '@nestjs/testing';
import { RootEnvironmentService } from './root.environment.service';

describe('RootEnvironmentService', () => {
  let service: RootEnvironmentService;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [RootEnvironmentService],
    }).compile();

    service = module.get<RootEnvironmentService>(RootEnvironmentService);
  });

  it('should be defined', () => {
    expect(service).toBeDefined();
  });
});

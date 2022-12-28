import { Test, TestingModule } from '@nestjs/testing';
import { RootTeamService } from './root-team.service';

describe('RootTeamService', () => {
  let service: RootTeamService;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [RootTeamService],
    }).compile();

    service = module.get<RootTeamService>(RootTeamService);
  });

  it('should be defined', () => {
    expect(service).toBeDefined();
  });
});

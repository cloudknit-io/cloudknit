import { Test, TestingModule } from '@nestjs/testing';
import { RootTeamController } from './root.team.controller';
import { RootTeamService } from './root.team.service';

describe('RootTeamController', () => {
  let controller: RootTeamController;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [RootTeamController],
      providers: [RootTeamService],
    }).compile();

    controller = module.get<RootTeamController>(RootTeamController);
  });

  it('should be defined', () => {
    expect(controller).toBeDefined();
  });
});

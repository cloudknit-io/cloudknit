import { Test, TestingModule } from '@nestjs/testing';
import { RootEnvironmentController } from './root.environment.controller';
import { RootEnvironmentService } from './root.environment.service';

describe('RootEnvironmentController', () => {
  let controller: RootEnvironmentController;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [RootEnvironmentController],
      providers: [RootEnvironmentService],
    }).compile();

    controller = module.get<RootEnvironmentController>(RootEnvironmentController);
  });

  it('should be defined', () => {
    expect(controller).toBeDefined();
  });
});

import { Test, TestingModule } from '@nestjs/testing';
import { OperationsController } from './operations.controller';

describe('OperationsController', () => {
  let controller: OperationsController;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [OperationsController],
    }).compile();

    controller = module.get<OperationsController>(OperationsController);
  });

  it('should be defined', () => {
    expect(controller).toBeDefined();
  });
});

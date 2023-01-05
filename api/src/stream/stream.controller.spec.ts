import { Test, TestingModule } from '@nestjs/testing';
import { StreamController } from './stream.controller';
import { StreamService } from './stream.service';

describe('StreamController', () => {
  let controller: StreamController;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [StreamController],
      providers: [StreamService],
    }).compile();

    controller = module.get<StreamController>(StreamController);
  });

  it('should be defined', () => {
    expect(controller).toBeDefined();
  });
});

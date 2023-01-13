import { Test, TestingModule } from '@nestjs/testing';
import { ComponentService } from './component.service';

describe('ComponentService', () => {
  let service: ComponentService;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [ComponentService],
    }).compile();

    service = module.get<ComponentService>(ComponentService);
  });

  it('should be defined', () => {
    expect(service).toBeDefined();
  });
});

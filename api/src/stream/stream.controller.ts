import { Controller } from '@nestjs/common';
import { ApiTags } from '@nestjs/swagger';
import { StreamService } from './stream.service';

@Controller({
  version: '1',
})
@ApiTags('stream')
export class StreamController {
  constructor(private readonly sseSvc: StreamService) {}
}

interface MessageEvent {
  data: string | object;
  id?: string;
  type?: string;
  retry?: number;
}

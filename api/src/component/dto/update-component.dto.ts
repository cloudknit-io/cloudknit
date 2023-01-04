import { PartialType } from '@nestjs/swagger';
import { CreateComponentDto } from './create-component.dto';

export class UpdateComponentDto extends PartialType(CreateComponentDto) {
  status?: string;
  duration?: number;
  isDestroyed?: boolean;
}

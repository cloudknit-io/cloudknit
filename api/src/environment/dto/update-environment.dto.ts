import { PartialType } from '@nestjs/swagger';
import { CreateEnvironmentDto } from './create-environment.dto';
import { EnvSpecComponentDto } from './env-spec.dto';

export class UpdateEnvironmentDto extends PartialType(CreateEnvironmentDto) {
  name: string;
  duration: number;
  isDeleted: boolean;
}

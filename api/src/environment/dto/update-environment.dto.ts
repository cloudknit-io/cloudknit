import { PartialType } from '@nestjs/swagger';
import { EnvironmentReconcile } from 'src/typeorm';
import { CreateEnvironmentDto } from './create-environment.dto';

export class UpdateEnvironmentDto extends PartialType(CreateEnvironmentDto) {
  isDeleted?: boolean;
  latestEnvRecon?: EnvironmentReconcile;
  lastReconcileDatetime?: string;
}

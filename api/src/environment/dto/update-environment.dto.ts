import { PartialType } from '@nestjs/swagger';
import { DagDto } from 'src/reconciliation/dtos/environment-dag.dto';
import { CreateEnvironmentDto } from './create-environment.dto';

export class UpdateEnvironmentDto extends PartialType(CreateEnvironmentDto) {
  name: string;
  duration: number;
  isDeleted: boolean;
  dag: DagDto;
}

import { PartialType } from '@nestjs/swagger';
import { CostResource } from 'src/costing/dtos/Resource.dto';
import { CreateComponentDto } from './create-component.dto';

export class UpdateComponentDto extends PartialType(CreateComponentDto) {
  status?: string;
  duration?: number;
  lastWorkflowRunId?: number;
  estimatedCost?: number;
  costResources?: CostResource[];
  isDestroyed?: boolean;
}

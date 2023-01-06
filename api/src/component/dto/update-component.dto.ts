import { PartialType, ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';
import { CreateComponentDto } from './create-component.dto';

export class UpdateComponentDto extends PartialType(CreateComponentDto) {
  status?: string;
  duration?: number;
  lastWorkflowRunId?: string;
  estimatedCost?: number;
  costResources?: CostResource[];
  isDestroyed?: boolean;
}

export class CostResource {
  @ApiProperty()
  name: string
  @ApiPropertyOptional()
  hourlyCost?: string
  @ApiPropertyOptional()
  monthlyCost?: string
  @ApiPropertyOptional()
  subresources?: CostResource[]
  @ApiPropertyOptional()
  costComponents?: CostComponent[]
  @ApiPropertyOptional()
  metadata?: object
}

export class CostComponent {
  @ApiProperty()
  name: string
  @ApiProperty()
  price: string
  @ApiProperty()
  unit: string
  @ApiPropertyOptional()
  hourlyCost?: string
  @ApiPropertyOptional()
  hourlyQuantity?: string
  @ApiPropertyOptional()
  monthlyCost?: string
  @ApiPropertyOptional()
  monthlyQuantity?: string
}

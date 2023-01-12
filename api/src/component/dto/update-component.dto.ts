import { PartialType, ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';
import { Type } from 'class-transformer';
import { IsBoolean, IsNumber, IsObject, IsOptional, IsString, ValidateNested } from 'class-validator';
import { CreateComponentDto } from './create-component.dto';

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

export class CostResource {
  @ApiProperty()
  @IsString()
  name: string

  @ApiPropertyOptional()
  @IsOptional()
  @IsString()
  hourlyCost?: string

  @ApiPropertyOptional()
  @IsOptional()
  @IsString()
  monthlyCost?: string

  @ApiPropertyOptional()
  @IsOptional()
  @ValidateNested()
  @Type(() => CostResource)
  subresources?: CostResource[]

  @ApiPropertyOptional({ type: [CostComponent] })
  @IsOptional()
  @ValidateNested()
  @Type(() => CostComponent)
  costComponents?: CostComponent[]

  @ApiPropertyOptional()
  @IsOptional()
  @IsObject()
  metadata?: object
}
export class UpdateComponentDto extends PartialType(CreateComponentDto) {
  @ApiPropertyOptional()
  @IsOptional()
  @IsString()
  status?: string;

  @ApiPropertyOptional()
  @IsOptional()
  @IsNumber({
    maxDecimalPlaces: 0
  })
  duration?: number;

  @ApiPropertyOptional()
  @IsOptional()
  @IsString()
  lastWorkflowRunId?: string;

  @ApiPropertyOptional()
  @IsOptional()
  @IsNumber({
    maxDecimalPlaces: 5
  })
  estimatedCost?: number;

  @ApiPropertyOptional({ type: [CostResource] })
  @IsOptional()
  @ValidateNested()
  @Type(() => CostResource)
  costResources?: CostResource[];

  @ApiPropertyOptional()
  @IsOptional()
  @IsBoolean()
  isDestroyed?: boolean;
}

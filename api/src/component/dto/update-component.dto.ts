import { PartialType, ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';
import { Type } from 'class-transformer';
import { IsBoolean, IsNumber, IsOptional, IsString, ValidateNested } from 'class-validator';
import { CreateComponentDto } from './create-component.dto';

export class UpdateComponentDto extends PartialType(CreateComponentDto) {
  @IsOptional()
  @IsString()
  status?: string;

  @IsOptional()
  @IsNumber({
    maxDecimalPlaces: 0
  })
  duration?: number;

  @IsOptional()
  @IsString()
  lastWorkflowRunId?: string;

  @IsOptional()
  @IsNumber({
    maxDecimalPlaces: 5
  })
  estimatedCost?: number;

  @IsOptional()
  @ValidateNested()
  @Type(() => CostResource)
  costResources?: CostResource[];

  @IsOptional()
  @IsBoolean()
  isDestroyed?: boolean;
}

export class CostResource {
  @IsString()
  name: string

  @IsOptional()
  @IsString()
  hourlyCost?: string

  @IsOptional()
  @IsString()
  monthlyCost?: string

  @IsOptional()
  @ValidateNested()
  @Type(() => CostResource)
  subresources?: CostResource[]

  @IsOptional()
  @ValidateNested()
  @Type(() => CostComponent)
  costComponents?: CostComponent[]

  @IsOptional()
  @IsString()
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

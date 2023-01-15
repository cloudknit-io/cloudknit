import { ApiProperty, PartialType } from '@nestjs/swagger';
import { CreateTeamDto } from './create-team.dto';

export class UpdateTeamDto extends PartialType(CreateTeamDto) {
  @ApiProperty({ required: false })
  isDeleted?: boolean = false;

  @ApiProperty({ required: false })
  estimatedCost?: number;
}

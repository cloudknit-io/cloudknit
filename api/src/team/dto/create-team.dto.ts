import { ApiProperty } from '@nestjs/swagger';
import { Organization } from 'src/typeorm';

export class CreateTeamDto {
  organization: Organization;

  @ApiProperty({ required: false })
  name: string;

  @ApiProperty({ required: false })
  repo: string;

  @ApiProperty({ required: false })
  repo_path: string;

  @ApiProperty({ required: false })
  teardownProtection: boolean;
}

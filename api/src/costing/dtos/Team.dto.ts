import { ApiProperty } from "@nestjs/swagger";
import { EnvironmentDto } from "./Environment.dto";

export class TeamDto {
  @ApiProperty({
    name: 'Team Name',
  })
  teamName: string = '';

  @ApiProperty({
    name: 'Environments',
  })
  environments: EnvironmentDto[] = [];
}

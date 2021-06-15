import { ApiProperty } from '@nestjs/swagger'
import { ComponentDto } from './Component.dto'
import { TeamDto } from './Team.dto'

export class EnvironmentDto {
  @ApiProperty({
    name: 'Environment Name',
  })
  environmentName: string = ''

  @ApiProperty()
  team: TeamDto = null;

  @ApiProperty({
    name: 'Components',
  })
  components: ComponentDto[] = []
}

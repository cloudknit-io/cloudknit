import { ApiProperty } from '@nestjs/swagger'

export class CreateOrganizationDto {
  @ApiProperty()
  name: string

  @ApiProperty()
  clientId?: string

  @ApiProperty()
  clientSecret?: string
  
  @ApiProperty()
  githubRepo?: string

  @ApiProperty()
  githubPath?: string

  @ApiProperty()
  githubSource?: string
}

import { ApiProperty } from '@nestjs/swagger'

export class CreateOrganizationDto {
  @ApiProperty()
  name: string
  
  @ApiProperty()
  githubRepo: string

  @ApiProperty()
  termsAgreedUserId?: number
}

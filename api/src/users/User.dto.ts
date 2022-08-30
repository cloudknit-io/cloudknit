import { ApiProperty } from '@nestjs/swagger'

export class CreateUserDto {
  @ApiProperty()
  username: string

  @ApiProperty()
  email: string

  @ApiProperty({
    default: false
  })
  termAgreementStatus?: boolean

  @ApiProperty({
    default: 'User'
  })
  role: string
}

import { ApiProperty } from '@nestjs/swagger';

export class CreateUserDto {
  @ApiProperty()
  username: string;

  @ApiProperty()
  email: string;

  @ApiProperty({
    default: 'User',
  })
  role: string;

  @ApiProperty()
  name: string;
}

export class PatchUserDto {
  @ApiProperty({
    default: null,
  })
  email: string;

  @ApiProperty({
    default: null,
  })
  role: string;

  @ApiProperty({
    default: null,
  })
  name: string;
}

import { ApiProperty } from '@nestjs/swagger';
import { UserRole } from 'src/types';

export class CreateUserDto {
  @ApiProperty()
  username: string;

  @ApiProperty()
  email: string;

  @ApiProperty({
    default: UserRole.USER,
  })
  role: UserRole;

  @ApiProperty()
  name: string;
}

export class CreatePlaygroundUserDto {
  @ApiProperty()
  ipv4: string;
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

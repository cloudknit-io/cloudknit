import { ApiProperty } from '@nestjs/swagger'
import { IsNotEmpty, IsNumber, IsOptional, IsString } from 'class-validator'

export class CreateOrganizationDto {
  @IsNotEmpty()
  @IsString()
  @ApiProperty()
  name: string
  
  @IsNotEmpty()
  @IsString()
  @ApiProperty()
  githubRepo: string

  @IsOptional()
  @IsNumber()
  @ApiProperty()
  termsAgreedUserId?: number

  @IsOptional()
  @IsString()
  githubOrgName?: string
}


export class PatchOrganizationDto {
    @ApiProperty()
    githubRepo?: string

    @ApiProperty()
    provisioned?: boolean
}

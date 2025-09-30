import { IsNotEmpty, IsString } from 'class-validator';
import { ApiProperty } from '@nestjs/swagger';

export class LoginUserDto {
  @ApiProperty({
    example: 'johndoe',
    description: 'Username dell\'utente',
  })
  @IsNotEmpty()
  @IsString()
  username: string;

  @ApiProperty({
    example: 'Password123!',
    description: 'Password dell\'utente',
  })
  @IsNotEmpty()
  @IsString()
  password: string;
}
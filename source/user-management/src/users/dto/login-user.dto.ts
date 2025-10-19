import { IsNotEmpty, IsString } from 'class-validator';
import { ApiProperty } from '@nestjs/swagger';

export class LoginUserDto {
  @ApiProperty({
    example: 'alessandro',
    description: 'Username dell\'utente',
  })
  @IsNotEmpty()
  @IsString()
  username: string;

  @ApiProperty({
    example: 'alessandro123',
    description: 'Password dell\'utente',
  })
  @IsNotEmpty()
  @IsString()
  password: string;
}
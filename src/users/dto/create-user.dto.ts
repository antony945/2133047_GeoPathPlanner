import { IsEmail, IsNotEmpty, IsString, MinLength } from 'class-validator';
import { ApiProperty } from '@nestjs/swagger';

export class CreateUserDto {
  @ApiProperty({
    example: 'johndoe',
    description: 'Username univoco dell\'utente',
  })
  @IsNotEmpty()
  @IsString()
  username: string;

  @ApiProperty({
    example: 'Mario',
    description: 'Nome dell\'utente',
  })
  @IsNotEmpty()
  @IsString()
  nome: string;

  @ApiProperty({
    example: 'Rossi',
    description: 'Cognome dell\'utente',
  })
  @IsNotEmpty()
  @IsString()
  cognome: string;

  @ApiProperty({
    example: 'mario.rossi@example.com',
    description: 'Email dell\'utente',
  })
  @IsNotEmpty()
  @IsEmail()
  email: string;

  @ApiProperty({
    example: 'Password123!',
    description: 'Password dell\'utente (minimo 6 caratteri)',
  })
  @IsNotEmpty()
  @IsString()
  @MinLength(6)
  password: string;
}
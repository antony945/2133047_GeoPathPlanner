import { IsEmail, IsOptional, IsString } from 'class-validator';
import { ApiPropertyOptional } from '@nestjs/swagger';

export class UpdateUserDto {
  @ApiPropertyOptional({
    example: 'Alessandro',
    description: 'Nome dell\'utente',
  })
  @IsOptional()
  @IsString()
  nome?: string;

  @ApiPropertyOptional({
    example: 'Colantuoni',
    description: 'Cognome dell\'utente',
  })
  @IsOptional()
  @IsString()
  cognome?: string;

  @ApiPropertyOptional({
    example: 'ale@gmail.com',
    description: 'Email dell\'utente',
  })
  @IsOptional()
  @IsEmail()
  email?: string;

  @ApiPropertyOptional({
    example: 'Italia',
    description: 'Paese dell\'utente',
  })
  @IsOptional()
  @IsString()
  country?: string;
}
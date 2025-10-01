import { IsNotEmpty, IsString, MinLength } from 'class-validator';
import { ApiProperty } from '@nestjs/swagger';

export class ChangePasswordDto {
  @ApiProperty({
    example: 'OldPassword123!',
    description: 'Password attuale dell\'utente',
  })
  @IsNotEmpty()
  @IsString()
  oldPassword: string;

  @ApiProperty({
    example: 'NewPassword123!',
    description: 'Nuova password (minimo 6 caratteri)',
  })
  @IsNotEmpty()
  @IsString()
  @MinLength(6)
  newPassword: string;
}
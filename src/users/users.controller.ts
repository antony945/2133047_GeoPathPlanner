import { Controller, Get, Put, Patch, Param, Body, UseGuards, ClassSerializerInterceptor, UseInterceptors, Request } from '@nestjs/common';
import { ApiTags, ApiOperation, ApiResponse, ApiBearerAuth } from '@nestjs/swagger';
import { UsersService } from './users.service';
import { JwtAuthGuard } from '../auth/guards/jwt-auth.guard';
import { UpdateUserDto } from './dto/update-user.dto';
import { ChangePasswordDto } from './dto/change-password.dto';

@ApiTags('users')
@Controller('users')
@UseInterceptors(ClassSerializerInterceptor) // Per escludere la password dalle risposte
export class UsersController {
  constructor(private readonly usersService: UsersService) {}

  @Get()
  @UseGuards(JwtAuthGuard)
  @ApiBearerAuth()
  @ApiOperation({ summary: 'Ottieni tutti gli utenti' })
  @ApiResponse({ 
    status: 200, 
    description: 'Lista di tutti gli utenti',
    schema: {
      type: 'array',
      items: {
        type: 'object',
        properties: {
          id: { type: 'string' },
          username: { type: 'string' },
          nome: { type: 'string' },
          cognome: { type: 'string' },
          email: { type: 'string' },
          country: { type: 'string' },
          createdAt: { type: 'string', format: 'date-time' },
          updatedAt: { type: 'string', format: 'date-time' },
        },
      },
    },
  })
  @ApiResponse({ status: 401, description: 'Non autorizzato' })
  async findAll() {
    return this.usersService.findAll();
  }

  @Get(':id')
  @UseGuards(JwtAuthGuard)
  @ApiBearerAuth()
  @ApiOperation({ summary: 'Ottieni un utente specifico tramite ID' })
  @ApiResponse({ 
    status: 200, 
    description: 'Dettagli dell\'utente',
    schema: {
      type: 'object',
      properties: {
        id: { type: 'string' },
        username: { type: 'string' },
        nome: { type: 'string' },
        cognome: { type: 'string' },
        email: { type: 'string' },
        country: { type: 'string' },
        createdAt: { type: 'string', format: 'date-time' },
        updatedAt: { type: 'string', format: 'date-time' },
      },
    },
  })
  @ApiResponse({ status: 401, description: 'Non autorizzato' })
  @ApiResponse({ status: 404, description: 'Utente non trovato' })
  async findOne(@Param('id') id: string) {
    return this.usersService.findOne(id);
  }

  @Put(':id')
  @UseGuards(JwtAuthGuard)
  @ApiBearerAuth()
  @ApiOperation({ summary: 'Modifica i dati di un utente' })
  @ApiResponse({ 
    status: 200, 
    description: 'Utente modificato con successo',
    schema: {
      type: 'object',
      properties: {
        id: { type: 'string' },
        username: { type: 'string' },
        nome: { type: 'string' },
        cognome: { type: 'string' },
        email: { type: 'string' },
        country: { type: 'string' },
        createdAt: { type: 'string', format: 'date-time' },
        updatedAt: { type: 'string', format: 'date-time' },
      },
    },
  })
  @ApiResponse({ status: 401, description: 'Non autorizzato' })
  @ApiResponse({ status: 404, description: 'Utente non trovato' })
  @ApiResponse({ status: 409, description: 'Email già in uso' })
  async update(@Param('id') id: string, @Body() updateUserDto: UpdateUserDto) {
    return this.usersService.update(id, updateUserDto);
  }

  @Patch(':id/password')
  @UseGuards(JwtAuthGuard)
  @ApiBearerAuth()
  @ApiOperation({ summary: 'Modifica la password di un utente' })
  @ApiResponse({ 
    status: 200, 
    description: 'Password modificata con successo',
    schema: {
      type: 'object',
      properties: {
        message: { type: 'string', example: 'Password modificata con successo' }
      }
    }
  })
  @ApiResponse({ status: 401, description: 'Password attuale non corretta o non autorizzato' })
  @ApiResponse({ status: 404, description: 'Utente non trovato' })
  async changePassword(
    @Param('id') id: string, 
    @Body() changePasswordDto: ChangePasswordDto
  ) {
    return this.usersService.changePassword(id, changePasswordDto);
  }
}
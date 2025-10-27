import { Injectable, ConflictException, NotFoundException, UnauthorizedException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { User } from './entities/user.entity';
import { CreateUserDto } from './dto/create-user.dto';
import { UpdateUserDto } from './dto/update-user.dto';
import { ChangePasswordDto } from './dto/change-password.dto';
import * as bcrypt from 'bcrypt';

@Injectable()
export class UsersService {
  constructor(
    @InjectRepository(User)
    private usersRepository: Repository<User>,
  ) {}

  async create(createUserDto: CreateUserDto): Promise<User> {
    // Verifica se username o email esistono già
    const existingUser = await this.usersRepository.findOne({
      where: [
        { username: createUserDto.username },
        { email: createUserDto.email }
      ]
    });

    if (existingUser) {
      if (existingUser.username === createUserDto.username) {
        throw new ConflictException('Username già in uso');
      }
      if (existingUser.email === createUserDto.email) {
        throw new ConflictException('Email già in uso');
      }
    }

    // Hash della password
    const hashedPassword = await bcrypt.hash(createUserDto.password, 10);

    // Crea nuovo utente
    const user = this.usersRepository.create({
      ...createUserDto,
      password: hashedPassword,
    });

    return await this.usersRepository.save(user);
  }

  async findAll(): Promise<User[]> {
    return await this.usersRepository.find();
  }

  async findOne(id: string): Promise<User> {
    const user = await this.usersRepository.findOne({
      where: { id }
    });

    if (!user) {
      throw new NotFoundException(`Utente con ID ${id} non trovato`);
    }

    return user;
  }

  async update(id: string, updateUserDto: UpdateUserDto): Promise<User> {
    const user = await this.findOne(id);

    // Se si sta cambiando l'email, verifica che non esista già
    if (updateUserDto.email && updateUserDto.email !== user.email) {
      const existingUser = await this.usersRepository.findOne({
        where: { email: updateUserDto.email }
      });
      
      if (existingUser) {
        throw new ConflictException('Email già in uso');
      }
    }

    // Aggiorna i campi
    Object.assign(user, updateUserDto);
    
    return await this.usersRepository.save(user);
  }

  async changePassword(id: string, changePasswordDto: ChangePasswordDto): Promise<{ message: string }> {
    const user = await this.usersRepository.findOne({
      where: { id },
      select: ['id', 'password'] // Include la password per la verifica
    });

    if (!user) {
      throw new NotFoundException(`Utente con ID ${id} non trovato`);
    }

    // Verifica la password attuale
    const isPasswordValid = await bcrypt.compare(changePasswordDto.oldPassword, user.password);
    
    if (!isPasswordValid) {
      throw new UnauthorizedException('Password attuale non corretta');
    }

    // Hash della nuova password
    const hashedNewPassword = await bcrypt.hash(changePasswordDto.newPassword, 10);
    
    // Aggiorna la password
    user.password = hashedNewPassword;
    await this.usersRepository.save(user);

    return { message: 'Password modificata con successo' };
  }

  async findByUsername(username: string): Promise<User | undefined> {
    const user = await this.usersRepository.findOne({
      where: { username }
    });
    return user === null ? undefined : user;
  }

  async validateUser(username: string, password: string): Promise<User | null> {
    const user = await this.findByUsername(username);
    
    if (user && await bcrypt.compare(password, user.password)) {
      return user;
    }
    
    return null;
  }
}
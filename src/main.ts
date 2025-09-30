import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { ValidationPipe } from '@nestjs/common';
import { DocumentBuilder, SwaggerModule } from '@nestjs/swagger';
import { ConfigService } from '@nestjs/config';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);
  
  // Ottieni il servizio di configurazione
  const configService = app.get(ConfigService);
  
  // Abilita CORS (configura secondo le tue esigenze)
  app.enableCors();
  
  // Abilita la validazione globale
  app.useGlobalPipes(new ValidationPipe({
    whitelist: true, // Rimuove proprietà non definite nel DTO
    transform: true, // Trasforma automaticamente i tipi
    forbidNonWhitelisted: true, // Lancia un errore se ci sono proprietà extra
  }));
  
  // Configurazione Swagger
  const config = new DocumentBuilder()
    .setTitle('User Authentication API')
    .setDescription('API per la gestione di registrazione e autenticazione utenti')
    .setVersion('1.0')
    .addBearerAuth()
    .build();
    
  const document = SwaggerModule.createDocument(app, config);
  SwaggerModule.setup('api', app, document);
  
  // Avvia il server
  const port = configService.get('PORT') || 3000;
  await app.listen(port);
  
  console.log(`Application is running on: http://localhost:${port}`);
  console.log(`Swagger documentation: http://localhost:${port}/api`);
}
bootstrap();
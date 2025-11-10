import { Injectable, CanActivate, ExecutionContext, ForbiddenException } from '@nestjs/common';
import { Observable } from 'rxjs';

@Injectable()
export class UserOwnerGuard implements CanActivate {
  canActivate(
    context: ExecutionContext,
  ): boolean | Promise<boolean> | Observable<boolean> {
    const request = context.switchToHttp().getRequest();
    const user = request.user;
    const userIdFromParams = request.params.id;

    // Se non c'è un ID nei parametri, permetti l'accesso (per route come /me)
    if (!userIdFromParams) {
      return true;
    }

    // Verifica che l'utente stia accedendo solo ai propri dati
    if (user.userId !== userIdFromParams) {
      throw new ForbiddenException('Non hai i permessi per accedere a questa risorsa');
    }

    return true;
  }
}
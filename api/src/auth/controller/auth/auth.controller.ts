import { Controller, Get, Res } from '@nestjs/common';
import { Response } from 'express';

@Controller('auth')
export class AuthController {
  /*
   * user will login via this route
   */
  @Get('login')
  public login() {
    return;
  }

  /*
   * Redirect route that OAuth2 will call
   */
  @Get('redirect')
  public redirect(@Res() res: Response) {
    return res.send(200);
  }

  @Get('status')
  public status() {}

  @Get('logout')
  public logout() {}
}

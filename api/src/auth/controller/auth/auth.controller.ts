import { Controller, Get, Res } from '@nestjs/common';
import { Response } from 'express';

@Controller('auth')
export class AuthController {
  /*
   * user will login via this route
   */
  @Get('login')
  public login() {
    const k8s = require("@kubernetes/client-node");
    const kc = new k8s.KubeConfig();
    kc.loadFromCluster();
    const k8sApi = kc.makeApiClient(k8s.CoreV1Api);
    k8sApi
      .listNamespacedPod("zlifecycle-ui")
      .then((res) => {
        console.log(res.body.items[0].metadata);
      })
      .catch((err) => {
        console.log(err);
      });
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

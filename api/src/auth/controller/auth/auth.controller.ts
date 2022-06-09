import { Body, Controller, Get, Param, Post, Res } from "@nestjs/common";
import { Response } from "express";
import { AuthService } from "src/auth/services/auth/auth.service";

@Controller("auth")
export class AuthController {
  constructor(private readonly authService: AuthService) {}

  /*
   * user will login via this route
   */
  @Get("login/:accessKeyId/:secretAccessKey")
  public async login(
    @Param("accessKeyId") accessKeyId: string,
    @Param("secretAccessKey") secretAccessKey: string
  ) {
    console.log(Buffer.from(accessKeyId, "base64").toString());
    console.log(Buffer.from(secretAccessKey, "base64").toString());
    const separator = "[compuzest-shared]";
    const k8s = require("@kubernetes/client-node");
    const kc = new k8s.KubeConfig();
    kc.loadFromCluster();
    const k8sApi = kc.makeApiClient(k8s.CoreV1Api);

    const credentials = await k8sApi
      .readNamespacedSecret("aws-credentials-file", "argocd")
      .then((x) => x.body)
      .catch(() => null);

    const secret2 = await k8sApi
      .readNamespacedSecret("aws-creds", "argocd")
      .then((x) => x.body)
      .catch(() => null);

    const updates = [];
    updates.push(
      this.updateSecret(
        k8sApi,
        secret2,
        {
          aws_access_key_id: accessKeyId,
          naws_secret_access_key: secretAccessKey,
        },
        "argocd",
        "aws-creds"
      )
    );

    if (credentials) {
      console.log(credentials);
      console.log(credentials.data);
      console.log(credentials.data.credentials);
      const decoded = Buffer.from(
        credentials.data.credentials,
        "base64"
      ).toString();
      const splitTokens = decoded.split(separator);
      const updatedCreds = splitTokens[0].replace(
        /aws_access_key_id = \S+\naws_secret_access_key = \S+/,
        `aws_access_key_id = ${Buffer.from(
          accessKeyId,
          "base64"
        ).toString()}\naws_secret_access_key = ${Buffer.from(
          secretAccessKey,
          "base64"
        ).toString()}`
      );
      splitTokens[0] = updatedCreds;
      const encoded = Buffer.from(splitTokens.join(separator)).toString(
        "base64"
      );
      updates.push(
        this.updateSecret(
          k8sApi,
          credentials,
          { credentials: encoded },
          "argocd",
          "aws-credentials-file"
        )
      );
    }

    const res = await Promise.all(updates);

    return res;
  }

  @Post("termAgreementStatus")
  public async getTermAgreementStatus(@Body() body) {
    return this.authService.getTermAgreementStatus(body);
  }

  @Post("setTermAgreementStatus")
  public async setTermAgreementStatus(@Body() body) {
    return await this.authService.setTermAgreementStatus(body);
  }

  @Get("users/:organizationId")
  public async getUsers(@Param("organizationId") organizationId: string) {
    return this.authService.getUserList(organizationId);
  }

  @Get("user/:username")
  public async getUser(@Param("username") username: string) {
    return this.authService.getUser({ username });
  }

  @Post("add")
  public async addUser(@Body() body) {
    return await this.authService.addUser(body);
  }

  private async updateSecret(
    client: any,
    secret: any,
    inputs: any,
    namespace: string,
    name: string
  ) {
    if (!secret) {
      return client
        .createNamespacedSecret(namespace, {
          metadata: {
            name,
            namespace,
          },
          type: "Opaque",
          data: inputs,
        })
        .then((x) => x.body);
    }

    secret.data = inputs;

    return client
      .replaceNamespacedSecret(name, namespace, secret)
      .then((x) => x.body);
  }

  /*
   * Redirect route that OAuth2 will call
   */
  @Get("redirect")
  public redirect(@Res() res: Response) {
    return res.send(200);
  }

  @Get("status")
  public status() {}

  @Get("logout")
  public logout() {}
}

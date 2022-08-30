import { Body, Controller, Delete, Get, Logger, NotFoundException, Param, Post, Request, Res } from "@nestjs/common";
import { CreateUserDto } from "src/users/User.dto";
import { AuthService } from "./auth.service";

@Controller({
  version: '1'
})
export class AuthController {
  private readonly logger = new Logger(AuthController.name);

  constructor(private readonly authService: AuthService) {}

  /*
   * user will login via this route
   */
  @Get("login/:accessKeyId/:secretAccessKey")
  public async login(
    @Param("accessKeyId") accessKeyId: string,
    @Param("secretAccessKey") secretAccessKey: string
  ) {
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
  public async getTermAgreementStatus(@Request() req, @Body() body) {
    return this.authService.getTermAgreementStatus(req.org, body);
  }

  @Post("setTermAgreementStatus")
  public async setTermAgreementStatus(@Request() req, @Body() body) {
    return await this.authService.setTermAgreementStatus(req.org, body);
  }

  @Get("users")
  public async getUsers(@Request() req) {
    return this.authService.getOrgUserList(req.org);
  }

  @Get("users/:username")
  public async getUser(@Request() req, @Param("username") username: string) {
    const user = await this.authService.getOrgUser(req.org, username);

    if (!user) {
      throw new NotFoundException('could not find user');
    }

    return user
  }

  @Post("users")
  public async createUser(@Request() req, @Body() user: CreateUserDto) {
    return await this.authService.createOrgUser(req.org, user);
  }

  @Delete("users/:username")
  public async deleteUser(@Request() req, @Param('username') username: string) {
    return await this.authService.deleteUser(req.org, username);
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
}

import { Controller, Get, Param, Res } from "@nestjs/common";
import { Response } from "express";

@Controller("auth")
export class AuthController {
  /*
   * user will login via this route
   */
  @Get("login")
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
      const decoded = atob(credentials.data.credentials);
      const splitTokens = decoded.split(separator);
      const updatedCreds = splitTokens[0].replace(
        /aws_access_key_id = \S+\naws_secret_access_key = \S+/,
        `aws_access_key_id = ${atob(
          accessKeyId
        )}\naws_secret_access_key = ${atob(secretAccessKey)}`
      );
      splitTokens[0] = updatedCreds;
      const encoded = btoa(splitTokens.join(separator));
      updates.push(this.updateSecret(
        k8sApi,
        credentials,
        { credentials: encoded },
        "argocd",
        "aws-credentials-file"
      ));
    }

    const res = await Promise.all(updates);

    return res;
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

import { Injectable } from "@nestjs/common";
import { AWSError } from "aws-sdk";
import { GetParametersByPathRequest } from "aws-sdk/clients/ssm";
import { AwsSecretDto } from "./dtos/aws-secret.dto";
import { AWSSSMHandler } from "./utilities/awsSsmHandler";

@Injectable()
export class SecretsService {
  awsSecretSeparator = "[compuzest-shared]";
  k8sApi = null;
  ssm: AWSSSMHandler = null;
  constKeys = new Set([
    "aws_access_key_id",
    "aws_secret_access_key",
    "aws_session_token",
  ]);

  constructor() {
    const k8s = require("@kubernetes/client-node");
    const kc = new k8s.KubeConfig();
    kc.loadFromCluster();
    this.k8sApi = kc.makeApiClient(k8s.CoreV1Api);
    this.ssm = AWSSSMHandler.instance();
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

  private stringToBase64(value: string) {
    return Buffer.from(value).toString("base64");
  }

  private base64ToString(value: string) {
    return Buffer.from(value, "base64").toString();
  }

  private isConstKey(name: string) {
    const lastToken = name.split("/").slice(-1);
    return this.constKeys.has(lastToken[0]);
  }

  private mapToKeyValue(data: any) {
    const { Name } = data;
    const tokens = Name.split("/");
    let key = "";
    switch (tokens.length) {
      case 3:
        key = tokens[1];
        break;
      case 4:
        key = `${tokens[1]}:${tokens[2]}`;
        break;
      case 5:
        key = `${tokens[1]}:${tokens[2]}:${tokens[3]}`;
        break;
    }
    return {
      key,
      value: tokens.slice(-1)[0],
    };
  }

  private mapToEnvironments(data: any) {
    const { Name } = data;
    const tokens = Name.split("/");
    if (tokens.length === 5) {
      return [tokens[3], tokens[2]];
    }
    return null;
  }

  private getCredentialFileInput(
    credentials: string,
    accessKeyId: string,
    secretAccessKey: string
  ) {
    const decoded = this.base64ToString(credentials);
    const splitTokens = decoded.split(this.awsSecretSeparator);
    const updatedCreds = splitTokens[0].replace(
      /aws_access_key_id = \S+\naws_secret_access_key = \S+/,
      `aws_access_key_id = ${this.base64ToString(
        accessKeyId
      )}\naws_secret_access_key = ${this.base64ToString(secretAccessKey)}`
    );
    splitTokens[0] = updatedCreds;
    return this.stringToBase64(splitTokens.join(this.awsSecretSeparator));
  }

  public async createOrUpdateSecret(
    accessKeyId: string,
    secretAccessKey: string
  ) {
    const credentials = await this.k8sApi
      .readNamespacedSecret("aws-credentials-file", "argocd")
      .then((x) => x.body)
      .catch(() => null);

    const secret2 = await this.k8sApi
      .readNamespacedSecret("aws-creds", "argocd")
      .then((x) => x.body)
      .catch(() => null);

    const updates = [];
    updates.push(
      this.updateSecret(
        this.k8sApi,
        secret2,
        {
          aws_access_key_id: accessKeyId,
          aws_secret_access_key: secretAccessKey,
        },
        "argocd",
        "aws-creds"
      )
    );

    if (credentials) {
      const encoded = this.getCredentialFileInput(
        credentials.data.credentials,
        accessKeyId,
        secretAccessKey
      );
      updates.push(
        this.updateSecret(
          this.k8sApi,
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

  public async secretExist() {
    const credentials = await this.k8sApi
      .readNamespacedSecret("aws-credentials-file", "argocd")
      .then((x) => x.body)
      .catch(() => null);

    const secret2 = await this.k8sApi
      .readNamespacedSecret("aws-creds", "argocd")
      .then((x) => x.body)
      .catch(() => null);

    if (credentials && secret2) {
      return true;
    }
    return false;
  }

  public async ssmSecretExists(pathName: string) {
    try {
      const awsRes = await this.ssm.getParameter({
        Name: pathName,
      });
      return true;
    } catch (err) {
      const e = err as AWSError;
      if (e.code === "ParameterNotFound") {
        return false;
      } else {
        throw err;
      }
    }
  }

  public async ssmSecretsExists(pathNames: string[]) {
    try {
      const awsRes = await this.ssm.getParameters({
        Names: pathNames,
      });
      const resp = [];

      resp.push(
        ...awsRes.Parameters.map((e) => ({
          key: e.Name.split("/").slice(-1)[0],
          exists: true,
        }))
      );

      resp.push(
        ...awsRes.InvalidParameters.map((e) => ({
          key: e.split("/").slice(-1)[0],
          exists: false,
        }))
      );
      console.log(resp);
      return resp;
    } catch (err) {
      const e = err as AWSError;
      if (e.code === "ParameterNotFound") {
        return false;
      } else {
        throw err;
      }
    }
  }

  public async getSsmSecretsByPath(path: string) {
    try {
      const req: GetParametersByPathRequest = {
        Path: path,
        WithDecryption: false,
        Recursive: false,
      };
      const awsRes = await this.ssm.getParametersByPath(req);
      return awsRes.Parameters.filter((e) => !this.isConstKey(e.Name)).map(
        (e) => this.mapToKeyValue(e)
      );
    } catch (err) {
      const e = err as AWSError;
      if (e.code === "ParameterNotFound") {
        return false;
      } else {
        throw err;
      }
    }
  }

  public async getEnvironments(
    path: string,
    environments: Map<string, string> = new Map<string, string>(),
    nextToken: string = null
  ) {
    try {
      const req: GetParametersByPathRequest = {
        Path: path,
        WithDecryption: false,
        Recursive: true,
        NextToken: nextToken,
      };
      const awsRes = await this.ssm.getParametersByPath(req);
      awsRes.Parameters.forEach((e) => {
        const env = this.mapToEnvironments(e);
        if (env) {
          environments.set(env[0], env[1]);
        }
      });
      nextToken = awsRes.NextToken;
      if (nextToken) {
        await this.getEnvironments(path, environments, nextToken);
      }
      return [...environments.entries()];
    } catch (err) {
      const e = err as AWSError;
      if (e.code === "ParameterNotFound") {
        return false;
      } else {
        throw err;
      }
    }
  }

  public async putSsmSecrets(awsSecrets: AwsSecretDto[]) {
    const awsCalls = awsSecrets.map((secret) =>
      this.putSsmSecret(secret.path, secret.value, "SecureString")
    );
    const responses = await Promise.all(awsCalls);
    return !responses.some((response) => response === false);
  }

  public async putSsmSecret(
    pathName: string,
    value: string,
    type: "SecureString" | "StringList" | "String"
  ): Promise<boolean> {
    try {
      const awsRes = await this.ssm.putParameter({
        Name: pathName,
        Value: value,
        Overwrite: true,
        Type: type,
      });
      return true;
    } catch (err) {
      const e = err as AWSError;
      if (e.code === "ParameterNotFound") {
        return false;
      } else {
        throw err;
      }
    }
  }

  public async deleteSSMSecret(path: string) {
    try {
      const dp = await this.ssm.deleteParameter({
        Name: path,
      });
      console.log(dp);
      return dp;
    } catch (err) {
      console.log(err);
      const e = err as AWSError;
      if (e.code === "ParameterNotFound") {
        return false;
      } else {
        throw err;
      }
    }
  }
}

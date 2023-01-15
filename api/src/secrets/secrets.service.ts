import { Injectable, Logger } from '@nestjs/common';
import { AWSError } from 'aws-sdk';
import {
  GetParameterRequest,
  GetParametersByPathRequest,
  Parameter,
} from 'aws-sdk/clients/ssm';
import { Organization } from 'src/typeorm';
import { AwsSecretDto } from './dtos/aws-secret.dto';
import { AWSSSMHandler } from './utilities/awsSsmHandler';
@Injectable()
export class SecretsService {
  private readonly logger = new Logger(SecretsService.name);

  awsSecretSeparator = '[compuzest-shared]';
  ssm: AWSSSMHandler = null;
  constKeys = new Set([
    'aws_access_key_id',
    'aws_secret_access_key',
    'aws_session_token',
    'state_aws_access_key_id',
    'state_aws_secret_access_key',
    'state_bucket',
    'state_lock_table',
  ]);

  constructor() {
    this.ssm = AWSSSMHandler.instance();
  }

  private isConstKey(name: string) {
    const lastToken = name.split('/').slice(-1);
    return this.constKeys.has(lastToken[0]);
  }

  private getPath(org: Organization, path: string) {
    if (path) {
      if (path[0] === '/') {
        path = path.slice(1);
      }

      return `/${org.name}/${path}`;
    } else {
      return `/${org.name}`;
    }
  }

  private mapToKeyValue(data: Parameter) {
    const { Name, LastModifiedDate } = data;

    let tokens = Name.split('/');
    tokens = tokens.slice(2); // removes org and preceding empty str
    const value = tokens.pop();

    return {
      key: tokens.join(':'),
      value,
      lastModifiedDate: LastModifiedDate,
    };
  }

  private mapToEnvironments(data: any) {
    const { Name } = data;
    const tokens = Name.split('/');
    if (tokens.length === 5) {
      return [tokens[3], tokens[2]];
    }
    return null;
  }

  public async ssmSecretsExists(org: Organization, pathNames: string[]) {
    try {
      const awsRes = await this.ssm.getParameters({
        Names: pathNames.map((path) => this.getPath(org, path)),
      });
      const resp = [];

      resp.push(
        ...awsRes.Parameters.map((e) => ({
          key: e.Name.split('/').slice(-1)[0],
          exists: true,
          lastModifiedDate: e.LastModifiedDate,
        }))
      );

      resp.push(
        ...awsRes.InvalidParameters.map((e) => ({
          key: e.split('/').slice(-1)[0],
          exists: false,
        }))
      );

      return resp;
    } catch (err) {
      const e = err as AWSError;
      if (e.code === 'ParameterNotFound') {
        return false;
      } else {
        throw err;
      }
    }
  }

  public async getSsmSecret(org: Organization, path: string): Promise<string> {
    try {
      const req: GetParameterRequest = {
        Name: this.getPath(org, path),
        WithDecryption: true,
      };

      const awsRes = await this.ssm.getParameter(req);

      return awsRes.Parameter.Value;
    } catch (err) {
      const e = err as AWSError;
      if (e.code === 'ParameterNotFound') {
        return null;
      } else {
        throw err;
      }
    }
  }

  public async getSsmSecretsByPath(org: Organization, path: string) {
    try {
      const req: GetParametersByPathRequest = {
        Path: this.getPath(org, path),
        WithDecryption: false,
        Recursive: false,
      };

      const awsRes = await this.ssm.getParametersByPath(req);

      return awsRes.Parameters.filter((e) => !this.isConstKey(e.Name)).map(
        (e) => this.mapToKeyValue(e)
      );
    } catch (err) {
      const e = err as AWSError;
      if (e.code === 'ParameterNotFound') {
        return false;
      } else {
        throw err;
      }
    }
  }

  public async getEnvironments(
    org: Organization,
    path: string,
    environments: Map<string, string> = new Map<string, string>(),
    nextToken: string = null
  ) {
    try {
      const req: GetParametersByPathRequest = {
        Path: this.getPath(org, path),
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
        await this.getEnvironments(org, path, environments, nextToken);
      }

      return [...environments.entries()];
    } catch (err) {
      const e = err as AWSError;

      if (e.code === 'ParameterNotFound') {
        return false;
      } else {
        throw err;
      }
    }
  }

  public async putSsmSecrets(org: Organization, awsSecrets: AwsSecretDto[]) {
    const awsCalls = awsSecrets.map((secret) =>
      this.putSsmSecret(org, secret.path, secret.value, 'SecureString')
    );

    const responses = await Promise.all(awsCalls);

    return !responses.some((response) => response === false);
  }

  public async putSsmSecret(
    org: Organization,
    pathName: string,
    value: string,
    type: 'SecureString' | 'StringList' | 'String'
  ): Promise<boolean> {
    try {
      const awsRes = await this.ssm.putParameter({
        Name: this.getPath(org, pathName),
        Value: value,
        Overwrite: true,
        Type: type,
      });
      return true;
    } catch (err) {
      const e = err as AWSError;
      if (e.code === 'ParameterNotFound') {
        return false;
      } else {
        throw err;
      }
    }
  }

  public async deleteSSMSecret(org: Organization, path: string) {
    try {
      const dp = await this.ssm.deleteParameter({
        Name: this.getPath(org, path),
      });
      return dp;
    } catch (err) {
      const e = err as AWSError;
      if (e.code === 'ParameterNotFound') {
        return false;
      } else {
        throw err;
      }
    }
  }
}

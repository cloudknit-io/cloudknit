import { Injectable, Logger } from "@nestjs/common";
import { AWSError } from "aws-sdk";
import { GetParameterRequest } from "aws-sdk/clients/ssm";
import { AWSSSMHandler } from "src/secrets/utilities/awsSsmHandler";

@Injectable()
export class SystemService {
  private readonly logger = new Logger(SystemService.name);

  ssm: AWSSSMHandler = null;

  constructor() {
    this.ssm = AWSSSMHandler.instance();
  }

  private normalizePath(path: string) : string {
    const systemPrefix = '/system';

    if (path[0] === "/") {
      path = path.slice(1);
    }

    return `${systemPrefix}/${path}`;
  }

  public async getSsmSecret(path: string) : Promise<string> {
    try {
      const req: GetParameterRequest = {
        Name: this.normalizePath(path),
        WithDecryption: true,
      };
      
      const awsRes = await this.ssm.getParameter(req);
      
      return awsRes.Parameter.Value;
    } catch (err) {
      const e = err as AWSError;
      if (e.code === "ParameterNotFound") {
        return null;
      } else {
        throw err;
      }
    }
  }
}

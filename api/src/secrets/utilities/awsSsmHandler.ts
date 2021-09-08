import { SSM } from "aws-sdk";

export class AWSSSMHandler {
  private static _instance: AWSSSMHandler = null;
  private ssm: SSM = null;
  constructor() {
    this.ssm = new SSM({
      accessKeyId: process.env.AWS_ACCESS_KEY_ID,
      secretAccessKey: process.env.AWS_SECRET_ACCESS_KEY,
      region: "us-east-1",
    });
  }

  static instance(): AWSSSMHandler {
    if (!AWSSSMHandler._instance) {
      AWSSSMHandler._instance = new AWSSSMHandler();
    }
    return AWSSSMHandler._instance;
  }

  async getParameter(
    request: SSM.GetParameterRequest
  ): Promise<SSM.GetParameterResult> {
    return new Promise<SSM.GetParameterResult>((done, error) => {
      this.ssm.getParameter(request, (err, data) => {
        if (err) {
          error(err);
        }
        done(data);
      });
    });
  }


  async getParameters(
    request: SSM.GetParametersRequest
  ): Promise<SSM.GetParametersResult> {
    return new Promise<SSM.GetParametersResult>((done, error) => {
      this.ssm.getParameters(request, (err, data) => {
        if (err) {
          error(err);
        }
        done(data);
      });
    });
  }

  async putParameter(
    request: SSM.PutParameterRequest
  ): Promise<SSM.PutParameterResult> {
    return new Promise<SSM.PutParameterResult>((done, error) => {
      this.ssm.putParameter(request, (err, data) => {
        if (err) {
          error(err);
        }
        done(data);
      });
    });
  }
}

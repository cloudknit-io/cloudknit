import { SSM } from 'aws-sdk';
import { get } from 'src/config';

export class AWSSSMHandler {
  private static _instance: AWSSSMHandler = null;
  private ssm: SSM = null;
  private readonly config = get();

  constructor() {
    this.ssm = new SSM({
      credentials: {
        accessKeyId: this.config.AWS.accessKeyId,
        secretAccessKey: this.config.AWS.secretAccessKey,
        sessionToken: this.config.AWS.sessionToken,
      },
      region: 'us-east-1',
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

  async getParametersByPath(
    request: SSM.GetParametersByPathRequest
  ): Promise<SSM.GetParametersByPathResult> {
    return new Promise<SSM.GetParametersByPathResult>((done, error) => {
      this.ssm.getParametersByPath(request, (err, data) => {
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

  async deleteParameter(request: SSM.DeleteParameterRequest) {
    return new Promise<SSM.DeleteParameterResult>((done, error) => {
      this.ssm.deleteParameter(request, (err, data) => {
        if (err) {
          error(err);
        }
        done(data);
      });
    });
  }
}

import { S3 } from 'aws-sdk'
import { get } from 'src/config';

export interface FileInfo {
    data?: S3.GetObjectOutput;
    key: string;
    error?: any;
}

export class S3Handler {
  private _s3 = null
  private readonly config = get();
  private static _instance = null

  private constructor() {
    this._s3 = new S3({
      accessKeyId: this.config.AWS.accessKeyId,
      secretAccessKey: this.config.AWS.secretAccessKey,
      region: 'us-east-1',
    })
  }

  static instance(): S3Handler {
    if (!this._instance) {
      this._instance = new S3Handler()
    }
    return this._instance
  }

  public get s3(): S3 {
    return this._s3
  }

  async getObjectStream(bucket: string, fileName: string) {
    return this.s3.getObject(
      {
        Bucket: bucket,
        Key: fileName,
      },
      (err, data) => {
      },
    ).createReadStream();
  }

  async getObject(bucket: string, fileName: string): Promise<FileInfo> {
    try {
      return await this.downloadFile(bucket, fileName);
    } catch (err) {
      return {
        key: fileName,
        error: err,
      };
    }
    
  }

  async getObjects(bucket: string, prefix: string): Promise<FileInfo[]> {
    try {
      const fileContents: S3.ObjectList = await this.getObjectList(
        bucket,
        prefix,
      )

      if (!fileContents || fileContents.length === 0) {
        throw '';
      }

      const filesRequests = []

      for (let i = 0; i < fileContents.length; i++) {
        const element = fileContents[i]
        filesRequests.push(this.downloadFile(bucket, element.Key))
      }

      return await Promise.all(filesRequests)
    } catch (err) {
      throw err;
    }
  }

  private async getObjectList(
    bucket: string,
    prefix: string,
  ): Promise<S3.ObjectList> {
    return new Promise<S3.ObjectList>((resolve, reject) => {
      this.s3.listObjectsV2(
        {
          Bucket: bucket,
          Prefix: prefix,
        },
        (err, data) => {
          if (err) {
            reject(err)
            return
          }
          resolve(data.Contents)
        },
      )
    })
  }

  public async copyToS3(bucket: string, path: string, contents: Express.Multer.File) {
    const uploadProcess = this.s3.upload({
      Bucket: bucket,
      Body: contents.buffer,
      Key: path
    });

    const response = await uploadProcess.promise();
    return response.Key;
  }


  private async downloadFile(bucket: string, key: string): Promise<FileInfo> {
    return new Promise<FileInfo>((res, rej) => {
      this.s3.getObject(
        {
          Bucket: bucket,
          Key: key,
        },
        (err, data) => {
          if (err) {
            rej(err)
          }
          const fileInfo: FileInfo = {
              data,
              key
          };
          res(fileInfo);
        },
      )
    })
  }
}

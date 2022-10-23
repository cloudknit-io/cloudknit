import { Logger } from "@nestjs/common";

export type ApiConfig = {
  TypeORM: {
    host: string,
    port: number,
    username: string,
    password: string,
    database: string,
  },
  port: number,
  AWS: {
    accessKeyId: string,
    secretAccessKey: string
  },
  environment: string,
  argo: {
    wf: {
      url: string,
      namespace: string
    }
  }
}

let config: ApiConfig;

function getEnvVarOrFail(varName: string): string {
  const v = process.env[varName];

  if (!v) {
    throw new Error(`could not find ${varName} env var`);
  }

  return v;
}

export function init() {
  const logger = new Logger('config');

  config = {
    TypeORM: {
      host: getEnvVarOrFail('TYPEORM_HOST'),
      port: parseInt(getEnvVarOrFail('TYPEORM_PORT')),
      username: getEnvVarOrFail('TYPEORM_USERNAME'),
      password: getEnvVarOrFail('TYPEORM_PASSWORD'),
      database: getEnvVarOrFail('TYPEORM_DATABASE'),
    },
    port: parseInt(process.env.APP_PORT) || 3000,
    AWS: {
      accessKeyId: getEnvVarOrFail('AWS_ACCESS_KEY_ID'),
      secretAccessKey: getEnvVarOrFail('AWS_SECRET_ACCESS_KEY'),
    },
    argo: {
      wf: {
        url: getEnvVarOrFail('CK_ARGO_WF_URL'),
        namespace: getEnvVarOrFail('CK_ARGO_WF_NAMESPACE'),
      }
    },
    environment: getEnvVarOrFail('CK_ENVIRONMENT'),
  };
  
  logger.log('successfully configured');
}

export function get() : ApiConfig {
  if (!config) {
    init();
  }

  return config;
}

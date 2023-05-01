import { Logger } from '@nestjs/common';

export type ApiConfig = {
  TypeORM: {
    host: string;
    port: number;
    username: string;
    password: string;
    database: string;
    sync: boolean;
  };
  port: number;
  AWS: {
    accessKeyId: string;
    secretAccessKey: string;
    sessionToken: string;
  };
  environment: string;
  isLocal: boolean;
  argo: {
    wf: {
      skipProvision: boolean;
      url: string;
      orgUrl: Function;
      namespace: string;
    },
    cd: {
      url: string
    }
  };
};

let config: ApiConfig;

function getEnvVarOrFail(varName: string): string {
  const v = process.env[varName];

  if (!v) {
    throw new Error(`could not find ${varName} env var`);
  }

  return v;
}

function getEnvVarOrDefault(varName: string, dfault: any) {
  try {
    const val = getEnvVarOrFail(varName);
    return val;
  } catch {
    return dfault;
  }
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
      sessionToken: getEnvVarOrDefault('AWS_SESSION_TOKEN', null)
    },
    argo: {
      wf: {
        skipProvision:
          getEnvVarOrDefault('CK_ARGO_WF_SKIP_PROVISION', 'true') === 'true',
        url: getEnvVarOrFail('CK_ARGO_WF_URL'),
        orgUrl: (orgName: string) =>
          getEnvVarOrFail('CK_ARGO_WF_ORG_URL').replaceAll(':org', orgName),
        namespace: getEnvVarOrFail('CK_ARGO_WF_NAMESPACE'),
      },
      cd: {
        url: getEnvVarOrDefault('CK_ARGO_CD_URL', 'http://argocd-zlifecycle-server.zlifecycle-system.svc.cluster.local'),
      }
    },
    environment: getEnvVarOrFail('CK_ENVIRONMENT'),
    isLocal: getEnvVarOrDefault('IS_LOCAL', 'false') === 'true',
  };

  logger.log('successfully configured');
}

export function get(): ApiConfig {
  if (!config) {
    init();
  }

  return config;
}

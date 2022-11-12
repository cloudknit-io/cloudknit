import * as k8s from "@kubernetes/client-node";
import { CustomObjectsApi, KubeConfig } from "@kubernetes/client-node";
import { ApiConfig, get } from "src/config";
import { WinstonLogger } from "src/logger";

class ApiKubeConfig
{
  private readonly _logger = new WinstonLogger();
  private static _instance: ApiKubeConfig;
  private kc: KubeConfig;
  private config: ApiConfig;
  private custObjApi: CustomObjectsApi;

  private constructor()
  {
    this.config = get();
    this.kc = new k8s.KubeConfig();

    if (this.config.isLocal) {
      this._logger.log('Loading KubeConfig from default');
      this.kc.loadFromDefault();
    } else {
      this._logger.log('Loading KubeConfig from cluster');
      this.kc.loadFromCluster();
    }

    this.custObjApi = this.kc.makeApiClient(k8s.CustomObjectsApi);      
  }

  public get kubeConfig() : KubeConfig {
    return this.kc;
  }

  public get customObjectApi() : CustomObjectsApi {
    return this.custObjApi;
  }

  public get logger() : WinstonLogger {
    return this._logger;
  }

  public static get Instance()
  {
    return this._instance || (this._instance = new this());
  }
}

export const Instance = ApiKubeConfig.Instance;

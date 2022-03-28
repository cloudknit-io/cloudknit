import {
  CoreV1Api,
  CustomObjectsApi,
  dumpYaml,
  KubeConfig,
  loadYaml,
  V1ConfigMap,
  V1Secret,
} from "@kubernetes/client-node";
import { Injectable } from "@nestjs/common";
import { InjectRepository } from "@nestjs/typeorm";
import { Company } from "src/typeorm/company/Company";
import { Repository } from "typeorm";

@Injectable()
export class CompanyService {
  private k8sApi: CoreV1Api;
  private k8sCRDApi: CustomObjectsApi;
  private debugMode = false;

  constructor(@InjectRepository(Company) private orgRepo: Repository<Company>) {
    const kc = new KubeConfig();
    if (this.debugMode) kc.loadFromDefault();
    else kc.loadFromCluster();
    this.k8sApi = kc.makeApiClient<CoreV1Api>(CoreV1Api);
    this.k8sCRDApi = kc.makeApiClient<CustomObjectsApi>(CustomObjectsApi);
  }

  async saveOAuthCredentials({ company, clientId, clientSecret }) {
    const savedData = await this.orgRepo.save({
      name: company,
      clientId,
      clientSecret,
    });

    return savedData;
  }

  async saveGitHubCredentials({
    company,
    githubRepo,
    githubPath,
    githubSource,
  }) {
    const orgData = await this.orgRepo.findOne({
      where: {
        name: company,
      },
    });
    orgData.githubPath = githubPath;
    orgData.githubRepo = githubRepo;
    orgData.githubSource = githubSource;
    const savedData = await this.orgRepo.save(orgData);
    return savedData;
  }

  async patchOrganisationData({ company, namespace }) {
    const orgData = await this.orgRepo.findOne({
      where: {
        name: company,
      },
    });
    if (!orgData) {
      throw `No data found for ${company}!`;
    }
    if (!this.k8sApi) {
      throw "Cannot initialize k8s API!";
    }
    await this.patchArgoCdConfig(orgData, namespace);
    await this.patchBffSecret(orgData, namespace);
    return orgData;
  }

  private async patchBffSecret({ clientSecret, name, clientId }, namespace) {
    const secrets: V1Secret = await this.k8sApi
      .readNamespacedSecret("zlifecycle-web-bff-development", namespace)
      .then((x) => x.body)
      .catch(() => null);
    if (!secrets) {
      throw "Error while updating credentials";
    }
    secrets.data["OPENID_CLIENT_ID"] =  Buffer.from(clientId).toString("base64");
    secrets.data["OPENID_CLIENT_SECRET"] = Buffer.from(clientSecret).toString("base64");
    const updateResponse = await this.k8sApi.replaceNamespacedSecret(
      "zlifecycle-web-bff-development",
      namespace,
      secrets
    );
  }

  private async patchArgoCdConfig({ name, clientSecret, clientId }, namespace) {
    const cm: V1ConfigMap = await this.k8sApi
      .readNamespacedConfigMap("argocd-cm", namespace)
      .then((x) => x.body)
      .catch(() => null);
    if (!cm) {
      throw "Error while updating credentials";
    }

    const dexConfig = loadYaml(cm.data["dex.config"]);
    dexConfig["connectors"][0]["config"]["clientID"] = clientId;
    dexConfig["connectors"][0]["config"]["clientSecret"] = clientSecret;
    dexConfig["connectors"][0]["config"]["orgs"][0]["name"] = name;
    dexConfig["staticClients"][0]["secret"] = clientSecret;
    const dexConfigYaml = dumpYaml(dexConfig);
    cm.data["dex.config"] = dexConfigYaml;
    const updateResponse = await this.k8sApi.replaceNamespacedConfigMap(
      "argocd-cm",
      namespace,
      cm
    );
  }

  public async patchCRD({ org, company }) {
    const orgData = await this.orgRepo.findOne({
      where: {
        name: org,
      },
    });
    console.log(orgData);
    if (!orgData) {
      return false;
    }

    const patch = {
      spec: {
        configRepo: {
          path: orgData.githubPath,
          source: orgData.githubSource,
        },
      },
    };

    return await this.k8sCRDApi
      .patchNamespacedCustomObject(
        "stable.compuzest.com",
        "v1",
        `${company}-config`,
        "companies",
        company,
        patch,
        undefined,
        undefined,
        undefined,
        {
          headers: {
            "content-type": "application/merge-patch+json",
          },
        }
      )
      .then((x) => true)
      .catch((e) => {
        console.log(e);
        return false;
      });
  }
}

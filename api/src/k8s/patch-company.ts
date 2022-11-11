import * as k8s from "@kubernetes/client-node";
import { Organization } from "src/typeorm";
import kc from "src/k8s/kc";


export async function patchCompany(org: Organization, githubRepo: string) {
  const lowerOrgName = org.name.toLowerCase();
  const options = { "headers": { "Content-type": k8s.PatchUtils.PATCH_FORMAT_JSON_PATCH}};
  const group = 'stable.compuzest.com';
  const version = 'v1';
  const plural = 'companies';
  const namespace = `${lowerOrgName}-config`;

  const patch = [{
    "op": "replace",
    "path":"/spec/configRepo/source",
    "value": githubRepo
  }];

  try {
    const res = await kc.customObjectApi.patchNamespacedCustomObject(group, version, namespace, plural, lowerOrgName, patch, undefined, undefined, undefined, options);
    kc.logger.log(`Successfully updated company resource for ${lowerOrgName}`)
  } catch (error) {
    const {body} = error;
    kc.logger.error({message: error.message, body});
    throw error;
  }
}

import { match, MatchFunction } from "path-to-regexp";

export interface PathMapping {
  pathMatch: MatchFunction<object>;
  // eslint-disable-next-line no-unused-vars
  newPath: (params: any) => string;
}

const CD_PATH_MAPPINGS = [
  {
    path: "/cd/api/v1/stream/environments/:environmentId",
    newPath: (params: any) =>
      `/api/v1/stream/applications?name=${params.orgName}-${params.environmentId}&projects=${params.orgName}`,
  },
  {
    path: "/cd/api/v1/projects/:environmentId/sync",
    newPath: (params: any) =>
      `/api/v1/applications/${params.orgName}-${params.environmentId}/sync`,
  },

  {
    path: "/cd/api/v1/projects/:environmentId/delete",
    newPath: (params: any) =>
      `/api/v1/applications/${params.orgName}-${params.environmentId}/resource?name=${params.orgName}-${params.environmentId}&namespace=${params.orgName}-executor&resourceName=${params.orgName}-${params.environmentId}&version=v1alpha1&kind=Workflow&group=argoproj.io&force=true&orphan=false`,
  },
  {
    path: "/cd/api/v1/projects/watcher/:projectId/sync",
    newPath: (params: any) =>
      `/api/v1/applications/${params.orgName}-${params.projectId}-team-watcher/sync`,
  },
];

const WF_PATH_MAPPINGS = [
  {
    path: "/wf/api/v1/projects/:projectId/environments/:environmentId/config/:configId/:workflowId",
    newPath: (params: any) =>
      `/api/v1/workflows/${params.orgName}-executor/${params.workflowId}`,
  },
  {
    path: "/wf/api/v1/projects/:projectId/environments/:environmentId/config/:configId/:workflowId/log",
    newPath: (params: any) =>
      `/api/v1/workflows/${params.orgName}-executor/${params.workflowId}/log?logOptions.container=main&logOptions.follow=false`,
  },
  {
    path: "/wf/api/v1/stream/projects/:projectId/environments/:environmentId/config/:configId/:workflowId",
    newPath: (params: any) =>
      `/api/v1/workflow-events/${params.orgName}-executor?listOptions.fieldSelector=metadata.name=${params.workflowId}`,
  },
  {
    path: "/wf/api/v1/stream/projects/:projectId/environments/:environmentId/config/:configId/:workflowId/log/:nodeId",
    newPath: (params: any) =>
      `/api/v1/workflows/${params.orgName}-executor/${params.workflowId}/log?logOptions.container=main&logOptions.follow=true&podName=${params.nodeId}`,
  },
];

const AUDIT_PATH_MAPPINGS = [
  {
    path: "/reconciliation/api/v1/component/:teamId/:environmentId/:componentId",
    newPath: (params: any) =>
      `/v1/orgs/${params.orgId}/teams/${params.teamId}/environments/${params.environmentId}/components/${params.componentId}/audit`,
  },
  {
    path: "/reconciliation/api/v1/environment/:teamId/:environmentId",
    newPath: (params: any) =>
    `/v1/orgs/${params.orgId}/teams/${params.teamId}/environments/${params.environmentId}/audit`,
  },
  {
    path: "/reconciliation/api/v1/getStateFile/:team/:environment/:component",
    newPath: (params: any) =>
      `/v1/orgs/${params.orgId}/reconciliation/component/state-file/${params.team}/${params.environment}/${params.component}`,
  },
  {
    path: "/reconciliation/api/v1/getLogs/:team/:environment/:component/:id",
    newPath: (params: any) =>
      `/v1/orgs/${params.orgId}/reconciliation/component/logs/${params.team}/${params.environment}/${params.component}/${params.id}`,
  },
  {
    path: "/reconciliation/api/v1/getPlanLogs/:team/:environment/:component/:id/:latest",
    newPath: (params: any) =>
      `/v1/orgs/${params.orgId}/reconciliation/component/plan/logs/${params.team}/${params.environment}/${params.component}/${params.id}/${params.latest}`,
  },
  {
    path: "/reconciliation/api/v1/getApplyLogs/:team/:environment/:component/:id/:latest",
    newPath: (params: any) =>
      `/v1/orgs/${params.orgId}/reconciliation/component/apply/logs/${params.team}/${params.environment}/${params.component}/${params.id}/${params.latest}`,
  },
  {
    path: "/reconciliation/api/v1/component/approve/:compReconId",
    newPath: (params: any) =>
      `/v1/orgs/${params.orgId}/reconciliation/component/${params.compReconId}/approve`,
  },
  {
    path: "/reconciliation/api/v1/approved-by/:teamId/:envId/:compId",
    newPath: (params: any) =>
      `/v1/orgs/${params.orgId}/reconciliation/approved-by/${params.teamId}/${params.envId}/${params.compId}/${params.email}`,
  }
];

const SECRET_PATH_MAPPINGS = [
  {
    path: "/secrets/default",
    newPath: (params: any) => `/v1/orgs/${params.orgId}/secrets/default`,
  },
  {
    path: "/secrets/exists/aws-secret",
    newPath: (params: any) => `/v1/orgs/${params.orgId}/secrets/exists/aws-secret`,
  },
  {
    path: "/secrets/update/aws-secret",
    newPath: (params: any) => `/v1/orgs/${params.orgId}/secrets/update/aws-secret`,
  },
  {
    path: "/secrets/get/environments",
    newPath: (params: any) => `/v1/orgs/${params.orgId}/secrets/environments`,
  },
  {
    path: "/secrets/get/ssm-secrets",
    newPath: (params: any) => `/v1/orgs/${params.orgId}/secrets/get/ssm-secrets/`,
  },
  {
    path: "/secrets/delete/ssm-secret",
    newPath: (params: any) => `/v1/orgs/${params.orgId}/secrets/delete/ssm-secret`,
  },
];

const STATE_PATH_MAPPINGS = [
  {
    path: "/terraform/state",
    newPath: () => `/terraform/state`,
  },
  {
    path: "/terraform/state-old",
    newPath: () => `/state`,
  },
];

const TERRAFORM_PATH_MAPPINGS = [
  {
    path: "/terraform-external/modules/aws",
    newPath: () =>
      "v2/modules?filter%5Bnamespace%5D=terraform-aws-modules&include=latest-version&page%5Bsize%5D=100&sort=-downloads",
  },
  {
    path: "/terraform-external/modules/aws/:module",
    newPath: (params: any) =>
      `v1/modules/terraform-aws-modules/${params.module}/aws`,
  },
];

const ORGANIZATION_PATH_MAPPINGS = [
  {
    path: "/orgs/:orgId",
    newPath: (params) => `/v1/orgs/${params.orgId}`,
  }
];

const USERS_PATH_MAPPINGS = [
  {
    path: "/users/add",
    newPath: (params: any) => `/v1/orgs/${params.orgId}/auth/users`,
  },
  {
    path: "/users/get",
    newPath: (params: any) => `/v1/orgs/${params.orgId}/auth/users`,
  },
  {
    path: "/users/delete/:username",
    newPath: (params) => `/v1/orgs/${params.orgId}/auth/users/${params.username}`,
  }
];

const EVENT_API = [
  {
    path: "/error-api",
    newPath: (params: any) => `/status?company=${params.orgName}`,
  },
  {
    path: "/events/stream",
    newPath: () => ``,
  },
];

const OPERATION_PATH_MAPPING = [
  {
    path: "/ops/is-provisioned",
    newPath: (params: any) => `/v1/orgs/${params.orgId}/ops/is-provisioned`,
  },
];

const API_PATH_MAPPINGS = [
  {
    path: "/api/teams",
    newPath: (params: any) => `v1/orgs/${params.orgId}/teams`,
  },
  {
    path: "/api/teams/:teamId/environments",
    newPath: (params: any) => `v1/orgs/${params.orgId}/teams/${params.teamId}/environments`,
  },
  {
    path: "/api/teams/:teamId/environments/:envId",
    newPath: (params: any) => `v1/orgs/${params.orgId}/teams/${params.teamId}/environments/${params.envId}`,
  },
  {
    path: "/api/teams/:teamId/environments/:envId/components",
    newPath: (params: any) => `v1/orgs/${params.orgId}/teams/${params.teamId}/environments/${params.envId}/components`,
  },
  {
    path: "/api/stream",
    newPath: (params: any) =>
      `v1/orgs/${params.orgId}/stream`,
  },
]

// eslint-disable-next-line no-unused-vars
function mapToRegex(mapping: {
  path: string;
  newPath: (params: any) => string;
}) {
  return {
    pathMatch: match(mapping.path, {
      encode: encodeURI,
      decode: decodeURIComponent,
    }),
    newPath: mapping.newPath,
  };
}

export const CD_MAPPINGS: PathMapping[] = CD_PATH_MAPPINGS.map(mapToRegex);
export const WF_MAPPINGS: PathMapping[] = WF_PATH_MAPPINGS.map(mapToRegex);
export const USERS_MAPPINGS: PathMapping[] = USERS_PATH_MAPPINGS.map(mapToRegex);
export const AUDIT_MAPPINGS: PathMapping[] = AUDIT_PATH_MAPPINGS.map(mapToRegex);
export const SECRET_MAPPINGS: PathMapping[] = SECRET_PATH_MAPPINGS.map(mapToRegex);
export const TERRAFORM_MAPPINGS: PathMapping[] = TERRAFORM_PATH_MAPPINGS.map(mapToRegex);
export const STATE_MAPPINGS: PathMapping[] = STATE_PATH_MAPPINGS.map(mapToRegex);
export const ORGANIZATION_MAPPINGS: PathMapping[] = ORGANIZATION_PATH_MAPPINGS.map(mapToRegex);
export const EVENT_MAPPINGS: PathMapping[] = EVENT_API.map(mapToRegex);
export const OPERATION_MAPPINGS: PathMapping[] = OPERATION_PATH_MAPPING.map(mapToRegex);
export const API_MAPPINGS: PathMapping[] = API_PATH_MAPPINGS.map(mapToRegex);

import { match, MatchFunction } from "path-to-regexp";

export interface PathMapping {
  pathMatch: MatchFunction<object>;
  // eslint-disable-next-line no-unused-vars
  newPath: (params: any) => string;
}

const CD_PATH_MAPPINGS = [
  {
    path: "/cd/api/v1/workspace",
    newPath: (params: any) => `/api/v1/projects/${params.team}`,
  },
  {
    path: "/cd/api/v1/projects",
    newPath: (params: any) => `/api/v1/applications?selector=type=project`,
  },
  {
    path: "/cd/api/v1/stream/projects/:resourceVersion",
    newPath: (params: any) =>
      `/api/v1/stream/applications?resourceVersion=${params.resourceVersion}`,
  },
  {
    path: "/cd/api/v1/projects/:projectId",
    newPath: (params: any) =>
      `/api/v1/applications/${params.projectId}?selector=type=project&project=${params.projectId}`,
  },
  {
    path: "/cd/api/v1/environments",
    newPath: (params: any) => `/api/v1/applications?selector=type=environment&project=${params.orgName}`,
  },
  {
    path: "/cd/api/v1/stream/environments/:environmentId",
    newPath: (params: any) =>
      `/api/v1/stream/applications?name=${params.environmentId}`,
  },
  {
    path: "/cd/api/v1/projects/:projectId/environments",
    newPath: (params: any) =>
      `/api/v1/applications?selector=type=environment,project_id=${params.projectId}&project=${params.projectId}`,
  },
  {
    path: "/cd/api/v1/stream/projects/:projectId/environments/:resourceVersion",
    newPath: (params: any) =>
      `/api/v1/stream/applications?resourceVersion=${params.resourceVersion}&selector=type=environment,project_id=${params.projectId}&project=${params.projectId}`,
  },
  {
    path: "/cd/api/v1/projects/:projectId/environments/:environmentId",
    newPath: (params: any) =>
      `/api/v1/applications/${params.environmentId}?selector=type=config,project_id=${params.projectId}&project=${params.projectId}`,
  },
  {
    path: "/cd/api/v1/projects/:projectId/environments/:environmentId/config",
    newPath: (params: any) =>
      `/api/v1/applications?selector=type=config,project_id=${params.projectId},environment_id=${params.environmentId}&project=${params.projectId}`,
  },
  {
    path: "/cd/api/v1/config",
    newPath: (params: any) => `/api/v1/applications?selector=type=config`,
  },
  {
    path: "/cd/api/v1/stream/projects/:projectId/environments/:environmentId/config/:resourceVersion",
    newPath: (params: any) =>
      `/api/v1/stream/applications?resourceVersion=${params.resourceVersion}&selector=type=config,project_id=${params.projectId},environment_id=${params.environmentId}&project=${params.projectId}`,
  },
  {
    path: "/cd/api/v1/stream/watcher/projects/:projectId",
    newPath: (params: any) =>
      `/api/v1/stream/applications?name=${params.projectId}-team-watcher`,
  },
  {
    path: "/cd/api/v1/projects/:environmentId/sync",
    newPath: (params: any) =>
      `/api/v1/applications/${params.environmentId}/sync`,
  },
  {
    path: "/cd/api/v1/projects/:environmentId/delete",
    newPath: (params: any) =>
      `/api/v1/applications/${params.environmentId}/resource?name=${params.environmentId}&namespace=${params.orgName}-executor&resourceName=${params.environmentId}&version=v1alpha1&kind=Workflow&group=argoproj.io&force=true&orphan=false`,
  },
  {
    path: "/cd/api/v1/projects/watcher/:projectId/sync",
    newPath: (params: any) =>
      `/api/v1/applications/${params.projectId}-team-watcher/sync`,
  },
  {
    path: "/cd/api/v1/component/:componentName",
    newPath: (params: any) => `api/v1/applications/${params.componentName}`,
  },
  {
    path: "/cd/api/v1/applications/:config/resource-tree",
    newPath: (params: any) =>
      `api/v1/applications/${params.config}/resource-tree`,
  },
  {
    path: "/cd/api/v1/stream/applications/:config/resource-tree",
    newPath: (params: any) =>
      `api/v1/stream/applications?name=${params.config}`,
  },
  {
    path: "/cd/api/v1/applications/:config/events",
    newPath: (params: any) => `api/v1/applications/${params.config}/events`,
  },
];

const WF_PATH_MAPPINGS = [
  {
    path: "/wf/api/v1/projects/:projectId/environments/:environmentId/config/:configId",
    newPath: (params: any) =>
      `/api/v1/workflows/${params.orgName}-executor?listOptions.labelSelector=config_id=${params.configId},project_id=${params.projectId},environment_id=${params.environmentId},team=${params.team}`,
  },
  {
    path: "/wf/api/v1/projects/:projectId/environments/:environmentId/config/:configId/:workflowId",
    newPath: (params: any) =>
      `/api/v1/workflows/${params.orgName}-executor/${params.workflowId}`,
  },
  {
    path: "/wf/api/v1/projects/:projectId/environments/:environmentId/config/:configId/:workflowId/approve",
    newPath: (params: any) =>
      `/api/v1/workflows/${params.orgName}-executor/${params.workflowId}/resume`,
  },
  {
    path: "/wf/api/v1/projects/:projectId/environments/:environmentId/config/:configId/:workflowId/decline",
    newPath: (params: any) =>
      `/api/v1/workflows/${params.orgName}-executor/${params.workflowId}/stop`,
  },
  {
    path: "/wf/api/v1/projects/:projectId/environments/:environmentId/config/:configId/:workflowId/log",
    newPath: (params: any) =>
      `/api/v1/workflows/${params.orgName}-executor/${params.workflowId}/log?logOptions.container=main&logOptions.follow=false`,
  },
  {
    path: "/wf/api/v1/projects/:projectId/environments/:environmentId/config/:configId/:workflowId/log/:nodeId",
    newPath: (params: any) =>
      `/api/v1/workflows/${params.orgName}-executor/${params.workflowId}/${params.nodeId}/log?logOptions.container=main&logOptions.follow=false`,
  },
  {
    path: "/wf/api/v1/stream/projects/:projectId/environments/:environmentId/config/:configId/:workflowId/log",
    newPath: (params: any) =>
      `/api/v1/workflows/${params.orgName}-executor/${params.workflowId}/log?logOptions.container=main&logOptions.follow=true`,
  },
  {
    path: "/wf/api/v1/stream/projects/:projectId/environments/:environmentId/config/:configId/:workflowId/log/:nodeId",
    newPath: (params: any) =>
      `/api/v1/workflows/${params.orgName}-executor/${params.workflowId}/log?logOptions.container=main&logOptions.follow=true&podName=${params.nodeId}`,
  },
  {
    path: "/wf/api/v1/stream/projects/:projectId/environments/:environmentId/config/:configId/:workflowId",
    newPath: (params: any) =>
      `/api/v1/workflow-events/${params.orgName}-executor?listOptions.fieldSelector=metadata.name=${params.workflowId}`,
  },
];

const COSTING_PATH_MAPPINGS = [
  {
    path: "/costing/api/v1/all",
    newPath: (params) => `/v1/orgs/${params.orgId}/costing/all`,
  },
  {
    path: "/costing/api/v1/team/:name",
    newPath: (params: any) => `/v1/orgs/${params.orgId}/costing/team/${params.name}`,
  },
  {
    path: "/costing/api/v1/environment/:teamName/:environmentName",
    newPath: (params: any) =>
      `/v1/orgs/${params.orgId}/costing/environment/${params.teamName}/${params.environmentName}`,
  },
  {
    path: "/costing/api/v1/component/:componentId",
    newPath: (params: any) => `/v1/orgs/${params.orgId}/costing/component/${params.componentId}`,
  },
  {
    path: "/costing/api/v1/resources/:componentId",
    newPath: (params: any) => `/v1/orgs/${params.orgId}/costing/resources/${params.componentId}`,
  },
  {
    path: "/costing/stream/api/v1/all",
    newPath: (params: any) => `/v1/orgs/${params.orgId}/costing/stream/all`,
  },
  {
    path: "/costing/stream/api/v1/notify",
    newPath: (params: any) => `/v1/orgs/${params.orgId}/costing/stream/notify`,
  },
  {
    path: "/costing/stream/api/v1/team/:name",
    newPath: (params: any) => `/v1/orgs/${params.orgId}/costing/stream/team/${params.name}`,
  },
  {
    path: "/costing/stream/api/v1/environment/:teamName/:environmentName",
    newPath: (params: any) =>
      `/v1/orgs/${params.orgId}/costing/stream/environment/${params.teamName}/${params.environmentName}`,
  },
];
// /v1/orgs/1/reconciliation/environments/equity-error-2
// /reconciliation/api/v1/environments/equity-error-2
const AUDIT_PATH_MAPPINGS = [
  {
    path: "/reconciliation/api/v1/environments/:id",
    newPath: (params: any) =>
      `/v1/orgs/${params.orgId}/reconciliation/environments/${params.id}`,
  },
  {
    path: "/reconciliation/api/v1/components/:id",
    newPath: (params: any) =>
      `/v1/orgs/${params.orgId}/reconciliation/components/${params.id}`,
  },
  {
    path: "/reconciliation/api/v1/component/:id",
    newPath: (params: any) =>
      `/v1/orgs/${params.orgId}/reconciliation/component/${params.id}`,
  },
  {
    path: "/reconciliation/api/v1/environment/:id",
    newPath: (params: any) =>
      `/v1/orgs/${params.orgId}/reconciliation/component/${params.id}`,
  },
  {
    path: "/reconciliation/api/v1​/notification/save",
    newPath: (params: any) => `/v1/orgs/${params.orgId}/reconciliation/notification/save`,
  },
  {
    path: "/reconciliation/api/v1/notifications/get/:teamName",
    newPath: (params: any) =>
      `/v1/orgs/${params.orgId}/reconciliation/notifications/get/${params.teamName}`,
  },
  {
    path: "/reconciliation/api/v1/notification/seen/:id",
    newPath: (params: any) =>
      `/v1/orgs/${params.orgId}/reconciliation/notification/seen/${params.id}`,
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
    path: "/reconciliation/api/v1/notifications/:teamName",
    newPath: (params: any) =>
      `/v1/orgs/${params.orgId}/reconciliation/notifications/${params.teamName}`,
  },
  {
    path: "/reconciliation/api/v1/components/notify/:id",
    newPath: (params: any) =>
      `/v1/orgs/${params.orgId}/reconciliation/components/notify/${params.id}`,
  },
  {
    path: "/reconciliation/api/v1/environments/notify/:id",
    newPath: (params: any) =>
      `/v1/orgs/${params.orgId}/reconciliation/environments/notify/${params.id}`,
  },
  {
    path: "/reconciliation/api/v1/approved-by/:id",
    newPath: (params: any) =>
      `/v1/orgs/${params.orgId}/reconciliation/approved-by/${params.id}/${params.email}`,
  },
  {
    path: "/reconciliation/api/v1/approved-by/:id/:rid",
    newPath: (params: any) =>
      `/v1/orgs/${params.orgId}/reconciliation/approved-by/${params.id}/${params.rid}`,
  },
];

const SECRET_PATH_MAPPINGS = [
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
export const COSTING_MAPPINGS: PathMapping[] = COSTING_PATH_MAPPINGS.map(mapToRegex);
export const AUDIT_MAPPINGS: PathMapping[] = AUDIT_PATH_MAPPINGS.map(mapToRegex);
export const SECRET_MAPPINGS: PathMapping[] = SECRET_PATH_MAPPINGS.map(mapToRegex);
export const TERRAFORM_MAPPINGS: PathMapping[] = TERRAFORM_PATH_MAPPINGS.map(mapToRegex);
export const STATE_MAPPINGS: PathMapping[] = STATE_PATH_MAPPINGS.map(mapToRegex);
export const ORGANIZATION_MAPPINGS: PathMapping[] = ORGANIZATION_PATH_MAPPINGS.map(mapToRegex);
export const EVENT_MAPPINGS: PathMapping[] = EVENT_API.map(mapToRegex);
export const OPERATION_MAPPINGS: PathMapping[] = OPERATION_PATH_MAPPING.map(mapToRegex);

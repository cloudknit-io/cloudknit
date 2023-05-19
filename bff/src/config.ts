const config = {
  SELECTED_ORG_HEADER: "CLOUDKNIT-SELECTED-ORG",
  AUTH0_WEB_BASE_URL: process.env.AUTH0_WEB_BASE_URL,
  AUTH0_WEB_CLIENT_ID: process.env.AUTH0_WEB_CLIENT_ID,
  AUTH0_WEB_SECRET: process.env.AUTH0_WEB_SECRET,
  AUTH0_API_CLIENT_ID: process.env.AUTH0_API_CLIENT_ID,
  AUTH0_API_SECRET: process.env.AUTH0_API_SECRET,
  AUTH0_API_AUDIENCE: process.env.AUTH0_API_AUDIENCE,
  AUTH0_ISSUER_BASE_URL: process.env.AUTH0_ISSUER_BASE_URL,
  WEB_URL: process.env.SITE_URL,
  API_URL: `${process.env.ZLIFECYCLE_API_URL}/v1`,
  ARGOCD_URL: process.env.ARGO_CD_API_URL,
  PLAYGROUND_APP: false, //process.env.CK_PLAYGROUND == 'true',
  argoWFUrl: (orgName: string) =>
    process.env.ARGO_WORKFLOW_API_URL.replaceAll(":org", orgName),
  stateMgrUrl: (orgName: string) =>
    process.env.ZLIFECYCLE_STATE_MANAGER_URL.replaceAll(":org", orgName),
  isProd: (): boolean => process.env.NODE_ENV === "PRODUCTION",
  isDebug: (): boolean => process.env.LOG_LEVEL === "debug",
};
export default config;

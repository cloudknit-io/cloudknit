const config = {
    SELECTED_ORG_HEADER: 'CLOUDKNIT-SELECTED-ORG',
    AUTH0_BASE_URL: process.env.AUTH0_BASE_URL,
    AUTH0_CLIENT_ID: process.env.AUTH0_CLIENT_ID,
    AUTH0_ISSUER_BASE_URL: process.env.AUTH0_ISSUER_BASE_URL,
    AUTH0_SECRET: process.env.AUTH0_SECRET,
    WEB_URL: process.env.SITE_URL,
    ARGOCD_URL: process.env.ARGO_CD_API_URL,
    argoWFUrl: (orgName: string) => process.env.ARGO_WORKFLOW_API_URL.replaceAll(':org', orgName),
    stateMgrUrl: (orgName: string) => process.env.ZLIFECYCLE_STATE_MANAGER_URL.replaceAll(':org', orgName),
    isProd: () : boolean => process.env.NODE_ENV === 'PRODUCTION',
    isDebug: () : boolean => process.env.LOG_LEVEL === 'debug'
}

export default config;

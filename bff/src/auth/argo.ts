import axios from 'axios';
import config from '../config';
import logger from '../utils/logger';

let sessions = {};
let lock = false;

async function argoCdLogin(org: string, username: string, password: string) {
  const url = `${config.argoCDUrl(org)}/api/v1/session`;

  try {
    const resp = await axios.post(url, {
      username,
      password
    });
  
    const { token } = resp.data;
  
    return token;
  } catch (err) {
    logger.error('could not login to argocd', { org, error: err.message });
  }

  return null;
}

function wait(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

async function getArgoCDPassword(org: string): Promise<string> {
  try {
    const url = `${process.env.ZLIFECYCLE_API_URL}/v1/orgs/${org}/secrets/get/ssm-secret`;
    const resp = await axios.post(url, {
      path: "/argocd/zlapi/password"
    });
  
    const { token } = resp.data;
  
    return token;
  } catch (err) {
    logger.error('could not login to argocd', { org, error: err.message });
  }
}

async function createSession(orgName: string) {  
  const session = { 
    token: await argoCdLogin(orgName, 'zlapi', await getArgoCDPassword(orgName)), 
    ttl: Date.now() + 10800000 // 3 hours
  };

  sessions[orgName] = session;

  return session;
}

function isExpired(orgName: string) {
  const session = sessions[orgName];
  const now = Date.now();

  return !session || !session.ttl || now > session.ttl;
}

async function getArgoToken (orgName: string): Promise<string> {
  while (lock === true) {
    await wait(100);
  }

  lock = true;

  if (isExpired(orgName)) {
    logger.info('Refreshed ArgoCD Token', {org: orgName});
    await createSession(orgName);
  }

  lock = false;

  return sessions[orgName].token;
}

export async function getArgoCDAuthHeader (orgName: string): Promise<any> {
  return {
    authorization: `Bearer ${await getArgoToken(orgName)}`,
  };
}

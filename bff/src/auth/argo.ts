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

    if (!resp.data || !resp.data.value) {
      throw new Error('password was not return from api');
    }
    
    const { value } = resp.data;
  
    return value;
  } catch (err) {
    logger.error('could not retrieve argocd password from api', { org, error: err.message });
    throw err;
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

  try {
    if (isExpired(orgName)) {
      await createSession(orgName);
      logger.info('Refreshed ArgoCD Token', {org: orgName});
    }
    
    return sessions[orgName].token;
  } catch (err) {
    logger.error('failed to create ArgoCD session', {org: orgName, err: err.message});
  } finally {
    lock = false;
  }

  return null;
}

export async function getArgoCDAuthHeader (orgName: string): Promise<any> {
  return {
    authorization: `Bearer ${await getArgoToken(orgName)}`,
  };
}

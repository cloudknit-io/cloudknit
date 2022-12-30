import config from "../config";
import * as express from 'express';
import { Organization, User } from "../models/user.interface";
import { BFFRequest } from "../types";
import { getUser } from "../auth/auth";
import axios from "axios";
import logger from '../utils/logger';
import { getArgoCDAuthHeader } from "../auth/argo";

const orgFromReq = async (req: BFFRequest) : Promise<Organization> => {
  if (!req.cookies[config.SELECTED_ORG_HEADER]) {
      return;
  }

  const orgName = req.cookies[config.SELECTED_ORG_HEADER];
  let org = req.appSession.organizations.find((org) => org.name === orgName);

  if (org) {
    return org;
  }

  // org was not in appSession
  // check database to see if user is part of org
  const user = userFromReq(req);
  const dbUser = await getUser(user.username);

  org = dbUser.organizations.find((org) => org.name === orgName)

  if (!org) {
    return null;
  }

  return org;
};

const userFromReq = (req: BFFRequest) : User => {
  if (!req.appSession) {
      return;
  }

  return req.appSession.user
};

const handleNoOrg = (res: express.Response) => {
  res.status(401).send({error: 'no organization has been selected'});
}

const getSystemSSMSecret = async (orgName: string, secretPath: string) : Promise<string> => {
  try {
    const url = `${process.env.ZLIFECYCLE_API_URL}/v1/system/ssmsecret`;
    const resp = await axios.get(url, {
      params: {
        path: secretPath
      }
    });
    
    const { value } = resp.data;
  
    return value;
  } catch (err) {
    logger.error('could not retrieve SSM value from api', { org: orgName, error: err.message });
    return null;
  }
};

const getSSMSecret = async (orgName: string, secretPath: string) : Promise<string> => {
  try {
    const url = `${process.env.ZLIFECYCLE_API_URL}/v1/orgs/${orgName}/secrets/get/ssm-secret`;
    const resp = await axios.post(url, {
      path: secretPath
    });
    
    const { value } = resp.data;
  
    return value;
  } catch (err) {
    logger.error('could not retrieve SSM value from api', { org: orgName, error: err.message });
    return null;
  }
};

const getOrg = async (orgName: string) : Promise<Organization> => {
  try {
    const url = `${process.env.ZLIFECYCLE_API_URL}/v1/orgs/${orgName}`;
    const resp = await axios.get(url);
    
    return resp.data;
  } catch (err) {
    logger.error('could not retrieve organization from api', { org: orgName, error: err.message });
    return null;
  }
}
const syncWatcher = async (orgName: string, teamName: string) => {
  try {
  const {authorization} = await getArgoCDAuthHeader(orgName);

  const requestPayload = {
    appNamespace: "zlifecycle-system",
    revision: "HEAD",
    prune: false,
    dryRun: false,
    strategy: {
      hook: {
        force: false,
      },
    },
    resources: null,
    syncOptions: {
      items: ["ServerSideApply=true"],
    },
    retryStrategy: {
      limit: 1,
      backoff: {
        duration: "5s",
        maxDuration: "3m0s",
        factor: 2,
      },
    },
  };
  await axios.post(
    `${config.ARGOCD_URL}/api/v1/applications/${orgName}-${teamName}-team-watcher/sync`,
    requestPayload,
    {
      headers: {
        authorization,
      },
    }
  );
  } catch (err) {
    logger.error('could not sync watcher', { orgName, teamName, err });
  }
};

export default {
  orgFromReq,
  userFromReq,
  handleNoOrg,
  getSystemSSMSecret,
  getSSMSecret,
  getOrg,
  syncWatcher,
};
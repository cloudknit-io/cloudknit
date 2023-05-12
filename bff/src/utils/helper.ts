import config from "../config";
import * as express from "express";
import { Organization, User } from "../models/user.interface";
import { AppSession, BFFRequest } from "../types";
import { getUser } from "../auth/auth";
import axios from "axios";
import logger from "../utils/logger";
import { getArgoCDAuthHeader } from "../auth/argo";
import ckConfig from "../config";

const orgFromReq = async (req: BFFRequest): Promise<Organization> => {
  if (!req.cookies[config.SELECTED_ORG_HEADER]) {
    return;
  }

  const orgName = req.cookies[config.SELECTED_ORG_HEADER];
  const session = appSession(req);

  if (!session) {
    return null;
  }

  let org = appSession(req).organizations.find((org) => org.name === orgName);

  if (org) {
    return org;
  }

  // org was not in appSession
  // check database to see if user is part of org
  const user = userFromReq(req);
  const dbUser = await getUser(user.username);

  org = dbUser.organizations.find((org) => org.name === orgName);

  if (!org) {
    return null;
  }

  return org;
};

const userFromReq = (req: BFFRequest): User => {
  if (!appSession(req)) {
    return;
  }

  return appSession(req).user;
};

const appSessionFromReq = (req: BFFRequest): AppSession => {
  if (!req.session) {
    return;
  }

  const appSession: AppSession = {
    access_token: req.session.passport.user.tokens.access_token,
    id_token: req.session.passport.user.tokens.id_token,
    token_type: req.session.passport.user.tokens.token_type,
    expires_at: req.session.passport.user.tokens.expires_at.toString(),
    organizations: [],
    user: null,
    refresh_token: "",
  };

  return appSession;
};

const handleNoOrg = (res: express.Response) => {
  res.status(401).send({ error: "no organization has been selected" });
};

const getSystemSSMSecret = async (
  orgName: string,
  secretPath: string
): Promise<string> => {
  try {
    const url = `${process.env.ZLIFECYCLE_API_URL}/v1/system/ssmsecret`;
    const resp = await axios.get(url, {
      params: {
        path: secretPath,
      },
    });

    const { value } = resp.data;

    return value;
  } catch (err) {
    logger.error("could not retrieve SSM value from api", {
      org: orgName,
      error: err.message,
    });
    return null;
  }
};

const getSSMSecret = async (
  orgName: string,
  secretPath: string
): Promise<string> => {
  try {
    const url = `${process.env.ZLIFECYCLE_API_URL}/v1/orgs/${orgName}/secrets/get/ssm-secret`;
    const resp = await axios.post(url, {
      path: secretPath,
    });

    const { value } = resp.data;

    return value;
  } catch (err) {
    logger.error("could not retrieve SSM value from api", {
      org: orgName,
      error: err.message,
    });
    return null;
  }
};

const getOrg = async (orgName: string): Promise<Organization> => {
  try {
    const url = `${process.env.ZLIFECYCLE_API_URL}/v1/orgs/${orgName}`;
    const resp = await axios.get(url);

    return resp.data;
  } catch (err) {
    logger.error("could not retrieve organization from api", {
      org: orgName,
      error: err.message,
    });
    return null;
  }
};

const syncWatcher = async (orgName: string, teamName: string) => {
  try {
    const { authorization } = await getArgoCDAuthHeader(orgName);
    await axios.post(
      `${config.ARGOCD_URL}/api/v1/applications/${orgName}-${teamName}-team-watcher/sync`,
      {},
      {
        headers: {
          authorization,
        },
      }
    );
  } catch (err) {
    logger.error("could not sync watcher", { orgName, teamName, err });
  }
};

const isOktaAuth = () =>
  ckConfig.AUTH0_ISSUER_BASE_URL.includes("oktapreview.com") ||
  ckConfig.AUTH0_ISSUER_BASE_URL.includes("okta.com");

const isGuestAuth = () => true; //ckConfig.PLAYGROUND_APP === "true";

export const appSession = (req: BFFRequest): any => {
  if (isOktaAuth() || isGuestAuth()) {
    return req.session.appSession;
  }
  return req.appSession;
};

export const oidcUser = (req: BFFRequest) => {
  if (isOktaAuth()) {
    return req.session.passport.user;
  }
  if (isGuestAuth()) {
    return req.session.appSession.user;
  }
  return req.oidc.user;
};

export default {
  orgFromReq,
  userFromReq,
  handleNoOrg,
  getSystemSSMSecret,
  getSSMSecret,
  getOrg,
  syncWatcher,
  appSessionFromReq,
  isOktaAuth,
  isGuestAuth
};

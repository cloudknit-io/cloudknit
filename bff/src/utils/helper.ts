import axios from "axios";
import config from "../config";
import * as express from 'express';
import { Organization, User } from "../models/user.interface";
import { BFFRequest } from "../types";
import { getUser } from "../auth/auth";

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
  res.status(401).send();
}

export default {
  orgFromReq,
  userFromReq,
  handleNoOrg
}

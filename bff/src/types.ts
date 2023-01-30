import { Session } from "express-openid-connect";
import { Organization, User } from "./models/user.interface";
import * as express from "express";

type AppSession = Session & {
  user: User;
  organizations: Array<Organization>;
};

type ApiAuth = {
  ckOrgId: number;
  iss: string;
  sub: string;
  aud: string;
  iat: number;
  exp: number;
  azp: string;
  gty: string;
  permissions: string[];
};

export type BFFRequest = express.Request & { appSession: AppSession };

export type ExternalAPIRequest = express.Request & { auth: ApiAuth };

import { Session } from "express-openid-connect";
import { Organization, User } from "./models/user.interface";
import * as express from "express";

export type AppSession = Session & {
  user: User;
  organizations: Array<Organization>;
};

export interface ExpressSession {
  cookie: Cookie;
  passport: Passport;
  appSession: AppSession;
}
export interface Cookie {
  originalMaxAge?: null;
  expires?: null;
  httpOnly: boolean;
  path: string;
}
export interface Passport {
  user: SessionUser;
}
export interface SessionUser {
  userinfo: Userinfo;
  tokens: Tokens;
}
export interface Userinfo {
  sub: string;
  name: string;
  locale: string;
  email: string;
  preferred_username: string;
  given_name: string;
  family_name: string;
  zoneinfo: string;
  email_verified: boolean;
}
export interface Tokens {
  token_type: string;
  expires_at: number;
  access_token: string;
  scope: string;
  id_token: string;
}


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

export type BFFRequest = express.Request & { appSession: AppSession } & { session: ExpressSession };

export type ExternalAPIRequest = express.Request & { auth: ApiAuth };

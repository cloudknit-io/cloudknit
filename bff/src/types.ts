import { Session } from "express-openid-connect";
import { Organization, User } from "./models/user.interface";
import * as express from 'express';

type AppSession = Session & {
  user: User,
  organizations: Array<Organization>
}

export type BFFRequest = express.Request & { appSession: AppSession };

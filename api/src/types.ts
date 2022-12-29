import { Request } from "express";
import { Environment, Team } from "./typeorm";
import { Organization } from "./typeorm/Organization.entity";

export type APIRequest = Request & { org: Organization, team: Team, env?: Environment }

export const SqlErrorCodes = {
  NO_DEFAULT: 'ER_NO_DEFAULT_FOR_FIELD',
  DUP_ENTRY: 'ER_DUP_ENTRY'
}

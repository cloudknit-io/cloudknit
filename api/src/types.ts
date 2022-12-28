import { Request } from "express";
import { Team } from "./typeorm";
import { Organization } from "./typeorm/Organization.entity";

export type APIRequest = Request & { org: Organization, team: Team }

export const SqlErrorCodes = {
  NO_DEFAULT: 'ER_NO_DEFAULT_FOR_FIELD'
}

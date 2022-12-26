import { Request } from "express";
import { Organization } from "./typeorm/Organization.entity";

export type APIRequest = Request & { org: Organization }

export const SqlErrorCodes = {
  NO_DEFAULT: 'ER_NO_DEFAULT_FOR_FIELD'
}

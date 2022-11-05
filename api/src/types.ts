import { Request } from "express";
import { Organization } from "./typeorm/Organization.entity";

export type APIRequest = Request & { org: Organization }

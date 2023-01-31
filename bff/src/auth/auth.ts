import axios from "axios";
import { OpenidRequest, OpenidResponse, Session } from "express-openid-connect";
import { User } from "../models/user.interface";
import logger from "../utils/logger";
import * as jose from "jose";
import ckConfig from "../config";
import * as createError from "http-errors";
import { BFFRequest } from "../types";
import helper from "../utils/helper";
import * as jwt from "express-jwt";
import * as jwks from "jwks-rsa";

export async function getUser(username: string): Promise<User> {
  try {
    const user = await axios.get(
      `${process.env.ZLIFECYCLE_API_URL}/v1/users/${username}`
    );

    return user.data;
  } catch (err) {
    if (axios.isAxiosError(err)) {
      // @ts-ignore
      logger.error("get user error", { error: err.toJSON().message });
    } else {
      logger.error("get user error", { error: { message: err.message } });
    }

    return null;
  }
}

async function createUser(
  username: string,
  email: string,
  role: string,
  name: string
): Promise<User> {
  try {
    const user = await axios.post(
      `${process.env.ZLIFECYCLE_API_URL}/v1/users/`,
      {
        username,
        email,
        role,
        name,
      }
    );

    return user.data;
  } catch (err) {
    if (axios.isAxiosError(err)) {
      // @ts-ignore
      logger.error("create user error", { error: err.toJSON().message });
    } else {
      logger.error("create user error", { error: { message: err.message } });
    }

    return null;
  }
}

export const organizationMW = (req: BFFRequest, res, next) => {
  if (!helper.orgFromReq(req)) {
    logger.error("no org selected");
    helper.handleNoOrg(res);
    return;
  }

  next();
};

export function getAuth0Config() {
  return {
    authRequired: true,
    auth0Logout: true,
    authorizationParams: {
      response_type: "id_token",
      response_mode: "form_post",
      scope: "openid profile email",
      connection: "github",
    },
    errorOnRequiredAuth: false,
    secret: ckConfig.AUTH0_WEB_SECRET,
    baseURL: ckConfig.AUTH0_WEB_BASE_URL,
    clientID: ckConfig.AUTH0_WEB_CLIENT_ID,
    issuerBaseURL: ckConfig.AUTH0_ISSUER_BASE_URL,
    routes: {
      login: "/auth/login",
      logout: "/auth/logout",
      callback: "/auth/callback",
    },
    // https://auth0.github.io/express-openid-connect/interfaces/configparams.html#aftercallback
    afterCallback: async (
      req: OpenidRequest,
      res: OpenidResponse,
      session: Session
    ) => {
      const claims = jose.decodeJwt(session.id_token);

      if (!claims.email_verified) {
        logger.error(`email is not verified for user ${claims.email}`, {
          data: claims,
        });
        throw new createError.Unauthorized();
      }

      // @ts-ignore
      let user = await getUser(claims.nickname);

      if (!user) {
        try {
          user = await createUser(
            // @ts-ignore
            claims.nickname,
            claims.email,
            "Admin",
            claims.name
          );

          if (!user) {
            throw new createError.Unauthorized();
          }

          logger.info(`create user ${user.username}`, { user });

          return {
            ...session,
            user,
            organizations: [],
          };
        } catch (err) {
          logger.error(`could not create user ${claims.nickname}`, {
            error: err.message,
          });
          throw new createError.Unauthorized();
        }
      }

      if (!user.organizations || user.organizations.length == 0) {
        logger.info(`user ${claims.nickname} does not have any organizations`);

        return {
          ...session,
          user,
          organizations: [],
        };
      }

      return {
        ...session,
        user,
        organizations: user.organizations,
      };
    },
  };
}

export function apiAuthMw() {
  return jwt.expressjwt({
    // @ts-ignore
    secret: jwks.expressJwtSecret({
      cache: true,
      rateLimit: true,
      jwksRequestsPerMinute: 5,
      jwksUri: new URL("/.well-known/jwks.json", ckConfig.AUTH0_ISSUER_BASE_URL)
        .href,
    }),
    audience: ckConfig.AUTH0_API_AUDIENCE,
    issuer: new URL("/", ckConfig.AUTH0_ISSUER_BASE_URL).href,
    algorithms: ["RS256"],
  });
}

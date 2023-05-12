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
import * as express from "express";
import { auth, requiresAuth } from "express-openid-connect";
import * as session from "express-session";
import { ExpressOIDC } from "@okta/oidc-middleware";
import * as crypto from "crypto";

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

export async function getPlaygroundUser(username: string): Promise<User> {
  try {
    const user = await axios.get(
      `${process.env.ZLIFECYCLE_API_URL}/v1/orgs/1/auth/users/${username}`
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

export async function createPlaygroundUser(ipv4: string): Promise<User> {
  const url = `${process.env.ZLIFECYCLE_API_URL}/v1/orgs/1/auth/playground/users`;
  console.log(url);
  try {
    const user = await axios.post(url, {
      ipv4,
    });

    console.log(user);

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

function getOktaAuthMW() {
  const oidc = new ExpressOIDC({
    issuer: ckConfig.AUTH0_ISSUER_BASE_URL,
    client_id: ckConfig.AUTH0_WEB_CLIENT_ID,
    client_secret: ckConfig.AUTH0_WEB_SECRET,
    appBaseUrl: ckConfig.AUTH0_WEB_BASE_URL,
    scope: "openid profile email",
    routes: {
      login: {
        // handled by this module
        path: "/auth/login",
      },
      loginCallback: {
        // handled by this module
        path: "/auth/callback",
        handler: async (req: BFFRequest, res: any, next: any) => {
          const session = req.session;
          const userInfo = session.passport.user.userinfo;

          if (!userInfo.email_verified) {
            logger.error(`email is not verified for user ${userInfo.email}`, {
              data: userInfo,
            });
            throw new createError.Unauthorized();
          }

          // @ts-ignore
          let user = await getUser(userInfo.preferred_username);
          if (!user) {
            try {
              user = await createUser(
                // @ts-ignore
                userInfo.preferred_username,
                userInfo.email,
                "Admin",
                userInfo.name
              );

              if (!user) {
                throw new createError.Unauthorized();
              }

              logger.info(`create user ${user.username}`, { user });
            } catch (err) {
              logger.error(
                `could not create user ${userInfo.preferred_username}`,
                {
                  error: err.message,
                }
              );
              throw new createError.Unauthorized();
            }
          }

          if (!user.organizations || user.organizations.length == 0) {
            logger.info(
              `user ${userInfo.preferred_username} does not have any organizations`
            );
          }

          req.session.appSession = {
            ...helper.appSessionFromReq(req),
            user,
            organizations: user.organizations || [],
          };
          next();
        },
        // handled by your application
        afterCallback: "/",
      },
    },
  });

  return oidc.router;
}

export async function guestAuthMW(req, res, next) {
  const ipv4 = getClientIP(req);
  logger.info("GUEST AUTH MW", {
    ipv4,
  });
  if (!ipv4) {
    res.status(500).send();
    return;
  }
  logger.info("Getting user for: ", {
    ipv4,
  });
  let user = await getPlaygroundUser(ipv4);
  if (!user) {
    logger.info("User not found for: ", {
      ipv4,
    });
    logger.info("Creating user for", {
      ipv4,
    });
    user = await createPlaygroundUser(ipv4);
  }
  if (user) {
    // Setting the appsession
    req.session.appSession = {
      user,
      organizations: user.organizations,
    };
  }
  logger.info("Current user info", {
    user,
  });
  next();
}

export function setUpAuth(app: express.Express, authRouter: express.Router) {
  if (helper.isGuestAuth()) {
    const MemoryStore = require("memorystore")(session);
    app.use(
      session({
        secret: crypto.randomUUID(),
        resave: false,
        saveUninitialized: false,
        cookie: { maxAge: 86400000 },
        store: new MemoryStore({
          checkPeriod: 86400000,
        }),
      })
    );
    app.use(guestAuthMW);
  } else if (helper.isOktaAuth()) {
    const MemoryStore = require("memorystore")(session);
    app.use(
      session({
        secret: crypto.randomUUID(),
        resave: false,
        saveUninitialized: false,
        cookie: { maxAge: 86400000 },
        store: new MemoryStore({
          checkPeriod: 86400000,
        }),
      })
    );

    app.use(getOktaAuthMW());

    authRouter.get("/auth/logout", (req: any, res: any, next: any) => {
      req.logout(function (err) {
        if (err) {
          return next(err);
        }
        res.redirect("/");
      });
    });
  } else {
    // auth0 router attaches /auth/login, /logout, and /callback routes to the baseURL
    app.use(auth(getAuth0Config()));

    authRouter.use(requiresAuth());
  }
}

function getClientIP(req) {
  if (!req.headers["x-forwarded-for"]) {
    return null;
  }

  // x-forwarded-for header returns the list of ips that our request has been forwarded by,
  // first being clients and then it depends on the no. of proxies that have forwarded it.
  const addresses = req.headers["x-forwarded-for"];

  // Getting the first ip since that is where the request has originated from.
  return addresses.split(",")[0];
}

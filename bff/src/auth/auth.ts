import axios from "axios";
import { OpenidRequest, OpenidResponse, Session } from "express-openid-connect";
import { User } from "../models/user.interface";
import logger from "../utils/logger";
import * as jose from 'jose';
import zlConfig from '../config';
import * as createError from 'http-errors';
import { BFFRequest } from "../types";
import helper from '../utils/helper';

export async function getUser(username: string): Promise<User> {
  try {
    const user = await axios.get(
      `${process.env.ZLIFECYCLE_API_URL}/v1/users/${username}`
    );

    return user.data;
  } catch (err) {
    if (axios.isAxiosError(err)) {
      // @ts-ignore
      logger.error('get user error', { error: err.toJSON().message });
    } else {
      logger.error('get user error', { error: { message: err.message } });
    }
    
    return null
  }
}

async function createUser(username: string, email: string, role: string, name: string): Promise<User> {
  try {
    const user = await axios.post(
      `${process.env.ZLIFECYCLE_API_URL}/v1/users/`,
      {
        username,
        email,
        role,
        name
      }
    );

    return user.data;
  } catch (err) {
    if (axios.isAxiosError(err)) {
      // @ts-ignore
      logger.error('create user error', { error: err.toJSON().message });
    } else {
      logger.error('create user error', { error: { message: err.message } });
    }

    return null
  }
}

export const organizationMW = (req: BFFRequest, res, next) => {
  if (!helper.orgFromReq(req)) {
    logger.error('no org selected')
    res.status(401).send({error: 'no organization has been selected'});
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
      connection: "github"
    },
    errorOnRequiredAuth: false,
    secret: zlConfig.AUTH0_SECRET,
    baseURL: zlConfig.AUTH0_BASE_URL,
    clientID: zlConfig.AUTH0_CLIENT_ID,
    issuerBaseURL: zlConfig.AUTH0_ISSUER_BASE_URL,
    routes: {
      login: '/auth/login',
      logout: '/auth/logout',
      callback: '/auth/callback'
    },
    // https://auth0.github.io/express-openid-connect/interfaces/configparams.html#aftercallback
    afterCallback: async (req: OpenidRequest, res: OpenidResponse, session: Session) => {
      const claims = jose.decodeJwt(session.id_token);

      if (!claims.email_verified) {
        logger.error(`email is not verified for user ${claims.email}`, { data: claims });
        throw new createError.Unauthorized();
      }
  
      // @ts-ignore
      let user = await getUser(claims.nickname);
  
      if (!user) {
        try {
          // @ts-ignore
          user = await createUser(claims.nickname, claims.email, 'Admin', claims.name);

          if (!user) {
            throw new createError.Unauthorized();
          }

          logger.info(`create user ${user.username}`, { user });

          return {
            ...session,
            user,
            organizations: []
          };
        } catch (err) {
          logger.error(`could not create user ${claims.nickname}`, { error: err.message })
          throw new createError.Unauthorized();
        }
      }
  
      if (!user.organizations || user.organizations.length == 0) {
        logger.info(`user ${claims.nickname} does not have any organizations`);
        
        return {
          ...session,
          user,
          organizations: []
        };
      }
      
      return {
        ...session,
        user,
        organizations: user.organizations
      };
    }
  };
}

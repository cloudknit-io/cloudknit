import { createLogger, format, transports } from 'winston';
import config from '../config';
import helper from './helper';
import { BFFRequest } from '../types';

const {
  combine, timestamp, prettyPrint, json
} = format;

const logger = createLogger({
  level: config.isDebug() ? 'debug' : 'info',
  format: combine(
    json(),
    timestamp(),
    config.isProd() ? null : prettyPrint({ colorize: true }),
  ),
  transports: [
    new transports.Console(),
  ]
});

const reqLogger = createLogger({
  level: 'info',
  format: combine(
    json(),
    timestamp()
  ),
  transports: [
    new transports.Console(),
  ]
});

export default logger;

export async function AuthRequestLogger(req: BFFRequest, res, next) {
  const reqUser = helper.userFromReq(req);
  const reqOrg = await helper.orgFromReq(req);

  let user, org = {};

  if (reqUser) {
    user = {
      username: reqUser.username,
      id: reqUser.id
    };
  }

  if (reqOrg) {
    org = {
      name: reqOrg.name,
      id: reqOrg.id
    };
  }

  reqLogger.info('authorized req log', { user, org, method: req.method, path: req.path });

  next();
}

export function ErrorLogger(err: Error, req, res, next) {
  if (err.message.startsWith('Authentication is required for this route')) {
    res.status(403);
    res.send();
    return;
  }
  
  logger.error({ type: "DefaultErrorHandler", message: err.message, stack: err.stack });

  res.status(500);
  res.json({ error: "internal server error" });
};

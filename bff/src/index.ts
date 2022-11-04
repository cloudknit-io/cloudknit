import './configuration';

import * as express from 'express';
import * as correlator from 'express-correlation-id';
import * as cors from 'cors';
import * as cookieParser from 'cookie-parser';
import { AuthRequestLogger, ErrorLogger } from './utils/logger';
import logger from './utils/logger';
import { getUser } from './auth/auth';
import zlConfig from './config';
import {orgRoutes, noOrgRoutes} from './proxy/proxy';
import { auth, requiresAuth } from 'express-openid-connect';
import AuthRoutes from './controllers/auth.controller';
import { getAuth0Config, organizationMW } from './auth/auth';
import helper from './utils/helper';

const app = express();
const authRouter = express.Router();

app.use(cors({
  origin: [/http(s)?:\/\/(.+\.)?zlifecycle\.app(:\d{1,5})?$/, /http(s)?:\/\/(.+\.)?localhost(:\d{1,5})?$/],
  methods: 'GET, HEAD, PUT, PATCH, POST, DELETE',
  preflightContinue: false,
  credentials: true,
}));

app.use(correlator({ header: 'x-correlation-id' }));
app.use(cookieParser());

// PUBLIC ROUTE
app.get('/', (req: any, res) => {
  res.redirect(zlConfig.WEB_URL);
});

// auth0 router attaches /auth/login, /logout, and /callback routes to the baseURL
app.use(auth(getAuth0Config()));

authRouter.use(requiresAuth());

authRouter.get('/session', async (req: any, res) => {
  const reqUser = helper.userFromReq(req);
  const dbUser = await getUser(reqUser.username);

  const user = {
    ...dbUser,
    picture: req.oidc.user.picture,
  };

  res.json(user);
});

app.use(noOrgRoutes(authRouter));

authRouter.use(AuthRequestLogger);
authRouter.use(organizationMW); // checks for selectedOrg cookie, throws 401 if not present

app.use('/auth', AuthRoutes(authRouter));
app.use('/', orgRoutes(authRouter));

// replaces expresses default error handler
app.use(ErrorLogger);

app.listen(process.env.PORT, () => {
  logger.info(`${process.env.npm_package_name} app listening on port ${process.env.PORT}! with ENV: ${process.env.NODE_ENV}`,
    {
      name: process.env.npm_package_name,
    })
  });

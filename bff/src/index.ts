import "./configuration";

import * as cookieParser from "cookie-parser";
import * as cors from "cors";
import * as express from "express";
import * as correlator from "express-correlation-id";
import { apiAuthMw, getUser, organizationMW, setUpAuth } from "./auth/auth";
import zlConfig from "./config";
import AuthRoutes from "./controllers/auth.controller";
import {
  externalApiRoutes,
  handlePublicRoutes,
  noOrgRoutes,
  orgRoutes,
} from "./proxy/proxy";
import helper, { oidcUser } from "./utils/helper";
import logger, { AuthRequestLogger, ErrorLogger } from "./utils/logger";
import { setUpSSE } from "./sse/sse";

const app = express();
const publicRouter = express.Router();
const externalRouter = express.Router();
const authRouter = express.Router();

app.use(
  cors({
    origin: [
      /http(s)?:\/\/(.+\.)?zlifecycle\.app(:\d{1,5})?$/,
      /http(s)?:\/\/(.+\.)?localhost(:\d{1,5})?$/,
    ],
    methods: "GET, HEAD, PUT, PATCH, POST, DELETE",
    preflightContinue: false,
    credentials: true,
  })
);

app.use(correlator({ header: "x-correlation-id" }));
app.use(cookieParser());

// PUBLIC ROUTE
app.get("/", (req: any, res) => {
  res.redirect(zlConfig.WEB_URL);
});

app.use("/public", handlePublicRoutes(publicRouter));
externalRouter.use(AuthRequestLogger);
app.use("/ext/api", apiAuthMw(), externalApiRoutes(externalRouter));

setUpAuth(app, authRouter);

authRouter.get("/session", async (req: any, res) => {
  const reqUser = helper.userFromReq(req);
  if (!reqUser) {
    res.json(null);
    return;
  }
  const dbUser = await getUser(reqUser.username);

  const user = {
    ...dbUser,
    picture: oidcUser(req).picture,
  };

  res.json(user);
});
setUpSSE(authRouter);
app.use(noOrgRoutes(authRouter));

authRouter.use(AuthRequestLogger);
authRouter.use(organizationMW); // checks for selectedOrg cookie, throws 401 if not present

app.use("/auth", AuthRoutes(authRouter));
app.use("/", orgRoutes(authRouter));

// replaces expresses default error handler
app.use(ErrorLogger);

app.listen(process.env.PORT, () => {
  logger.info(
    `${process.env.npm_package_name} app listening on port ${process.env.PORT}! with ENV: ${process.env.NODE_ENV}`,
    {
      name: process.env.npm_package_name,
    }
  );
});

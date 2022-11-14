import axios from 'axios';
import * as express from 'express';
import { ClientRequest, IncomingMessage, ServerResponse } from "http";
import { createProxyMiddleware } from "http-proxy-middleware";
import { getArgoCDAuthHeader } from '../auth/argo';
import config from '../config';
import helper from '../utils/helper';
import { Organization, User } from '../models/user.interface';
import { BFFRequest } from '../types';
import logger from '../utils/logger';
import {
  WF_MAPPINGS,
  CD_MAPPINGS,
  PathMapping,
  COSTING_MAPPINGS,
  AUDIT_MAPPINGS,
  SECRET_MAPPINGS,
  TERRAFORM_MAPPINGS,
  STATE_MAPPINGS,
  ORGANIZATION_MAPPINGS,
  EVENT_MAPPINGS,
  USERS_MAPPINGS,
} from "./pathMappings";
import { getUser } from '../auth/auth';

const proxies = new Map();

const pathRewrite = (route: string, mappings: PathMapping[], extraParams: any=null) => (path: string) => {
  for (let i = 0; i < mappings.length; i++) {
    const parts = path.split("?"); // split on get params
    const pathMatch = mappings[i].pathMatch(parts[0]);

    // Path Match
    if (pathMatch !== false) {
      // Use params from path match and add team to params
      const params = { team: 'default', ...pathMatch.params, ...extraParams };

      const url = mappings[i].newPath(params);

      if (parts.length > 1) {
        return `${url}?${parts[1]}`; // add get params back to url if they exist
      }
      
      return url;
    }
  }

  return path.replace(route, "");
};

var enableCors = function (proxyRes, req) {
  if (req.headers["access-control-request-method"]) {
    proxyRes.headers["access-control-allow-methods"] =
      req.headers["access-control-request-method"];
  }

  if (req.headers["access-control-request-headers"]) {
    proxyRes.headers["access-control-allow-headers"] =
      req.headers["access-control-request-headers"];
  }

  if (req.headers.origin) {
    proxyRes.headers["access-control-allow-origin"] = req.headers.origin;
    proxyRes.headers["access-control-allow-credentials"] = "true";
  }
};

const createProxy = function() {  
  return (org: Organization, path: string, context?: any) => {
    const orgPath = `${org.id}-${path}`;
    const organization = { name: org.name, id: org.id };

    if (!proxies.has(orgPath)) {
      proxies.set(orgPath, createProxyMiddleware({
        ...context,
        logLevel: 'info',
        onProxyRes: (proxyRes: IncomingMessage, req: IncomingMessage, res: ServerResponse) => {
          const code = proxyRes.statusCode;
          if (code >= 400 || res.statusCode >= 400) {
            if (path === "/cd" && code === 401) {
              // sending a 401 foces web into a re-auth loop with Auth0
              // when what's actually happened is an auth issue with ArgoCD not Auth0
              proxyRes.statusCode = 500;
            }

            if (code >= 400) {
              logger.info('PROXY RESPONSE', {
                proxyStatusCode: code,
                proxyPath: path,
                serverStatusCode: res.statusCode,
                serverPath: req.url,
                organization
              });
            }
          }
        },
        onProxyReq: async (proxyReq: ClientRequest, req: IncomingMessage, res: ServerResponse) => {
          if (res.statusCode >= 400) {
            logger.info('PROXY REQUEST', {
              proxyPath: proxyReq.path,
              serverStatusCode: res.statusCode,
              serverPath: req.url,
              organization
            });
          }
        }
      }));
    }

    return proxies.get(orgPath);
  }
}();

export function handlePublicRoutes(router: express.Router) : express.Router {
  // GitHub webhook proxy
  router.post('/webhook/:orgName/argocd', async (req: express.Request, res: express.Response, next) => {
    const orgName = req.params.orgName;
    const org = await helper.getOrg(orgName);

    if (!org || !orgName) {
      helper.handleNoOrg(res);
      return;
    }

    logger.log('Webhook request headers', {headers: req.headers});
    logger.log('Webhook request body', {body: req.body});

    const argoCdUrl = `${config.argoCDUrl(org.name)}/api/webhook`;

    try {
      await axios.post(argoCdUrl, {
        ...req.body
      }, {
        headers: {
          "X-GitHub-Event": req.header('X-GitHub-Event'),
          'X-GitHub-Delivery': req.header('X-GitHub-Delivery'),
          'X-Hub-Signature-256': req.header('X-Hub-Signature-256'),
        }
      });

      res.status(200).send();
    } catch (error) {
      logger.error('GitHub webhook error', { org, error });
      res.status(500).send();
      return;
    }
  });

  return router;
}

export function noOrgRoutes(router: express.Router) {
  // adds new organization
  router.post("/registration/orgs", express.json(), async (req: BFFRequest, res) => {
    const user = helper.userFromReq(req);
    try {
      const org = await axios.post(
        `${process.env.ZLIFECYCLE_API_URL}/v1/orgs/`,
        {
          name: req.body.name,
          githubRepo: req.body.githubRepo,
          termsAgreedUserId: user.id,
        }
      );

      res.json(org.data).send();
    } catch (err) {
      logger.error('create org error', { error: err.response });
      res.status(500).json(err.response.data).send();
    }
  });

  // sets selected org header
  router.post("/auth/select-org", express.json(), async (req: BFFRequest, res) => {
    const newOrgSelection = req.body.selectOrg;

    if (!newOrgSelection) {
      res.status(400).json({ message: 'selectOrg is empty'}).send();
      return;
    }

    const user = helper.userFromReq(req);
    let orgSelection;

    // query database to get org list
    const dbUser = await getUser(user.username);
    orgSelection = dbUser.organizations.find((org) => org.name === newOrgSelection);
    
    // set selectedOrg header
    if (!orgSelection) {
      res.status(400).json({ message: `${newOrgSelection} could not be selected`}).send();
      logger.error(`${newOrgSelection} could not be selected`, { orgs: req.appSession.organizations})
      return;
    }

    logger.info({message: `selected org ${orgSelection.name} for user ${user.username}`});
    
    res.cookie(config.SELECTED_ORG_HEADER, orgSelection.name, {
      httpOnly: true,
      secure: true,
      sameSite: true
    });

    res.send();
  });

  return router;
}

export function orgRoutes(router: express.Router) {  
  router.use("/wf", async (req: BFFRequest, res, next) => {    
    const org = await helper.orgFromReq(req);

    if (!org) {
      helper.handleNoOrg(res);
      return;
    }

    return (
      createProxy(org, "/wf", {
        target: config.argoWFUrl(org.name),
        pathRewrite: pathRewrite("/wf", WF_MAPPINGS, { orgName: org.name }),
        cookieDomainRewrite: "",
        onProxyRes: enableCors,
        changeOrigin: true,
      }) as any
    )(req, res, next);
  });

  router.use("/cd", async (req: BFFRequest, res, next) => {
    /* 
    Since http-proxy-middleware's are cached we need a way to inject ArgoCD tokens
    into the cached request headers. Otherwise, the cached jwt, which has a 24h TTL, 
    would expire.

    Here, we set the `authorization` header and get a valid ArgoCD token on each call.
    */
    const org = await helper.orgFromReq(req);

    if (!org) {
      helper.handleNoOrg(res);
      return;
    }
    const {authorization} = await getArgoCDAuthHeader(org.name);

    req.headers['authorization'] = authorization;

    next();
  }, async (req: BFFRequest, res, next) => {
    const org = await helper.orgFromReq(req);

    if (!org) {
      helper.handleNoOrg(res);
      return;
    }

    return (
      createProxy(org, "/cd", {
        target: config.argoCDUrl(org.name),
        changeOrigin: true,
        secure: true,
        cookieDomainRewrite: "",
        onProxyRes: enableCors,
        pathRewrite: pathRewrite("/cd", CD_MAPPINGS, { orgId: org.id, orgName: org.name }),
      }) as any
    )(req, res, next);
  });

  router.use("/costing", async (req: BFFRequest, res, next) => {
    const org = await helper.orgFromReq(req);

    if (!org) {
      helper.handleNoOrg(res);
      return;
    }

    return (
      createProxy(org, "/costing", {
        target: process.env.ZLIFECYCLE_API_URL,
        changeOrigin: true,
        secure: true,
        cookieDomainRewrite: "",
        onProxyRes: enableCors,
        pathRewrite: pathRewrite("/costing", COSTING_MAPPINGS, { orgId: org.id }),
      }) as any
    )(req, res, next);
  });

  router.use("/reconciliation", async (req: BFFRequest, res, next) => {
    const org = await helper.orgFromReq(req);

    if (!org) {
      helper.handleNoOrg(res);
      return;
    }
    const user = helper.userFromReq(req);

    return (
      createProxy(org, "/reconciliation", {
        target: process.env.ZLIFECYCLE_API_URL,
        changeOrigin: true,
        secure: true,
        cookieDomainRewrite: "",
        onProxyRes: enableCors,
        pathRewrite: pathRewrite("/", AUDIT_MAPPINGS, { orgId: org.id, email: user.email }),
      }) as any
    )(req, res, next);
  });

  router.use("/secrets", async (req: BFFRequest, res, next) => {
    const org = await helper.orgFromReq(req);

    if (!org) {
      helper.handleNoOrg(res);
      return;
    }

    return (
      createProxy(org, "/secrets", {
        target: process.env.ZLIFECYCLE_API_URL,
        changeOrigin: true,
        secure: true,
        cookieDomainRewrite: "",
        onProxyRes: enableCors,
        pathRewrite: pathRewrite("/", SECRET_MAPPINGS, { orgId: org.id }),
      }) as any
    )(req, res, next);
  });

  router.use("/orgs", async (req: BFFRequest, res, next) => {
    const org = await helper.orgFromReq(req);

    if (!org) {
      helper.handleNoOrg(res);
      return;
    }

    return (
      createProxy(org, "/orgs", {
        target: process.env.ZLIFECYCLE_API_URL,
        changeOrigin: true,
        secure: true,
        cookieDomainRewrite: "",
        onProxyRes: enableCors,
        pathRewrite: pathRewrite("/", ORGANIZATION_MAPPINGS, { orgId: org.id }),
      }) as any
    )(req, res, next);
  });

  router.use("/terraform-external", async (req: BFFRequest, res, next) => {
    const org = await helper.orgFromReq(req);

    if (!org) {
      helper.handleNoOrg(res);
      return;
    }

    return (
      createProxy(org, "/terraform-external", {
        target: "https://registry.terraform.io/",
        changeOrigin: true,
        secure: true,
        cookieDomainRewrite: "",
        onProxyRes: enableCors,
        pathRewrite: pathRewrite("/", TERRAFORM_MAPPINGS),
      }) as any
    )(req, res, next);
  });

  router.use("/terraform", async (req: BFFRequest, res, next) => {    
    const org = await helper.orgFromReq(req);

    if (!org) {
      helper.handleNoOrg(res);
      return;
    }

    return (
      createProxy(org, "/terraform", {
        target: config.stateMgrUrl(org.name),
        changeOrigin: true,
        secure: true,
        cookieDomainRewrite: "",
        onProxyRes: enableCors,
        pathRewrite: pathRewrite("/", STATE_MAPPINGS),
      }) as any
    )(req, res, next);
  });

  router.use("/users", async (req: BFFRequest, res, next) => {
    const org = await helper.orgFromReq(req);

    if (!org) {
      helper.handleNoOrg(res);
      return;
    }

    return (
      createProxy(org, "users", {
        target: process.env.ZLIFECYCLE_API_URL,
        changeOrigin: true,
        secure: true,
        cookieDomainRewrite: "",
        onProxyRes: enableCors,
        pathRewrite: pathRewrite("/", USERS_MAPPINGS, { orgId: org.id }),
      }) as any
    )(req, res, next);
  });

  router.use("/error-api", async (req: BFFRequest, res, next) => {
    const org = await helper.orgFromReq(req);

    if (!org) {
      helper.handleNoOrg(res);
      return;
    }

    return (
      createProxy(org, "/error-api", {
        target: `${config.eventApiUrl(org.name)}:8081`,
        changeOrigin: true,
        secure: false,
        cookieDomainRewrite: "",
        onProxyRes: enableCors,
        pathRewrite: { 
          "^/error-api": `/status?company=${org.name}`
        }
      })
    )(req, res, next);
  });

  router.use("/events", async (req: BFFRequest, res, next) => {
    const org = await helper.orgFromReq(req);

    if (!org) {
      helper.handleNoOrg(res);
      return;
    }

    return (
      createProxy(org, "/events", {
        target: `${config.eventApiUrl(org.name)}:8082`,
        changeOrigin: true,
        secure: false,
        cookieDomainRewrite: "",
        onProxyRes: enableCors,
        pathRewrite: pathRewrite("/", EVENT_MAPPINGS, { orgName: org.name }),
      }) as any
    )(req, res, next);
  });

  return router;
}

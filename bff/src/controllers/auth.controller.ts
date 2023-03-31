import * as express from "express";
import axios from "axios";
import logger from "../utils/logger";
import helper, { appSession, oidcUser } from "../utils/helper";
import config from "../config";
import { getArgoCDAuthHeader } from "../auth/argo";

export default function createRoutes(router: express.Router): express.Router {
  router.get("/profile", async (req: any, res) => {
    const user = {
      ...appSession(req).user,
      oidc: oidcUser(req),
    };

    res.json(user);
  });

  router.get("/me", async (req: any, res) => {
    const { user } = appSession(req);

    try {
      const org = user.organizations[0].name;
      const { data } = await axios.get(`${config.ARGOCD_URL}/api/v1/projects`, {
        headers: { ...(await getArgoCDAuthHeader(org)) },
      });

      res.json({
        id: user.id,
        name: user.name,
        terms: user.termAgreementStatus,
        role: user.role,
        groups: data.items, // not sure if this is what we need??
        username: user.username,
      });
    } catch (err) {
      logger.error("Failed getting ArgoCD projects", {
        message: err.message,
        status: err.status || err.statusCode,
      });
      throw err;
    }
  });

  router.post(
    "/setTermsAndConditions",
    async (req: any, res): Promise<boolean> => {
      const { user, body } = req;

      try {
        const url = `${process.env.ZLIFECYCLE_API_URL}/auth/setTermAgreementStatus`;
        const resp = await axios.post(url, {
          username: user.username,
          email: user.email,
          company: process.env.COMPANY,
          ...body,
        });
      } catch (err) {
        logger.error("setTermsAndConditions error", err);
        throw err;
      }

      return true;
    }
  );

  router.get("/access-token", async (req: any, res) => {
    const org = await helper.orgFromReq(req);
    const url = new URL("/oauth/token", config.AUTH0_ISSUER_BASE_URL).href;

    const data = {
      client_id: config.AUTH0_API_CLIENT_ID,
      client_secret: config.AUTH0_API_SECRET,
      audience: config.AUTH0_API_AUDIENCE,
      grant_type: "client_credentials",
      ckOrgId: org.id,
    };

    try {
      const resp = await axios.post(url, data, {
        headers: { "content-type": "application/json" },
      });

      res.send(resp.data);
    } catch (err) {
      logger.error("could not create access token", err.response.data);
      res.status(500).send({ error: "could not create access token" });
    }
  });

  return router;
}

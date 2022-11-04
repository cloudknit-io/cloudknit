import * as express from 'express';
import axios from 'axios';
import logger from '../utils/logger';
import config from '../config';
import { getArgoCDAuthHeader } from '../auth/argo';

export default function createRoutes(router: express.Router) : express.Router {
  router.get('/profile', async (req: any, res) => {
    const user = {
      ...req.appSession.user,
      oidc: req.oidc.user
    }

    res.json(user);
  });

  router.get('/me', async (req: any, res) => {
    const { user } = req.appSession;

    try {
      const org = user.organizations[0].name;
      const { data } = await axios.get(`${config.argoCDUrl(org)}/api/v1/projects`, {
        headers: { ...await getArgoCDAuthHeader(org) },
      });

      res.json({
        id: user.id,
        name: user.name,
        terms: user.termAgreementStatus,
        role: user.role,
        groups: data.items, // not sure if this is what we need??
        username: user.username
      });
    } catch (err) {
      logger.error('Failed getting ArgoCD projects', { message: err.message, status: err.status || err.statusCode });
      throw err;
    }
  });

  router.post('/setTermsAndConditions', async (req: any, res) : Promise<boolean> => {
    const { user, body } = req;

    try {
      const url = `${process.env.ZLIFECYCLE_API_URL}/auth/setTermAgreementStatus`;
      const resp = await axios.post(url, {
        username: user.username,
        email: user.email,
        company: process.env.COMPANY,
        ...body
      });
    } catch (err) {
      logger.error('setTermsAndConditions error', err);
      throw err;
    }

    return true;
  });

  return router;
}

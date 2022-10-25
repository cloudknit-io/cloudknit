import * as https from 'https';
import http from 'http';
import axios from 'axios';
import { get } from 'src/config';
import { WinstonLogger } from 'src/logger';

const logger = new WinstonLogger();

export type ProvisionOrgWf = {
  orgName: string
};

/**
 * Takes Object:
 * {
 *  orgName: 'some-org'
 * }
 * 
 * Turns into string:
 * orgName=some-org
 * @param parameters Object
 * @returns Array of 'key=value' strings representing the passed in objects key/values
 */
export function generateParams(parameters: object): Array<string> {
  const params = [];

  for (const k of Object.keys(parameters)) {
    params.push(`${k}=${parameters[k]}`);
  }

  return params;
}

async function SubmitWorkflow(resourceName: string, entryPoint: string, parameters: object) {
  const config = get();
  const url = `${config.argo.wf.url}/api/v1/workflows/${config.argo.wf.namespace}/submit`;

  const httpsAgent = new https.Agent({
    requestCert: true,
    rejectUnauthorized: false,
  });

  const data = {
    "namespace": config.argo.wf.namespace,
    resourceName,
    "resourceKind": "WorkflowTemplate",
    "submitOptions": {
      entryPoint,
      parameters: generateParams(parameters)
    }
  };

  logger.debug({message: 'Submitting provision-org workflow', url, data});

  try {
    const resp = await axios.post(url, data, {
      httpsAgent: url.startsWith('https') ? httpsAgent : null
    });
  
    logger.log(`submitted ${resourceName} workflow which generated ${resp.data.metadata.name}`)
  } catch (error) {
    if (error.response) {
      logger.error({ message: 'error submitting workflow', error: error.response.data, data, url })
    }
    
    throw error;
  }
}

export async function SubmitProvisionOrg(params: ProvisionOrgWf) {
  return SubmitWorkflow('provision-org', 'provision-org', params);
}

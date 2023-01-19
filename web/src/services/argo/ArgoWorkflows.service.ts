import { WorkflowFeedbackPayload } from 'models/projects.models';
import { Response } from 'models/response.models';
import ApiClient from 'utils/apiClient';

export interface WorkflowPayload {
	projectId: string;
	environmentId: string;
	configId?: string;
	workflowId?: string;
}

export class ArgoWorkflowsService {

	static getConfigWorkflow({
		projectId,
		environmentId,
		configId,
		workflowId,
	}: WorkflowPayload): Promise<Response<any>> {
		return ApiClient.get<any>(
			`/wf/api/v1/projects/${projectId}/environments/${environmentId}/config/${configId}/${workflowId}`
		);
	}

	static getConfigWorkflowLog({
		projectId,
		environmentId,
		configId,
		workflowId,
	}: {
		projectId: string;
		environmentId: string;
		configId?: string;
		workflowId?: string;
	}): Promise<Response<any>> {
		return ApiClient.get(
			`/wf/api/v1/projects/${projectId}/environments/${environmentId}/config/${configId}/${workflowId}/log`
		);
	}
}

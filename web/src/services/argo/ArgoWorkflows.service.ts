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

	static approveConfigWorkflow({
		projectId,
		environmentId,
		configId,
		workflowId,
		data,
	}: {
		projectId: string;
		environmentId: string;
		configId: string;
		workflowId: string;
		data: WorkflowFeedbackPayload;
	}): Promise<Response<any>> {
		return ApiClient.put<any>(
			`/wf/api/v1/projects/${projectId}/environments/${environmentId}/config/${configId}/${workflowId}/approve`,
			data
		);
	}

	static declineConfigWorkflow({
		projectId,
		environmentId,
		configId,
		workflowId,
		data,
	}: {
		projectId: string;
		environmentId: string;
		configId: string;
		workflowId: string;
		data: WorkflowFeedbackPayload;
	}): Promise<Response<any>> {
		return ApiClient.put<any>(
			`/wf/api/v1/projects/${projectId}/environments/${environmentId}/config/${configId}/${workflowId}/decline`,
			data
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

	static getNodeLog({
		projectId,
		environmentId,
		configId,
		workflowId,
		nodeId,
	}: {
		projectId: string;
		environmentId: string;
		configId?: string;
		workflowId?: string;
		nodeId?: string;
	}): Promise<Response<any>> {
		return ApiClient.get(
			`/wf/api/v1/projects/${projectId}/environments/${environmentId}/config/${configId}/${workflowId}/log/${nodeId}`
		);
	}
}

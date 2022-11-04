import { Application, ApplicationList } from 'models/argo.models';
import { EnvironmentComponentsList } from 'models/projects.models';
import { Response } from 'models/response.models';
import { ArgoMapper } from 'services/argo/ArgoMapper';
import ApiClient from 'utils/apiClient';

export class ArgoComponentsService {
	// TODO: Remove magic strings
	static getComponents(projectId: string, environmentId: string): Promise<Response<EnvironmentComponentsList>> {
		if (environmentId.toLowerCase() === 'all' && projectId.toLowerCase() === 'all') {
			return ArgoComponentsService.getAllComponents();
		}
		return ApiClient.get<ApplicationList>(
			`/cd/api/v1/projects/${projectId}/environments/${environmentId}/config`
		).then(ArgoMapper.components);
	}

	static getAllComponents(): Promise<Response<EnvironmentComponentsList>> {
		return ApiClient.get<ApplicationList>(`/cd/api/v1/config`).then(ArgoMapper.components);
	}

	static getApplicationEvents(componentName: string) {
		return ApiClient.get<any>(`/cd/api/v1/applications/${componentName}/events`);
	}

	static getApplicationResourceTree(componentName: string) {
		return ApiClient.get<any>(`/cd/api/v1/applications/${componentName}/resource-tree`);
	}

	static async patchComponentStatus(componentName: string, status: string): Promise<void> {
		const applicationData = await ApiClient.get<any>(`/cd/api/v1/component/${componentName}`);
		if (applicationData.data) {
			const payload = applicationData.data || {};
			payload.apiVersion = 'argoproj.io/v1alpha1';
			payload.kind = 'Application';
			if (payload.metadata.labels) {
				payload.metadata.labels.component_status = status;
				await ApiClient.put(`/cd/api/v1/component/${componentName}`, payload);
			}
		}
	}
}

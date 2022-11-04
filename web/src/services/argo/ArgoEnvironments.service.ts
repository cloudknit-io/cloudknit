import { ApplicationList } from 'models/argo.models';
import { EnvironmentsList } from 'models/projects.models';
import { Response } from 'models/response.models';
import { ArgoMapper } from 'services/argo/ArgoMapper';
import ApiClient from 'utils/apiClient';

export class ArgoEnvironmentsService {
	static getEnvironments(projectId: string): Promise<Response<EnvironmentsList>> {
		let url = '/cd/api/v1/environments';
		if (projectId) {
			url = `/cd/api/v1/projects/${projectId}/environments`;
		}
		return ApiClient.get<ApplicationList>(url).then(ArgoMapper.environments);
	}

	static syncEnvironment(environmentId: string): Promise<Response<any>> {
		let url = '/cd/api/v1/projects';
		if (environmentId) {
			url = `/cd/api/v1/projects/${environmentId}/sync`;
		}
		//{"revision":"HEAD","prune":false,"dryRun":false,"strategy":{"hook":{}},"resources":null}
		return ApiClient.post<ApplicationList>(url);
	}

	static deleteEnvironment(environmentId: string): Promise<Response<any>> {
		let url = '/cd/api/v1/projects';
		if (environmentId) {
			url = `/cd/api/v1/projects/${environmentId}/delete`;
		}
		return ApiClient.delete<ApplicationList>(url);
	}
}

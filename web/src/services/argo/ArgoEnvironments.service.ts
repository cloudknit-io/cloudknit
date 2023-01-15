import { ApplicationList } from 'models/argo.models';
import { EnvironmentsList } from 'models/projects.models';
import { Response } from 'models/response.models';
import ApiClient from 'utils/apiClient';

export class ArgoEnvironmentsService {

	static syncEnvironment(environmentId: string): Promise<Response<any>> {
		let url = '/cd/api/v1/projects';
		if (environmentId) {
			url = `/cd/api/v1/projects/${environmentId}/sync`;
		}
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

import { ApplicationList } from 'models/argo.models';
import { TeamsList } from 'models/projects.models';
import { Response } from 'models/response.models';
import { ArgoMapper } from 'services/argo/ArgoMapper';
import ApiClient from 'utils/apiClient';
import { ENVIRONMENT_VARIABLES } from 'utils/environmentVariables';

export class ArgoTeamsService {
	static getProjects(): Promise<Response<TeamsList>> {
		return ApiClient.get<ApplicationList>('/cd/api/v1/projects').then(ArgoMapper.teams);
	}

	static async hardSyncTeam(projectId: string, envName: string) {
		const resp = await ApiClient.post<ApplicationList>(`/cd/api/v1/projects/watcher/${projectId}/sync`, {
			revision: 'HEAD',
			prune: false,
			dryRun: false,
			strategy: {
				hook: {
					force: true,
				},
			},
			resources: [{
				group: 'stable.compuzest.com',
				kind: 'Environment',
				name: `${ENVIRONMENT_VARIABLES.REACT_APP_CUSTOMER_NAME}-${projectId}-${envName}`,
				namespace: 'zlifecycle',
				version: 'v1',
			}],
			syncOptions: null,
		});

		console.log(resp);
	}
}

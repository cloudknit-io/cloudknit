import ApiClient from 'utils/apiClient';
import { ENVIRONMENT_VARIABLES } from 'utils/environmentVariables';

export class StateFileService {
	private static customerName = ENVIRONMENT_VARIABLES.REACT_APP_CUSTOMER_NAME;
	private static instance: StateFileService;
	private constructUri = (path: string) => `/terraform/state${path}`;
	private constructOldUri = (path: string) => `/terraform/state-old${path}`;
	static getInstance() {
		if (!StateFileService.instance) {
			StateFileService.instance = new StateFileService();
		}
		return StateFileService.instance;
	}

	async fetchStateListComponents(teamId: string, environmentId: string, componentId: string, ilRepo: string) {
		try {
			const uri = this.constructUri('');
			const oldUri = this.constructOldUri('');
			const payload = {
				zstate: {
					repoUrl: ilRepo,
					meta: {
						il: `${StateFileService.customerName}-il`,
						team: teamId,
						environment: environmentId,
						component: componentId,
					},
				},
			};
			const response = await Promise.any<any>([
				ApiClient.post(uri, payload, {
					headers: {
						'Content-Type': 'application/json',
					},
				}),
				ApiClient.post(oldUri, payload, {
					headers: {
						'Content-Type': 'application/json',
					},
				}),
			]);

			return Promise.resolve(response.data);
		} catch (err) {
			Promise.reject(err);
		}
	}

	async deleteComponents(
		teamId: string,
		environmentId: string,
		componentId: string,
		ilRepo: string,
		components: string[]
	) {
		try {
			const uri = this.constructUri('');
			const oldUri = this.constructOldUri('');
			const payload = {
				data: {
					zstate: {
						repoUrl: ilRepo,
						meta: {
							il: `${StateFileService.customerName}-il`,
							team: teamId,
							environment: environmentId,
							component: componentId,
						},
					},
					resources: components,
				},
			};
			const res = await Promise.any<any>([ApiClient.delete(uri, payload), ApiClient.delete(oldUri, payload)]);
			return Promise.resolve(res.data);
		} catch (err) {
			return Promise.reject(err);
		}
	}

	async validateYaml(yaml: string) {
		try {
			return new Promise<boolean>((r, e) => setTimeout(() => r(true), 4000));
		} catch (err) {
			return false;
		}
	}
}

import { Environment, Team } from 'models/entity.store';
import { Subject } from 'rxjs';
import { BaseService } from 'services/base/base.service';
import ApiClient from 'utils/apiClient';
import { EventClientCost, subscriberCost } from 'utils/apiClient/EventClient';

export class EntityService extends BaseService {
	private static instance: EntityService | null = null;
	private constructUri = (path: string) => `/api/${path}`;

	private constructor() {
		super(Number.MAX_SAFE_INTEGER, 'entity_cache_key');
	}

	static getInstance() {
		if (!EntityService.instance) {
			EntityService.instance = new EntityService();
		}
		return EntityService.instance;
	}

	async getTeams(): Promise<Team[]> {
		const url = this.constructUri(EntitytUriType.teams());
		try {
			const resp = await ApiClient.get<Team[]>(url);
			return resp.data;
		} catch (err) {
			console.error(err);
			return [];
		}
	}

	async getEnvironments(teamId: number): Promise<Environment[]> {
		const url = this.constructUri(EntitytUriType.environments(teamId));
		try {
			const { data } = await ApiClient.get<Environment[]>(url);
			return data;
		} catch (err) {
			console.error(err);
			return [];
		}
	}
}

class EntitytUriType {
	static teams = () => `teams`;
	static environments = (teamId: number) => `teams/${teamId}/environments`;
}

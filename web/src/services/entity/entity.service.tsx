import { Component, Environment, StreamDataWrapper, Team } from 'models/entity.type';
import { BaseService } from 'services/base/base.service';
import ApiClient from 'utils/apiClient';
import { EventClient } from 'utils/apiClient/EventClient';

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

	async getTeams(withComps: boolean = false): Promise<Team[]> {
		const url = this.constructUri(EntitytUriType.teams(withComps));
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

	async syncEnvironment(teamId: number, envId: number): Promise<Environment[]> {
		const url = this.constructUri(EntitytUriType.environment(teamId, envId));
		try {
			const { data } = await ApiClient.patch<Environment[]>(url, {
				isReconcile: true
			});
			return data;
		} catch (err) {
			console.error(err);
			return [];
		}
	}

	async getComponents(teamId: number, envId: number, withLastAuditStatus: boolean = false): Promise<Component[]> {
		const url = this.constructUri(EntitytUriType.components(teamId, envId, withLastAuditStatus));
		try {
			const { data } = await ApiClient.get<Component[]>(url);
			return data;
		} catch (err) {
			console.error(err);
			return [];
		}
	}

	stream(eventList: string[]) {
		const ec = new EventClient<StreamDataWrapper>(this.constructUri(EntitytUriType.stream()), eventList);
		return ec.listen();
	}
}

class EntitytUriType {
	static teams = (withComps: boolean) => `teams?withCost=true&withEnvironments=true&withComponents=${withComps}`;
	static environments = (teamId: number) => `teams/${teamId}/environments`;
	static environment = (teamId: number, envId: number) => `teams/${teamId}/environments/${envId}`;
	static components = (teamId: number, envId: number, withLastAuditStatus: boolean) =>
		`teams/${teamId}/environments/${envId}/components?withLastAuditStatus=${withLastAuditStatus}`;
	static stream = () => `stream`;
}

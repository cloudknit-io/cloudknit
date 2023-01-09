import { Environment, Team, Component, EnvAuditData, CompAuditData } from 'models/entity.store';
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

	async getComponents(teamId: number, envId: number): Promise<Component[]> {
		const url = this.constructUri(EntitytUriType.components(teamId, envId));
		try {
			const { data } = await ApiClient.get<Component[]>(url);
			return data;
		} catch (err) {
			console.error(err);
			return [];
		}
	}

	streamComponents() {
		const ec = new EventClient<Component>(this.constructUri(EntitytUriType.streamComponents()), 'Component');
		return ec.listen();
	}

	streamEnvironments() {
		const ec = new EventClient<Environment>(this.constructUri(EntitytUriType.streamEnvironments()), 'Environment');
		return ec.listen();
	}

	streamEnvironmentAudits() {
		const ec = new EventClient<EnvAuditData>(this.constructUri(EntitytUriType.streamEnvironmentAudits()), 'EnvironmentReconcile');
		return ec.listen();
	}

	streamComponentAudits() {
		const ec = new EventClient<CompAuditData>(this.constructUri(EntitytUriType.streamComponentAudits()), 'ComponentReconcile');
		return ec.listen();
	}
}

class EntitytUriType {
	static teams = () => `teams`;
	static environments = (teamId: number) => `teams/${teamId}/environments`;
	static components = (teamId: number, envId: number) => `teams/${teamId}/environments/${envId}/components`;
	static streamEnvironments = () => `stream/environments?teamName=0&envName=0`;
	static streamComponents = () => `stream/components?teamName=0&envName=0&compName=0`;
	static streamEnvironmentAudits = () => `stream/environment/audits`;
	static streamComponentAudits = () => `stream/component/audits`;
}

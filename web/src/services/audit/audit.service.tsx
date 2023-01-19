import AuthStore from 'auth/AuthStore';
import { CompAuditData } from 'models/entity.store';
import { BaseService } from 'services/base/base.service';
import ApiClient from 'utils/apiClient';
import { ENVIRONMENT_VARIABLES } from 'utils/environmentVariables';

export class AuditService extends BaseService {
	private static instance: AuditService | null = null;
	private constructUri = (path: string) => `/reconciliation/api/v1/${path}`;

	constructor() {
		super(10000, 'audit_cache_key');
	}

	static getInstance() {
		if (!AuditService.instance) {
			AuditService.instance = new AuditService();
		}
		return AuditService.instance;
	}

	async fetchLogs(teamId: string, environmentId: string, componentId: string, auditId: number): Promise<any> {
		if (this.getCachedValue(`logs_${auditId}`)) {
			return this.getCachedValue(`logs_${auditId}`);
		}
		const resp: any = await ApiClient.get(
			this.constructUri(AuditUriType.componentLogs(teamId, environmentId, componentId, auditId))
		);

		const { data } = resp;
		if (Array.isArray(data)) {
			this.setAuditCache(`logs_${auditId}`, resp);
		}
		return resp;
	}

	async fetchStateFile(teamId: string, environmentId: string, componentId: string, cacheKey: string) {
		if (this.getCachedValue(`tfstate-${cacheKey}`)) {
			return this.getCachedValue(`tfstate-${cacheKey}`);
		}

		const resp: any = await ApiClient.get(
			this.constructUri(AuditUriType.getStateFile(teamId, environmentId, componentId))
		);

		if (!resp.data.error) {
			this.setAuditCache(`tfstate-${cacheKey}`, resp.data);
		}

		return resp.data;
	}

	async fetchPlanLogs(
		teamId: string,
		environmentId: string,
		componentId: string,
		auditId: number,
		latest: boolean
	): Promise<any> {
		if (!latest && this.getCachedValue(`plan_log_${auditId}`)) {
			return this.getCachedValue(`plan_log_${auditId}`);
		}

		const resp: any = await ApiClient.get(
			this.constructUri(AuditUriType.planLogs(teamId, environmentId, componentId, auditId, latest))
		);

		const { data } = resp;
		if (Array.isArray(data) && !latest) {
			this.setAuditCache(`plan_log_${auditId}`, resp);
		}
		return resp;
	}

	async fetchApplyLogs(
		teamId: string,
		environmentId: string,
		componentId: string,
		auditId: number,
		latest: boolean
	): Promise<any> {
		if (!latest && this.getCachedValue(`apply_log_${auditId}`)) {
			return this.getCachedValue(`apply_log_${auditId}`);
		}
		const resp: any = await ApiClient.get(
			this.constructUri(AuditUriType.applyLogs(teamId, environmentId, componentId, auditId, latest))
		);
		const { data } = resp;
		if (Array.isArray(data) && !latest) {
			this.setAuditCache(`apply_log_${auditId}`, resp);
		}
		return resp;
	}

	async getComponent(id: number, envId: number, teamId: number): Promise<CompAuditData> {
		const { data } = await ApiClient.get<CompAuditData>(
			this.constructUri(AuditUriType.component(id, envId, teamId))
		);
		return data;
	}

	async getEnvironment(envId: number, teamId: number): Promise<CompAuditData> {
		const { data } = await ApiClient.get<CompAuditData>(
			this.constructUri(AuditUriType.environment(envId, teamId)));
		return data;
	}

	async approve(configReconcileId: number) {
		return ApiClient.post(this.constructUri(AuditUriType.approve(configReconcileId)), {
			email: AuthStore.getUser()?.email
		});
	}

	setAuditCache(key: string, value: any) {
		this.cache.set(key, value);
		this.localStorageCache.setItem(key, JSON.stringify(value));
	}
}

class AuditUriType {
	static customerName = ENVIRONMENT_VARIABLES.REACT_APP_CUSTOMER_NAME;
	static component = (componentId: number, envId: number, teamId: number) =>
		`component/${teamId}/${envId}/${componentId}`;
	static approve = (componentReconcileId: number) => `component/approve/${componentReconcileId}`;
	static componentLogs = (teamId: string, environmentId: string, componentId: string, id: number) =>
		`getLogs/${teamId}/${environmentId}/${componentId}/${id}`;
	static planLogs = (teamId: string, environmentId: string, componentId: string, id: number, latest: boolean) =>
		`getPlanLogs/${teamId}/${environmentId}/${componentId}/${id}/${latest}`;
	static applyLogs = (teamId: string, environmentId: string, componentId: string, id: number, latest: boolean) =>
		`getApplyLogs/${teamId}/${environmentId}/${componentId}/${id}/${latest}`;
	static getStateFile = (teamId: string, environmentId: string, componentId: string) =>
		`getStateFile/${teamId}/${environmentId}/${componentId}`;
	static environment = (envId: number, teamId: number) => `environment/${teamId}/${envId}`;
}

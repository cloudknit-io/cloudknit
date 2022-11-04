import { Subject } from 'rxjs';
import { BaseService } from 'services/base/base.service';
import ApiClient from 'utils/apiClient';
import { EventClientAudit } from 'utils/apiClient/EventClient';
import { ENVIRONMENT_VARIABLES } from 'utils/environmentVariables';
import { ReactComponent as EKS } from 'assets/visualization-demo/platform-eks.svg';
import { ReactComponent as EC2 } from 'assets/visualization-demo/platform-ec2.svg';
import { ReactComponent as Networking } from 'assets/visualization-demo/networking.svg';
import ReactDOMServer from 'react-dom/server';
import React from 'react';

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

	private startStream(id: string, url: string): void {
		const eca = new EventClientAudit(url);
		const observer = eca.listen();
		observer.subscribe(data => {
			if (this.streamMap.get(id)?.observers.length === 0) {
				observer.unsubscribe();
				eca.close();
				return;
			}
			this.streamMap.get(id)?.next({ type: 'update', data });
		});
	}

	async getComponentInfo(id: string) {
		const resp = await ApiClient.get<any>(this.constructUri(AuditUriType.componentInfo(id)), {});
		if (resp.data) {
			return resp.data;
		}
		return null;
	}

	async patchApprovedBy(id: string) {
		const resp = await ApiClient.patch<any>(this.constructUri(AuditUriType.patchApprovedBy(id)), {});
		if (resp.data) {
			return resp.data;
		}
		return null;
	}

	async getApprovedBy(id: string, reconcileId: string) {
		const resp = await ApiClient.get<any>(this.constructUri(AuditUriType.getApprovedBy(id, reconcileId)));
		if (resp.data) {
			return resp.data;
		}
		return null;
	}

	async getEnvironmentInfo(id: string) {
		const resp = await ApiClient.get<any>(this.constructUri(AuditUriType.environmentInfo(id)), {});
		if (resp.data) {
			return resp.data;
		}
		return null;
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

	async getVisualizationURL({ team, environment, component }: any) {
		const uri = this.constructUri(AuditUriType.getVisualization());
		const payload = {
			team,
			environment,
			component,
		};
		return ApiClient.post(uri, payload);
	}

	async getVisualizationSVG({ team, environment, component }: any) {
		const uri = this.constructUri(AuditUriType.getVisualizationSVG());

		const payload = {
			path: `${team}/${environment}/${component}/visualization/rover`,
		};

		return await ApiClient.post(uri, payload);
	}

	async getVisualizationSVGDemo({ team, environment, component }: any) {
		if (component.includes('eks')) {
			return { data: ReactDOMServer.renderToString(<EKS />) };
		} else if (component.includes('networking')) {
			return { data: ReactDOMServer.renderToString(<Networking />) };
		} else if (component.includes('ec2')) {
			return { data: ReactDOMServer.renderToString(<EC2 />) };
		}
		return null;
	}

	getComponent(id: string): Subject<any> | undefined {
		this.startStream(id, this.constructUri(AuditUriType.componentStream(id)));
		return this.getStream(id, this.constructUri(AuditUriType.component(id)));
	}

	getEnvironment(id: string): Subject<any> | undefined {
		this.startStream(id, this.constructUri(AuditUriType.environmentStream(id)));
		return this.getStream(id, this.constructUri(AuditUriType.environment(id)));
	}

	async initNotifications(teamName: string) {
		return await ApiClient.get(
			this.constructUri(AuditUriType.getNotification(teamName))
		);
	}

	getNotifications(teamName: string): Subject<any> | undefined {
		const key = `notification-${AuditUriType.customerName}-${teamName}`;
		if (this.streamMap.has(key)) {
			return this.streamMap.get(key);
		}
		this.streamMap.set(key, new Subject());
		this.startStream(key, this.constructUri(AuditUriType.notificationStream(teamName)));
		return this.streamMap.get(key);
	}

	async setSeenNotification(id: string) {
		const url = this.constructUri(AuditUriType.seenNotification(id));
		await ApiClient.get(url);
	}

	setAuditCache(key: string, value: any) {
		this.cache.set(key, value);
		this.localStorageCache.setItem(key, JSON.stringify(value));
	}
}

class AuditUriType {
	static customerName = ENVIRONMENT_VARIABLES.REACT_APP_CUSTOMER_NAME;
	static environmentInfo = (componentId: string) => `environments/${componentId}`;
	static componentInfo = (componentId: string) => `components/${componentId}`;
	static component = (componentId: string) => `component/${componentId}`;
	static patchApprovedBy = (componentId: string) => `approved-by/${componentId}`;
	static getApprovedBy = (componentId: string, reconcileId: string) => `approved-by/${componentId}/${reconcileId}`;
	static componentLogs = (teamId: string, environmentId: string, componentId: string, id: number) =>
		`getLogs/${teamId}/${environmentId}/${componentId}/${id}`;
	static planLogs = (teamId: string, environmentId: string, componentId: string, id: number, latest: boolean) =>
		`getPlanLogs/${teamId}/${environmentId}/${componentId}/${id}/${latest}`;
	static applyLogs = (teamId: string, environmentId: string, componentId: string, id: number, latest: boolean) =>
		`getApplyLogs/${teamId}/${environmentId}/${componentId}/${id}/${latest}`;
	static getStateFile = (teamId: string, environmentId: string, componentId: string) =>
		`getStateFile/${teamId}/${environmentId}/${componentId}`;
	static environment = (environmentId: string) => `environment/${environmentId}`;
	static getNotification = (teamName: string) => `notifications/get/${teamName}`;
	static seenNotification = (id: string) => `notification/seen/${id}`;
	static getVisualization = () => `visualization/get`;
	static getVisualizationSVG = () => `get/object`;
	static componentStream = (componentId: string) => `components/notify/${componentId}`;
	static environmentStream = (environmentId: string) => `environments/notify/${environmentId}`;
	static notificationStream = (teamName: string) => `notifications/${teamName}`;
}

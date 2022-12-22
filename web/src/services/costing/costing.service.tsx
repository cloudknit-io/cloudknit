import { Subject } from 'rxjs';
import { BaseService } from 'services/base/base.service';
import ApiClient from 'utils/apiClient';
import { EventClientCost, subscriberCost } from 'utils/apiClient/EventClient';

export class CostingService extends BaseService {
	private static instance: CostingService | null = null;
	private constructUri = (path: string) => `/costing/api/v1/${path}`;

	constructor() {
		super(Number.MAX_SAFE_INTEGER, 'costing_cache_key');
	}

	static getInstance() {
		if (!CostingService.instance) {
			CostingService.instance = new CostingService();
		}
		return CostingService.instance;
	}

	getTeamCostStream(teamName?: string): Subject<any> {
		if (!teamName) {
			throw 'Team Name cannot be empty';
		}
		return this.getStream<number>(teamName, this.constructUri(CostingtUriType.team(teamName)));
	}

	getEnvironmentCostStream(teamName = '', environmentName = ''): Subject<any> {
		if (!teamName || !environmentName) {
			throw 'Team Name and Environment Name cannot be empty';
		}
		const key = `${teamName}-${environmentName}`;
		const url = this.constructUri(CostingtUriType.environment(teamName, environmentName));
		return this.getStream<number>(key, url);
	}

	getComponentCostStream(teamName: string, envName: string, compName: string): Subject<any> {
		const key = `${teamName}-${envName}-${compName}`;
		const url = this.constructUri(CostingtUriType.component(teamName, envName, compName));
		return this.getStream<number>(key, url);
	}

	streamTeamCost(teamId: string): any {
		new EventClientCost(`/costing/stream/api/v1/team/${teamId}`).listen();
	}

	streamEnvironmentCost(teamId: string, environmentName: string): void {
		new EventClientCost(`/costing/stream/api/v1/environment?teamName=${teamId}&envName=${environmentName}`).listen();
	}

	streamNotification(): void {
		subscriberCost.subscribe(data => {
			this.notify(data);
		});
		new EventClientCost(`/costing/stream/api/v1/notify`).listen();
	}

	streamAll(): void {
		new EventClientCost(`/costing/stream/api/v1/all`).listen();
	}

	setComponentStatus(data: any) {
		return ApiClient.post('/costing/api/v1/saveComponent', data);
	}

	private notify(data: any = null) {
		if (!data || JSON.stringify(data) === '{}') {
			return;
		}

		const { team, environment, component } = data;
		const notifier = (key: string, data: any) => {
			if (this.streamMap.has(key)) {
				this.notifySubscribers(key, data, this.streamMap.get(key) as Subject<any>);
			}
		};
		notifier(team.teamId, team.cost);
		notifier(environment.environmentId, environment.cost);
		notifier(component.id, component);
		notifier(`${component.componentId}-resources`, {
			componentId: component.componentId,
			resources: component.resources,
		});
	}
}

class CostingtUriType {
	static environment = (teamId: string, environmentId: string) => `environment?teamName=${teamId}&envName=${environmentId}`;
	static component = (teamName: string, envName: string, compName: string) => `component?teamName=${teamName}&envName=${envName}&compName=${compName}`;
	static team = (teamId: string) => `team/${teamId}`;
}

import { Application, ApplicationList, ZSyncStatus } from 'models/argo.models';
import {
	ArgoItemGeneric,
	EnvironmentComponentItem,
	EnvironmentComponentsList,
	EnvironmentItem,
	EnvironmentsList,
	TeamItem,
	TeamsList,
} from 'models/projects.models';
import { Response } from 'models/response.models';

const genericMapper = (item: Application): ArgoItemGeneric => ({
	id: item.metadata.name,
	displayValue: item.metadata.labels?.env_name || item.metadata.name || '',
	name: item.metadata.annotations?.name || item.metadata.name || '',
	labels: item.metadata.labels || {},
	healthStatus: item.status.health.status,
	syncStatus: item.status.sync.status,
	syncFinishedAt: item.status.operationState?.finishedAt,
	resourceVersion: item.metadata.resourceVersion,
	runningStatus: item.metadata?.annotations?.status || '0',
	operationPhase: item.status.operationState?.phase,
	syncResult: item.status.operationState?.syncResult,
	conditions: item.status.conditions || [],
});

export class ArgoMapper {
	static async teams(response: Response<ApplicationList>): Promise<Response<TeamsList>> {
		return {
			...response,
			data: response.data.items.map(ArgoMapper.parseTeam),
		};
	}

	static parseTeam(item: Application): TeamItem {
		return {
			...genericMapper(item),
			name: item.metadata.annotations?.name || item.metadata.name || '',
			description: item.metadata.annotations?.description,
			teamName: item.metadata.annotations?.team,
			email: item.metadata.annotations?.email,
			history: item.status.history,
			resources: item.status.resources,
			repoUrl: item.spec?.source?.repoURL,
		};
	}

	static async environments(response: Response<ApplicationList>): Promise<Response<EnvironmentsList>> {
		return {
			...response,
			data: response.data.items.map(ArgoMapper.parseEnvironment),
		};
	}

	static parseEnvironment(item: Application): EnvironmentItem {
		return {
			...genericMapper(item),
			repository: item.spec.source.repoURL,
			path: item.spec.source.path,
		};
	}

	static async components(response: Response<ApplicationList>): Promise<Response<EnvironmentComponentsList>> {
		return {
			...response,
			data: response.data.items.map(ArgoMapper.parseComponent),
		};
	}

	static parseDependsOn(labels: { [name: string]: string } = {}): string[] {
		if (!labels.depends_on_0 && !labels.depends_on) {
			return [];
		} else if (labels.depends_on_0) {
			return Object.keys(labels)
				.filter(e => e.startsWith('depends_on_'))
				.map(e => labels[e]);
		} else {
			return labels.depends_on.split('..');
		}
	}
	static parseComponent(item: Application): EnvironmentComponentItem {
		const labels = item.metadata.labels || {};
		return {
			...genericMapper(item),
			modules: item.spec.source.repoURL,
			variables: item.spec.source.repoURL,
			status: item.status.health.status,
			isDestroy: labels.is_destroy === 'true',
			isSkipped: labels.is_skipped === 'true',
			dependsOn: ArgoMapper.parseDependsOn(item.metadata.labels),
			componentName: labels.component_name || '',
			componentStatus: ZSyncStatus.Unknown,
			componentCost: '0',
		};
	}
}

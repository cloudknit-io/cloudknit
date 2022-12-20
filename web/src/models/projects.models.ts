import {
	ApplicationCondition,
	HealthStatusCode,
	OperationPhase,
	ResourceStatus,
	RevisionHistory,
	SyncOperationResult,
	SyncStatusCode,
	ZSyncStatus,
} from 'models/argo.models';

export interface ArgoItemGeneric {
	id?: string;
	name: string;
	labels?: {
		[key: string]: string;
	};
	displayValue: string;
	healthStatus: HealthStatusCode;
	syncStatus: SyncStatusCode;
	syncFinishedAt: string | undefined;
	resourceVersion?: string;
	runningStatus: string;
	operationPhase: OperationPhase | undefined;
	syncResult: SyncOperationResult | undefined;
	conditions: ApplicationCondition[];
}

export type PageHeaderTabs = PageHeaderTab[];
export interface PageHeaderTab {
	name: string | undefined;
	path: string;
	active: boolean;
}
export interface TeamItem extends ArgoItemGeneric {
	description?: string;
	teamName?: string;
	email?: string;
	children?: EnvironmentItem[];
	history?: RevisionHistory[];
	cost?: string;
	resources?: ResourceStatus[];
	repoUrl?: string;
}

export type TeamsList = TeamItem[];

export interface EnvironmentItem extends ArgoItemGeneric {
	repository: string;
	path?: string;
}

export type EnvironmentsList = EnvironmentItem[];

export interface EnvironmentComponentItem extends ArgoItemGeneric {
	modules: string;
	dependsOn: string[];
	status: SyncStatusCode | undefined;
	variables: string;
	componentName: string;
	componentStatus: ZSyncStatus;
	componentCost: string;
	isDestroy: boolean;
	isSkipped: boolean;
	costResources: [];
}

export type EnvironmentComponentsList = EnvironmentComponentItem[];

export interface WorkflowFeedbackPayload {
	name: string;
	namespace: string;
	nodeFieldSelector: string;
}

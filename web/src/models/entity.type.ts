import { AuditStatus } from "./argo.models";

export type Team = {
	id: number;
	name: string;
	estimatedCost?: number;
	environments: Environment[];
};

export type Environment = {
	argoId: string;
	id: number;
	name: string;
	lastReconcileDatetime: Date;
	duration: number;
	latestEnvRecon: EnvAuditData;
	dag: DAG[];
	teamId: number;
	status: string;
	isDeleted: boolean;
	estimatedCost: number;
	components: Component[];
	errorMessage: string[];
};

export type DAG = {
	name: string;
	type: string;
	dependsOn: string[];
};

export type Component = {
	changeId: Symbol;
	argoId: string;
	teamId: number;
	id: number;
	name: string;
	type: string;
	status: string;
	estimatedCost: number;
	lastReconcileDatetime: Date;
	duration: number;
	isDestroyed: boolean;
	costResources: any;
	dependsOn: string[];
	envId: number;
	lastWorkflowRunId: string;
	latestCompRecon: CompAuditData;
	lastAuditStatus: AuditStatus;
};

export type Update = {
	teams: Team[];
	environments: Environment[];
	components: Component[];
};

export type AuditData = {
	reconcileId: number;
	duration: number;
	status: AuditStatus;
	startDateTime: Date;
	operation?: string;
	approvedBy?: string;
	estimatedCost: number;
};

export type EnvAuditData = {
	envId: number;
	dag: DAG[];
	errorMessage: [],
} & AuditData;

export type CompAuditData = {
	compId: number;
	isDestroyed: boolean;
	costResources: any;
	lastWorkflowRunId: string;
} & AuditData;

export type StreamDataWrapper = {
	data: EnvAuditData | CompAuditData | Component | Environment | Team;
	type: StreamTypeEnum;
};

export enum StreamTypeEnum {
	Team = 'Team',
	Environment = 'Environment',
	Component = 'Component',
	ComponentReconcile = 'ComponentReconcile',
	EnvironmentReconcile = 'EnvironmentReconcile',
	Empty = 'Empty',
}
import { ReactComponent as AWSIcon } from 'assets/images/icons/AWS.svg';
import { ReactComponent as MoreOptionsIcon } from 'assets/images/icons/more-options.svg';
import { TableColumn } from 'components/atoms/table/Table';
import { renderHealthStatus, renderSyncedStatus } from 'components/molecules/cards/renderFunctions';
import { AuditStatus, ZSyncStatus } from 'models/argo.models';
import { Component, EntityStore, Environment, Team } from 'models/entity.store';
import moment from 'moment';
import React from 'react';
import { ArgoWorkflowsService } from 'services/argo/ArgoWorkflows.service';
import { FeatureKeys, VisibleFeatures } from '../feature_toggle';

const ViewTypeMap: { [key: string]: number } = {
	Concise_Logs: 0,
	Detailed_Cost_Breakdown: 2,
	Audit_View: 3,
};

if (VisibleFeatures[FeatureKeys.DETAILED_LOGS]) {
	Reflect.set(ViewTypeMap, 'Detailed_Logs', 1);
}

if (VisibleFeatures[FeatureKeys.STATE_FILE]) {
	Reflect.set(ViewTypeMap, 'State_File', 4);
}

export const ViewTypeTabName = new Map<number, string>(
	Reflect.ownKeys(ViewTypeMap).map(e => {
		const k = e as string;
		return [ViewTypeMap[k], k.replace(/\_/g, ' ')];
	})
);
export const ViewType = ViewTypeMap;

export const renderSync = (data: any) => (
	<div className="d-flex">
		{renderHealthStatus(data.healthStatus, data.componentStatus)}
		{renderSyncedStatus(data.componentStatus, data.operationPhase, data.runningStatus, '', data)}
	</div>
);

export const renderSyncStatus = (data: Component) => (
	<div className="d-flex">
		{renderSyncedStatus(data.status as ZSyncStatus, '', '', '', data)}
	</div>
);

const renderServices = () => <AWSIcon />;

const renderActions = () => <MoreOptionsIcon />;

export const getTime = (time: string): string => {
	moment.locale('en', {
		relativeTime: {
			future: 'in %s',
			past: '%s ago',
			s: '%d seconds',
			m: '%d minute',
			mm: '%d minutes',
			h: '%d hour',
			hh: '%d hours',
			d: 'a day',
			dd: '%d days',
			M: 'a month',
			MM: '%d months',
			y: 'a year',
			yy: '%d years',
		},
	});
	return moment(time, moment.ISO_8601).fromNow();
};

export const configTableColumns: TableColumn[] = [
	{
		id: 'componentName',
		name: 'Name',
		// width: 250,
	},
	{
		id: 'services',
		name: 'Services',
		// width: 100,
		render: renderServices,
	},
	{
		id: 'healthStatus',
		name: 'Status',
		// width: 200,
		combine: true,
		render: renderSyncStatus,
	},
	{
		id: 'id',
		name: 'Cost',
	},
	{
		id: 'dependsOn',
		name: 'Depends On',
		render: data => {
			return (data || []).join(' | ');
		},
	},
	{
		id: 'actions',
		name: '',
		width: 30,
		render: renderActions,
	},
];

export const auditColumns = [
	{ id: 'reconcileId', name: 'Run ID', width: 50, render: (data: any) => `#${data}` },
	{
		id: 'operation',
		name: 'Operation',
		render: (data: any) => data,
	},
	{
		id: 'status',
		name: 'Status',
		render: (data: any) => {
			const s = data.toLowerCase();
			if (s === 'ended' || s.includes('success')) {
				return renderSyncedStatus(data);
			} else if (s === 'initialising...' || s === 'initializing') {
				return renderSyncedStatus(AuditStatus.Initializing);
			} else if (s === 'failed') {
				return renderSyncedStatus(AuditStatus.Failed);
			} else {
				return renderSyncedStatus(data);
			}
		},
	},
	{ id: 'approvedBy', name: 'Approved by', render: (data: any) => data || '--' },
	{
		id: 'startDateTime',
		name: 'Started',
		render: (data: any) => getTime(new Date(data).toISOString()),
	},
	{
		id: 'duration',
		name: 'Duration',
		render: (data: number) => momentHumanizer(data),
	},
];

export const eventColumns = [
	{ id: 'reason', name: 'REASON', render: (data: any) => data, width: 140 },
	{
		id: 'message',
		name: 'MESSAGE',
		render: (data: any) => data,
	},
	{
		id: 'count',
		name: 'COUNT',
		render: (data: any) => data,
		width: 60,
	},
	{
		id: 'firstTimestamp',
		name: 'FIRST OCCURRED',
		render: (data: any) => data,
		width: 60,
	},
	{
		id: 'lastTimestamp',
		name: 'LAST OCCURRED',
		render: (data: any) => data,
		width: 60,
	},
];

export const momentHumanizer = (data: number) => {
	if (data === -1) {
		return '---';
	}
	if (data < 60 * 1000) {
		return moment.utc(data).format('s[s]');
	}

	if (data < 60 * 60 * 1000) {
		return moment.utc(data).format('m[m] s[s]');
	}

	if (data < 24 * 60 * 60 * 1000) {
		return moment.utc(data).format('h[h] m[m]');
	}

	return moment.duration(data, 'milliseconds').humanize();
};

export const getSeparatedConfigId = (config: Component) => {
	const env = EntityStore.getInstance().getEnvironmentById(config.envId) as Environment;
	const team = EntityStore.getInstance().getTeam(env?.teamId) as Team;
	return {
		team: team.name,
		component: config.name,
		environment: env.name,
	};
};

interface FetchLogsParams {
	projectId: string;
	environmentId: string;
	configId?: string;
	workflowId?: string;
}

export interface ConfigParamsSet {
	projectId: string;
	environmentId: string;
	configId?: string;
	workflowId?: string;
}

const parse = async (params: FetchLogsParams, onlyPlan?: boolean): Promise<any> => {
	const planSplitter = '------------------------------------------------------------------------';
	const showLogsStart = '----->show_output_start<-----';
	const showLogsEnd = '----->show_output_end<-----';
	const { data } = await ArgoWorkflowsService.getConfigWorkflowLog(params);
	if (onlyPlan === undefined) {
		const splitData: [] = data.split('\n') || [];
		const planItems = [],
			logItems = [];
		let fillPlanItems = false,
			cp = false;
		for (let i = 0; i < splitData.length; i++) {
			const e: string = splitData[i];
			const content = JSON.parse(e || '{}')?.result?.content || '';
			if (content === showLogsStart) {
				cp = true;
				continue;
			}
			if (content === showLogsEnd) {
				cp = false;
			}
			if ((content === planSplitter || fillPlanItems) && cp) {
				fillPlanItems = true;
				planItems.push(content);
			}
			if (cp) {
				logItems.push(content);
			}
		}
		return { planItems: planItems.join('\n'), logItems: logItems.join('\n') };
	} else {
		//Write flow for separate get
	}
};

export const getWorkflowLogs = async (
	configParamsSet: ConfigParamsSet,
	fetchWorkflowData: Function,
	setPlans: Function,
	setLogs: Function
) => {
	const workflowDataSet = configParamsSet;
	const { planItems, logItems } = await parse({ ...workflowDataSet });
	setPlans(planItems || '//');
	setLogs(logItems || '//');
};

export const processNodeLogs = (data: any) => {
	const showLogsStart = '----->show_output_start<-----';
	const showLogsEnd = '----->show_output_end<-----';
	const splitData: [] = data.split('\n') || [];
	const logItems = [];
	let cp = false;
	for (let i = 0; i < splitData.length; i++) {
		const e = splitData[i];
		const content = JSON.parse(e || '{}')?.result?.content || '';
		if (content.includes(showLogsStart)) {
			cp = true;
			continue;
		}
		if (content.includes(showLogsEnd)) {
			cp = false;
		}
		if (cp) {
			logItems.push(content);
		}
	}
	return logItems.join('\n');
};

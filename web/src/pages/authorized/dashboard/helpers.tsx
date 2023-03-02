import { ReactComponent as AWSIcon } from 'assets/images/icons/AWS.svg';
import { ReactComponent as MoreOptionsIcon } from 'assets/images/icons/more-options.svg';
import { TableColumn } from 'components/atoms/table/Table';
import { renderHealthStatus, renderLabels, renderSyncedStatus } from 'components/molecules/cards/renderFunctions';
import { ESyncStatus, HealthStatuses, ZSyncStatus } from 'models/argo.models';
import React from 'react';

import { CircularClusterPacking } from './CircularClusterPacking';
import { CloudBarchartD3 } from './CloudBarchartD3';
import { HistoryCalender } from './HistoryCalender';
import { StatusDoughnut } from './StatusDoughnut';
import { SunburstD3 } from './SunburstD3';
import { TagBarchartD3 } from './TagBarchartD3';

const renderSync = (data: any) => (
	<>
		{renderHealthStatus(data.healthStatus)}
		{renderSyncedStatus(data.syncStatus, data.operationPhase, data.runningStatus)}
	</>
);

const renderServices = () => <AWSIcon />;

const renderActions = () => <MoreOptionsIcon />;

export const projectTableColumns: TableColumn[] = [
	{
		id: 'name',
		name: 'Name',
		width: 250,
	},
	{
		id: 'team',
		name: 'Team',
		width: 100,
	},
	{
		id: 'teamEmail',
		name: 'Team email',
		width: 100,
	},
	{
		id: 'services',
		name: 'Services',
		width: 100,
		render: renderServices,
	},
	{
		id: 'labels',
		name: 'Labels',
		width: 250,
		render: renderLabels,
	},
	{
		id: 'healthStatus',
		name: 'Status',
		width: 150,
		combine: true,
		render: renderSync,
	},
	{
		id: 'description',
		name: 'Description',
	},
	{
		id: 'actions',
		name: '',
		width: 30,
		render: renderActions,
	},
];

export const d3Charts = (
	hierarchicalData: any,
	componentData: any
): { id: string; label: string; jsx: JSX.Element }[] => [
	{
		id: 'cluster',
		label: 'Team/Env/Comp Pack',
		jsx: <CircularClusterPacking data={hierarchicalData} />,
	},
	{
		id: 'sunburst',
		label: 'Team/Env/Comp Sunburst',
		jsx: <SunburstD3 data={hierarchicalData} />,
	},
	{
		id: 'barchart',
		label: 'Component Categories',
		jsx: <TagBarchartD3 data={componentData} />,
	},
	{
		id: 'doughnut',
		label: 'Sync Status',
		jsx: (
			<StatusDoughnut
				data={((data: any = []) => {
					const status: any = ZSyncStatus;
					const dd = Object.keys(status).map((e: string) => ({
						name: e,
						value: data.filter((d: any) => d.componentStatus === status[e]).length,
						components: data.filter((d: any) => d.componentStatus === status[e]),
					}));
					return dd;
				})(componentData || [])}
			/>
		),
	},
	{
		id: 'chart',
		label: 'Cloud Environment Chart',
		jsx: <CloudBarchartD3 data={componentData} />,
	},
	{
		id: 'doughnut',
		label: 'Health Status',
		jsx: (
			<StatusDoughnut
				data={((data: any = []) => {
					const status: any = HealthStatuses;
					const dd = Object.keys(status).map((e: string) => ({
						name: e,
						value: data.filter((d: any) => d.healthStatus === status[e]).length,
						components: data.filter((d: any) => d.healthStatus === status[e]),
					}));
					return dd;
				})(componentData || [])}
			/>
		),
	},
	{
		id: 'calender',
		label: 'Team History',
		jsx: <HistoryCalender data={(hierarchicalData || []).map((e: any) => e.history)} />,
	},
];

export const getClassName = (status: string): string => {
	switch (status) {
		case ZSyncStatus.Initializing:
			return '--unknown';
		case ZSyncStatus.RunningPlan:
			return '--running';
		case ZSyncStatus.CalculatingCost:
			return '--waiting';
		case ZSyncStatus.WaitingForApproval:
			return '--pending';
		case ZSyncStatus.Provisioning:
			return '--waiting';
		case ZSyncStatus.Provisioned:
			return '--successful';
		case ZSyncStatus.Destroying:
			return '--waiting';
		case ZSyncStatus.Destroyed:
			return '--successful';
		case ZSyncStatus.PlanFailed:
		case ZSyncStatus.ApplyFailed:
		case ZSyncStatus.ProvisionFailed:
		case ZSyncStatus.ValidationFailed:
		case ZSyncStatus.DestroyFailed:
		case ZSyncStatus.OutOfSync:
		case ESyncStatus.OutOfSync:
			return '--failed';
		case ZSyncStatus.InSync:
		case ESyncStatus.Synced:
			return '--successful';
		default:
			return '--unknown';
	}
};

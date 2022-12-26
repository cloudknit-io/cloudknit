import './style.scss';

import { ReactComponent as HealthDegradedIcon } from 'assets/images/icons/card-status/health/Degraded.svg';
import { ReactComponent as Skipped } from 'assets/images/icons/skipped.svg';
import { ReactComponent as HealthyIcon } from 'assets/images/icons/card-status/health/Healthy.svg';
import { ReactComponent as HealthUnknownIcon } from 'assets/images/icons/card-status/health/Unknown health.svg';
import { ReactComponent as DeleteIcon } from 'assets/images/icons/card-status/sync/delete.svg';
import { ReactComponent as CalculatingCost } from 'assets/images/icons/card-status/sync/monetization_on.svg';
import { ReactComponent as OutOfSyncIcon } from 'assets/images/icons/card-status/sync/Not Sync.svg';
import { ReactComponent as SyncedIcon } from 'assets/images/icons/card-status/sync/Sync.svg';
import { ReactComponent as Waiting } from 'assets/images/icons/card-status/sync/timer.svg';
import { ReactComponent as Hourglass } from 'assets/images/icons/hourglass.svg';
import { ReactComponent as LoaderDestroy } from 'assets/images/icons/card-status/sync/loader-destroy.svg';
import { LabelColors, ZCardLabel } from 'components/atoms/card-label/CardLabel';
import { Loader } from 'components/atoms/loader/Loader';
import { ZText } from 'components/atoms/text/Text';
import { isNumber } from 'lodash';
import {
	AuditStatus,
	ESyncStatus,
	HealthStatusCode,
	OperationPhase,
	SyncStatusCode,
	ZSyncStatus,
} from 'models/argo.models';
import moment from 'moment';
import React, { FC, ReactNode } from 'react';
import { useState } from 'react';
import { WaitingLoader } from 'components/atoms/waiting/WaitingLoader';

const BACKEND_LABELS: string[] = [
	'zlifecycle.com/model',
	'argocd.argoproj.io/instance',
	'id',
	'project_id',
	'type',
	'environment_id',
	'last_workflow_run_id',
	'status',
];
const LABEL_COLORS: LabelColors[] = ['orange', 'light-green', 'pink', 'blue', 'violet'];

//TODO (T.P.) switch to a general file
moment.locale('en', {
	relativeTime: {
		future: 'in %s',
		past: '%s ago',
		s: 'seconds',
		ss: '%ss',
		m: 'a minute',
		mm: '%dm',
		h: 'an hour',
		hh: '%dh',
		d: 'a day',
		dd: '%dd',
		M: 'a month',
		MM: '%dM',
		y: 'a year',
		yy: '%dY',
	},
});

export const currency = (cost: number) => Number(cost).toFixed(2);

export const CostRenderer: FC<any> = ({ data }: any) => {
	const cost = Number(data?.cost || data || 0);
	return <>{isNumber(cost) ? `${cost == -1 ? 'N/A' : '$' + currency(cost)}` : 'calculating cost...'}</>;
};

export const renderLabels = (labels: { [name: string]: string }): ReactNode[] => {
	return Object.entries(labels)
		.filter(([key]) => {
			return !BACKEND_LABELS.includes(key);
		})
		.map(([key, value], index) => (
			<ZCardLabel
				key={'label' + key + index}
				text={key + '=' + value}
				color={LABEL_COLORS[index % LABEL_COLORS.length]}
			/>
		));
};

type StatusDisplayProps = {
	icon?: ReactNode;
	text: string;
	time?: string | undefined;
	title?: string;
};

const StatusDisplay: FC<StatusDisplayProps> = ({ icon, text, time, title }) => {
	const [showTooltip, setShowTooltip] = useState<boolean>(false);
	const getTime = (time: string): string => {
		if (!time) {
			return '';
		}
		return `| ${moment(time, moment.ISO_8601).fromNow()}`;
	};

	return (
		<div
			className="zlifecycle-status-display"
			title={`${text} ${getTime(time || '')}`}
			onMouseOver={() => setShowTooltip(true)}
			onMouseOut={() => setShowTooltip(false)}>
			{title && showTooltip && (
				<div className={`zlifecycle-status-display_tooltip ${showTooltip ? 'show' : ''}`}>{title}</div>
			)}
			{icon}
			<ZText.Body className="zlifecycle-status-display__text" size="14" lineHeight="18">
				{text}
			</ZText.Body>
			{time && (
				<ZText.Body className="zlifecycle-status-display--time" size="14" lineHeight="18">
					{getTime(time)}
				</ZText.Body>
			)}
		</div>
	);
};

export const renderHealthIcon = (healthStatus: HealthStatusCode): ReactNode => {
	switch (healthStatus) {
		case 'Healthy':
			return <HealthyIcon />;
		case 'Degraded':
			return <HealthDegradedIcon />;
		case 'Unknown':
			return <HealthUnknownIcon />;
		default:
			return <HealthUnknownIcon />;
	}
};

export const renderHealthStatus = (healthStatus: HealthStatusCode, componentStatus?: ZSyncStatus): ReactNode => {
	if (componentStatus === ZSyncStatus.Destroyed) {
		return <></>;
	}
	switch (healthStatus) {
		case 'Healthy':
			return <StatusDisplay text={'Healthy'} icon={<HealthyIcon />} />;
		case 'Degraded':
			return <StatusDisplay text={'Degraded'} icon={<HealthDegradedIcon />} />;
		case 'Unknown':
			return <StatusDisplay text={'Unknown'} icon={<HealthUnknownIcon />} />;
		default:
			return <StatusDisplay text={'Unknown'} icon={<HealthUnknownIcon />} />;
	}
};

export const renderEnvSyncedStatus = (
	componentStatus: ZSyncStatus | SyncStatusCode | AuditStatus,
	operationPhase?: any,
	runningStatus?: string,
	syncFinishedAt?: string,
	data?: any
): ReactNode => {
	// console.log(data);
	// console.log(syncFinishedAt);
	switch (componentStatus) {
		case ZSyncStatus.Initializing:
			return (
				<StatusDisplay text={'Initializing'} icon={<Loader height={20} width={20} />} time={syncFinishedAt} />
			);
		case ZSyncStatus.RunningPlan:
			return (
				<StatusDisplay text={'Running Plan'} icon={<Loader height={20} width={20} />} time={syncFinishedAt} />
			);
		case ZSyncStatus.RunningDestroyPlan:
			return (
				<StatusDisplay
					text={'Running Plan'}
					icon={<Loader height={20} width={20} />}
					time={syncFinishedAt}
				/>
			);
		case ZSyncStatus.CalculatingCost:
			return (
				<StatusDisplay
					text={'Calculating Est. Cost'}
					icon={<CalculatingCost height={20} width={20} />}
					time={syncFinishedAt}
				/>
			);
		case ZSyncStatus.WaitingForApproval:
			return (
				<StatusDisplay
					text={'Waiting For Approval'}
					icon={<Waiting height={20} width={20} />}
					time={syncFinishedAt}
				/>
			);
		case ZSyncStatus.Provisioning:
			return (
				<StatusDisplay text={'Provisioning'} icon={<Loader height={20} width={20} />} time={syncFinishedAt} />
			);
		case ZSyncStatus.Provisioned:
			return <StatusDisplay text={'Provisioned'} icon={<SyncedIcon />} time={syncFinishedAt} />;
		case ZSyncStatus.Destroying:
			return <StatusDisplay text={'Destroying'} icon={<Loader height={20} width={20} />} time={syncFinishedAt} />;
		case ZSyncStatus.Destroyed:
			return (
				<StatusDisplay text={'Destroyed'} icon={<DeleteIcon height={16} width={16} />} time={syncFinishedAt} />
			);
		case ZSyncStatus.PlanFailed:
			return <StatusDisplay text={'Plan Failed'} icon={<OutOfSyncIcon />} time={syncFinishedAt} />;
		case ZSyncStatus.ApplyFailed:
			return <StatusDisplay text={'Apply Failed'} icon={<OutOfSyncIcon />} time={syncFinishedAt} />;
		case ZSyncStatus.ProvisionFailed:
			return <StatusDisplay text={'Provision Failed'} icon={<OutOfSyncIcon />} time={syncFinishedAt} />;
		case ZSyncStatus.DestroyFailed:
			return <StatusDisplay text={'Destroy Failed'} icon={<OutOfSyncIcon />} time={syncFinishedAt} />;
		case ZSyncStatus.OutOfSync:
		case ESyncStatus.OutOfSync:
			return <StatusDisplay text={'Out of Sync'} icon={<OutOfSyncIcon />} time={syncFinishedAt} />;
		case ZSyncStatus.InSync:
		case ESyncStatus.Synced:
			return (
				<StatusDisplay
					text={'In Sync'}
					icon={<SyncedIcon height={16} width={16} className="calculate_cost_animate" />}
					time={syncFinishedAt}
				/>
			);
		default:
			const dependsOn = data?.dependsOn?.filter((d: any) => d !== 'root');
			if (dependsOn?.length > 0 && !data.isDestroy) {
				return (
					<StatusDisplay
						title={`${dependsOn.join(', ')}`}
						text={'Waiting for Parent'}
						icon={<Hourglass height={20} width={20} />}
					/>
				);
			} else {
				return (
					<StatusDisplay
						text={'Initializing'}
						icon={<Loader height={20} width={20} />}
						time={syncFinishedAt}
					/>
				);
			}
	}
};


export const renderSyncedStatus = (
	componentStatus: ZSyncStatus | SyncStatusCode | AuditStatus,
	operationPhase?: any,
	runningStatus?: string,
	syncFinishedAt?: string,
	data?: any
): ReactNode => {
	// console.log(data);
	// console.log(syncFinishedAt);
	switch (componentStatus) {
		case ZSyncStatus.Initializing:
			return <StatusDisplay text={'Initializing'} icon={<Loader height={20} width={20} />} time={syncFinishedAt} />
		case ZSyncStatus.InitializingApply:
			return <StatusDisplay text={'Initializing Apply'} icon={<Loader height={20} width={20} />} time={syncFinishedAt} />
		case AuditStatus.Initializing:
		case AuditStatus.Env_Destroying:
		case AuditStatus.Env_Provisioning:
		case AuditStatus.Destroying:
		case AuditStatus.Provisioning:
			return (
				<StatusDisplay text={'In Progress'} icon={<Loader height={20} width={20} />} time={syncFinishedAt} />
			);
		case AuditStatus.Destroyed:
		case AuditStatus.Provisioned:
		case AuditStatus.Success:
		case AuditStatus.Ended:
		case AuditStatus.Env_Destroy_Ended:
		case AuditStatus.Env_Provision_Ended:
			return (
				<StatusDisplay text={'Succeeded'} icon={<SyncedIcon />} time={syncFinishedAt} />
			);
		case ZSyncStatus.RunningPlan:
			return (
				<StatusDisplay text={'Running Plan'} icon={<Loader height={20} width={20} />} time={syncFinishedAt} />
			);
		case ZSyncStatus.RunningDestroyPlan:
			return (
				<StatusDisplay
					text={'Running Plan'}
					icon={<Loader height={20} width={20} />}
					time={syncFinishedAt}
				/>
			);
		case ZSyncStatus.CalculatingCost:
			return (
				<StatusDisplay
					text={'Calculating Est. Cost'}
					icon={<CalculatingCost height={20} width={20} />}
					time={syncFinishedAt}
				/>
			);
		case ZSyncStatus.WaitingForApproval:
			return (
				<StatusDisplay
					text={'Waiting For Approval'}
					icon={<Waiting height={20} width={20} />}
					time={syncFinishedAt}
				/>
			);
		case ZSyncStatus.Provisioning:
			return (
				<StatusDisplay text={'Applying'} icon={<Loader height={20} width={20} />} time={syncFinishedAt} />
			);
		case ZSyncStatus.Provisioned:
			return <StatusDisplay text={'Succeeded'} icon={<SyncedIcon />} time={syncFinishedAt} />;
		case ZSyncStatus.Destroying:
			return <StatusDisplay text={'Applying'} icon={<Loader height={20} width={20} />} time={syncFinishedAt} />;
		case ZSyncStatus.Destroyed:
			return (
				<StatusDisplay text={'Succeeded'} icon={<DeleteIcon height={16} width={16} />} time={syncFinishedAt} />
			);
		case AuditStatus.Skipped:
		case AuditStatus.SkippedDestroy:
			return (
				<StatusDisplay
					text={'Skipped Destroy'}
					icon={<Skipped height={16} width={16} />}
					time={syncFinishedAt}
				/>
			);
		case AuditStatus.SkippedProvision:
			return (
				<StatusDisplay
					text={'Skipped Provision'}
					icon={<Skipped height={16} width={16} />}
					time={syncFinishedAt}
				/>
			);
		case AuditStatus.SkippedReconcile:
			return (
				<StatusDisplay
					text={'Skipped Reconcile'}
					icon={<Skipped height={16} width={16} />}
					time={syncFinishedAt}
				/>
			);
		case ZSyncStatus.NotProvisioned:
			return (
				<StatusDisplay
					text={'Not Provisioned'}
					icon={<Skipped height={16} width={16} />}
					time={syncFinishedAt}
				/>
			);
		case AuditStatus.Failed:
			return <StatusDisplay text={'Failed'} icon={<OutOfSyncIcon />} />;
		case ZSyncStatus.PlanFailed:
		case AuditStatus.DestroyPlanFailed:
		case AuditStatus.ProvisionPlanFailed:
			return <StatusDisplay text={'Plan Failed'} icon={<OutOfSyncIcon />} time={syncFinishedAt} />;
		case ZSyncStatus.ApplyFailed:
		case AuditStatus.DestroyApplyFailed:
		case AuditStatus.ProvisionApplyFailed:
			return <StatusDisplay text={'Apply Failed'} icon={<OutOfSyncIcon />} time={syncFinishedAt} />;
		case ZSyncStatus.OutOfSync:
		case ESyncStatus.OutOfSync:
			return <StatusDisplay text={'Out of Sync'} icon={<OutOfSyncIcon />} time={syncFinishedAt} />;
		case AuditStatus.Failed:
			return <StatusDisplay text={'Failed'} icon={<OutOfSyncIcon />} time={syncFinishedAt} />;
		case ZSyncStatus.InSync:
		case ESyncStatus.Synced:
			return (
				<StatusDisplay
					text={'In Sync'}
					icon={<SyncedIcon height={16} width={16} className="calculate_cost_animate" />}
					time={syncFinishedAt}
				/>
			);
		default:
			const dependsOn = data?.dependsOn?.filter((d: any) => d !== 'root');
			if (dependsOn?.length > 0 && !data.isDestroy) {
				return (
					<StatusDisplay
						title={`${dependsOn.join(', ')}`}
						text={'Waiting for Parent'}
						icon={<Hourglass height={20} width={20} />}
					/>
				);
			} else {
				return (
					<StatusDisplay
						text={'Initializing'}
						icon={<Loader height={20} width={20} />}
						time={syncFinishedAt}
					/>
				);
			}
	}
};

export const getHealthStatusIcon = (healthStatus: any) => {
	switch (healthStatus) {
		case 'Healthy':
			return <HealthyIcon title="Healthy" />;
		case 'Degraded':
			return <HealthDegradedIcon title="Degraded" />;
		case 'Unknown':
			return <HealthUnknownIcon title="Unknown" />;
		default:
			return <HealthUnknownIcon title="Unknown" />;
	}
};

export const getSyncStatusIcon = (syncStatus: any, operation?: 'Destroy' | 'Provision') => {
	switch (syncStatus) {
		case ZSyncStatus.Initializing:
			return operation === 'Destroy' ? <LoaderDestroy height={16} width={16} title="Initializing" /> : <Loader height={16} width={16} title="Initializing" />;
		case ZSyncStatus.InitializingApply:
			return operation === 'Destroy' ? <LoaderDestroy height={16} width={16} title="Initializing Apply" /> : <Loader height={16} width={16} title="Initializing Apply" />;
		case ZSyncStatus.RunningPlan:
			return <Loader height={16} width={16} title="Running Plan" />;
		case ZSyncStatus.RunningDestroyPlan:
			return <LoaderDestroy height={16} width={16}  title="Running Destroy Plan" />;
		case ZSyncStatus.CalculatingCost:
			return <CalculatingCost height={16} width={16} title="Calculating Cost" />;
		case ZSyncStatus.WaitingForApproval:
			return <Waiting height={16} width={16} title="Waiting For Approval" />;
		case ZSyncStatus.Provisioning:
			return <Loader height={16} width={16} title="Provisioning" />;
		case ZSyncStatus.Provisioned:
			return <SyncedIcon title="Provisioned" />;
		case ZSyncStatus.Destroying:
			return <LoaderDestroy height={16} width={16} title="Destroying" />; // <Loader height={16} width={16} title="Destroying" />;
		case ZSyncStatus.Destroyed:
			return <DeleteIcon title="Destroyed" height={16} width={16} />;
		case ZSyncStatus.Skipped:
		case ZSyncStatus.NotProvisioned:
			return <Skipped height={16} width={16} title="Skipped" />;
		case ZSyncStatus.SkippedReconcile:
			return <Skipped height={16} width={16} title="SkippedReconcile" />;
		case ZSyncStatus.PlanFailed:
		case ZSyncStatus.ApplyFailed:
		case ZSyncStatus.ProvisionFailed:
		case ZSyncStatus.DestroyFailed:
			return <OutOfSyncIcon title="Failed" />;
		case ZSyncStatus.OutOfSync:
		case ESyncStatus.OutOfSync:
			return <OutOfSyncIcon title="Out of Sync" />;
		case ZSyncStatus.InSync:
		case ESyncStatus.Synced:
			return <SyncedIcon height={16} width={16} />;
		default:
			return <WaitingLoader radius={3} title="Waiting For Parent" />;
	}
};

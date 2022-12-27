import { ESyncStatus, ZSyncStatus } from 'models/argo.models';
import moment from 'moment';
import React, { ReactElement } from 'react';
import { getSVGNode } from './shape-helper';
import { DagNodeProps, ZDagNode } from 'components/molecules/dag-node/DagNode';
import { ZDagAppNode } from 'components/molecules/dag-node/DagAppNode';
import ReactDOM from 'react-dom';
import { AuditService } from 'services/audit/audit.service';
import { Subject } from 'rxjs';

export const getClassName = (status: string): string => {
	switch (status) {
		case ZSyncStatus.Initializing:
		case ZSyncStatus.InitializingApply:
			return '--initializing';
		case ZSyncStatus.RunningPlan:
		case ZSyncStatus.RunningDestroyPlan:
			return '--running';
		case ZSyncStatus.CalculatingCost:
			return '--running';
		case ZSyncStatus.WaitingForApproval:
			return '--pending';
		case ZSyncStatus.Provisioning:
		case ZSyncStatus.Destroying:
			return '--waiting';
		case ZSyncStatus.Provisioned:
			return '--successful';
		case ZSyncStatus.Destroyed:
		case ZSyncStatus.NotProvisioned:
			return '--destroyed';
		case ZSyncStatus.Skipped:
			return '--skipped';
		case ZSyncStatus.SkippedReconcile:
			return '--skipped-reconcile';
		case ZSyncStatus.PlanFailed:
		case ZSyncStatus.ApplyFailed:
		case ZSyncStatus.ProvisionFailed:
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

export const getTextWidth = (name: string): number => {
	const canvas = document.createElement('canvas');
	const ctx: CanvasRenderingContext2D = canvas.getContext('2d') as CanvasRenderingContext2D;
	ctx.font = 'bold 15px "DM Sans"';
	const { width } = ctx.measureText(name);
	canvas.remove();
	return width;
};

export const getEnvironment = async (envName: string, teamName: string) => {
	return await AuditService.getInstance().getEnvironmentInfo(envName, teamName);
};

export const getLastReconcileTime = async (id: string, envName: string, teamName: string, defaultValue: string) => {
	try {
		const info = await AuditService.getInstance().getComponentInfo(
			id,
			envName,
			teamName
		);
		if (!info) {
			return defaultValue;
		}
		return info['lastReconcileDatetime'];
	} catch {
		return defaultValue;
	}
};

export const getTime = (time: string): string => {
	return moment(time, moment.ISO_8601).fromNow();
};

type DagCacheProps = {
	container: HTMLElement;
	updater: Subject<DagNodeProps>;
	node: JSX.Element;
};

const dagNodeCache = new Map<string, DagCacheProps>();

export const cleanDagNodeCache = (projectId: string) => {
	if (projectId && !dagNodeCache.has(projectId)) {
		const dagNodes = [...dagNodeCache.values()];
		dagNodes.forEach(d => ReactDOM.unmountComponentAtNode(d.container));
		dagNodeCache.clear();
	}
};

export function updateNodeFigure({
	id,
	name,
	icon,
	componentStatus,
	displayValue,
	syncStatus,
	healthStatusIcon,
	syncFinishedAt,
	projectId,
	isSkipped,
	labels,
	expandIcon,
	onNodeClick,
	estimatedCost,
	isDestroy
}: any) {
	const nodeId = id === 'root' ? name : displayValue;
	const props = dagNodeCache.get(nodeId);
	props?.updater?.next({
		componentStatus,
		displayValue,
		Icon: icon,
		id,
		isSkipped,
		name,
		syncFinishedAt,
		SyncStatus: syncStatus,
		projectId,
		onNodeClick,
		estimatedCost,
		operation: labels.env_status === 'destroying' || isDestroy ? 'Destroy' : 'Provision',
		updater: props.updater,
	});
	return dagNodeCache.get(nodeId);
}

function createNodeFigure({
	id,
	name,
	icon,
	componentStatus,
	displayValue,
	syncStatus,
	healthStatusIcon,
	syncFinishedAt,
	projectId,
	isSkipped,
	labels,
	expandIcon,
	onNodeClick,
	estimatedCost,
	isDestroy
}: any) {
	const nodeId = id === 'root' ? name : displayValue;
	const isApp = labels?.component_type !== 'terraform';
	const groupNode = getSVGNode(
		{
			class: 'roundedCorners',
			id: id,
			style: 'cursor: pointer',
		},
		'g'
	);

	let width = getTextWidth(name);
	const timeWidth = 40 + getTextWidth(getTime(syncFinishedAt));
	width = timeWidth > width ? timeWidth : width;
	const rectContainer = getSVGNode(
		{
			height: '70px',
			width: `${width + 100 > 160 ? width + 100 : 160}px`,
			fill: 'transparent',
			rx: '10',
		},
		'rect'
	);
	let node = <></>;
	const updater = new Subject<DagNodeProps>();
	if (isApp) {
		node = (
			<ZDagAppNode
				componentStatus={componentStatus}
				displayValue={displayValue}
				Icon={icon}
				id={id}
				isSkipped={isSkipped}
				name={name}
				syncFinishedAt={syncFinishedAt}
				HealthIcon={healthStatusIcon}
				ExpandIcon={expandIcon}
				SyncStatus={syncStatus}
				onNodeClick={onNodeClick}
			/>
		);
	} else {
		node = (
			<ZDagNode
				componentStatus={componentStatus}
				displayValue={displayValue}
				Icon={icon}
				id={id}
				isSkipped={isSkipped}
				name={name}
				syncFinishedAt={syncFinishedAt}
				SyncStatus={syncStatus}
				projectId={projectId}
				onNodeClick={onNodeClick}
				estimatedCost={estimatedCost}
				labels={labels}
				operation={labels.env_status === 'destroying' || isDestroy ? 'Destroy' : 'Provision'}
				updater={updater}
			/>
		);
		dagNodeCache.set(nodeId, { container: groupNode, updater, node });
	}

	ReactDOM.render(node, groupNode);
	groupNode.appendChild(rectContainer);
	return groupNode;
}

export default createNodeFigure;

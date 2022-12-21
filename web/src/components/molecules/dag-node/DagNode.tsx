import React, { FC, useEffect, useState } from 'react';
import { CostRenderer, getSyncStatusIcon } from '../cards/renderFunctions';
import { CostingService } from 'services/costing/costing.service';
import ReactDOMServer from 'react-dom/server';
import { FeatureKeys, featureToggled } from 'pages/authorized/feature_toggle';
import {
	getClassName,
	getEnvironment,
	getLastReconcileTime,
	getTextWidth,
	getTime,
} from 'components/organisms/tree-view/node-figure-helper';
import { ZSyncStatus } from 'models/argo.models';
import { Subject, Subscription, pipe } from 'rxjs';
import { debounceTime } from 'rxjs/operators';

export type DagNodeProps = {
	id: string;
	componentStatus: ZSyncStatus;
	name: string;
	syncFinishedAt: string;
	isSkipped: boolean;
	displayValue: string;
	projectId: string;
	SyncStatus: ZSyncStatus;
	Icon: JSX.Element;
	estimatedCost: string;
	operation: 'Provision' | 'Destroy';
	onNodeClick: (...params: any) => any;
	updater: Subject<DagNodeProps>;
};

export const ZDagNode: FC<DagNodeProps> = ({
	id,
	componentStatus,
	name,
	syncFinishedAt,
	isSkipped,
	displayValue,
	SyncStatus,
	Icon,
	projectId,
	operation,
	estimatedCost,
	onNodeClick,
	updater,
}: DagNodeProps) => {
	const [cost, updateCost] = useState<number | null>(
		id === 'root' ?
		CostingService.getInstance().getCachedValue(name) : Number(estimatedCost || 0));
	const [syncTime, setSyncTime] = useState<any>(syncFinishedAt);
	const [status, setStatus] = useState<ZSyncStatus>(componentStatus);
	const [syncStatus, setSyncStatus] = useState<ZSyncStatus>(SyncStatus);
	const [skippedStatus, setSkippedStatus] = useState<boolean>(isSkipped);
	const [operationType, setOperationType] = useState<"Provision" | "Destroy">(operation);
	let width = getTextWidth(name);
	const timeWidth = 40 + getTextWidth(getTime(syncTime));
	width = timeWidth > width ? timeWidth : width;
	const rectWidth = width + 100 > 160 ? width + 100 : 160;

	useEffect(() => {
		let $subscription: Subscription[] = [];
		if (id === 'root') {
			getEnvironment(name, projectId)
				.then(env => setSyncTime(env['lastReconcileDatetime']))
				.catch(e => setSyncTime(syncFinishedAt));
			$subscription.push(
				CostingService.getInstance()
					.getEnvironmentCostStream(projectId, name.replace(projectId + '-', ''))
					.subscribe(data => {
						updateCost(data);
					})
			);
		} else {
			getLastReconcileTime(displayValue, syncFinishedAt).then(r => setSyncTime(r));
		}
		$subscription.push(
			updater.pipe(debounceTime(1000)).subscribe(async (data: DagNodeProps) => {
				if (data.estimatedCost) {
					updateCost(Number(data.estimatedCost || 0))
				}
				setStatus(data.componentStatus);
				setSyncStatus(data.SyncStatus);
				setSkippedStatus(data.isSkipped);
				setOperationType(data.operation);
				if (id === 'root') {
					getEnvironment(name, projectId)
						.then(env => setSyncTime(env['lastReconcileDatetime']))
						.catch(e => setSyncTime(syncFinishedAt));
				} else {
					const time = await getLastReconcileTime(displayValue, syncFinishedAt);
					setSyncTime(time);
				}
			})
		);

		return () => $subscription.forEach(e => e.unsubscribe());
	}, []);

	return (
		<g
			className="react-dag-node"
			filter={`${status === ZSyncStatus.Destroyed ? 'grayscale(1)' : ''}`}
			onClick={e => {
				onNodeClick(id);
			}}>
			<rect
				height={70}
				width={rectWidth}
				fill="#fff"
				rx="10"
				className={`node node__pod ${id === 'root' ? 'root' : 'node__pod'}${getClassName(status || '')} ${
					skippedStatus ? ' striped' : ''
				}`}
			/>
			<text x="65" y="22" fill="#323232" fontFamily="DM Sans" fontWeight={'light'} fontSize="15px">
				{name}
			</text>
			{syncFinishedAt && (
				<text x="85" y="61" fill="#323232" fontFamily="DM Sans" fontWeight={'light'} fontSize="14px">
					{' | ' + getTime(syncTime)}
				</text>
			)}
			{operationType === 'Destroy' && status !== ZSyncStatus.Destroyed && (
				<g transform={`translate(${rectWidth - 20},${51})`}>{getSyncStatusIcon(ZSyncStatus.Destroyed)}</g>
			)}
			<text x="65" y="42" fill="#323232" fontFamily="DM Sans" fontWeight={'light'} fontSize="14px">
				<CostRenderer data={cost} />
			</text>
			{syncStatus && <g transform={`translate(${65},${48})`}>{getSyncStatusIcon(syncStatus, operationType)}</g>}
			{Icon && (
				<g
					onClick={e => {
						if (!featureToggled(FeatureKeys.VISUALIZATION)) {
							return;
						}
						if (id === 'root') {
							return;
						}
						(e.nativeEvent as any)['visualizationHandler'] = true;
					}}
					transform={`translate(${12},${12}) scale(${id === 'root' ? 0.25 : 0.35})`}
					dangerouslySetInnerHTML={{
						__html: ReactDOMServer.renderToString(Icon),
					}}></g>
			)}
		</g>
	);
};

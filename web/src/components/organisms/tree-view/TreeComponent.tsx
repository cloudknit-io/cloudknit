import { ReactComponent as ArrowUp } from 'assets/images/icons/chevron-right.svg';
import { ReactComponent as ComputeIcon } from 'assets/images/icons/DAG-View/config.svg';
import { ReactComponent as LayersIcon } from 'assets/images/icons/DAG-View/environment-icon.svg';
import { NotificationsApi, Select } from 'components/argo-core';
import { getHealthStatusIcon, getSyncStatusIcon } from 'components/molecules/cards/renderFunctions';
import { useApi } from 'hooks/use-api/useApi';
import { EnvironmentComponentItem, EnvironmentItem } from 'models/projects.models';
import React, { FC, useEffect, useMemo, useState } from 'react';
import { ArgoEnvironmentsService } from 'services/argo/ArgoEnvironments.service';
import { ReactComponent as SyncIcon } from 'assets/images/icons/sync-icon.svg';

import Tree from './TreeView';
import { ESyncStatus, OperationPhase, ResourceResult, ZSyncStatus } from 'models/argo.models';
import { getEnvironmentErrorCondition, syncMe } from 'pages/authorized/environments/helpers';
import { Context } from 'context/argo/ArgoUi';
import { subscriberWatcher } from 'utils/apiClient/EventClient';
import { ReactComponent as Expand } from 'assets/images/icons/expand.svg';
import { ReactComponent as Add } from 'assets/images/icons/add.svg';
import { ReactComponent as Subtract } from 'assets/images/icons/subtract.svg';
import { ReactComponent as AppIcon } from 'assets/images/icons/DAG-View/Layers.svg';
import { cleanDagNodeCache } from './node-figure-helper';
import { Component, EntityStore, Environment } from 'models/entity.store';

const curveTypes = [
	{ value: '0', title: 'Curve Basis' },
	{ value: '1', title: 'Curve Bundle' },
	{ value: '2', title: 'Curve Cardinal' },
	{ value: '3', title: 'Curve CatmullRom' },
	{ value: '4', title: 'Curve Linear' },
	{ value: '5', title: 'Curve MonotoneX' },
	{ value: '6', title: 'Curve Natural' },
	{ value: '7', title: 'Curve Step' },
];

const networkTypes = [
	{
		value: 'network-simplex',
		title: 'Network Simplex',
	},
	{
		value: 'tight-tree',
		title: 'Tight Tree',
	},
	{
		value: 'longest-path',
		title: 'Longest Path',
	},
];

export const getShape = (name: any) => {
	switch (name) {
		case 'network':
			return 'ellipse';
		default:
			return 'ellipse';
	}
};

interface Props {
	environmentId: string;
	nodes: any;
	onNodeClick: any;
	environmentItem?: Environment;
}

export const TreeComponent: FC<Props> = ({ environmentId, nodes, onNodeClick, environmentItem }: Props) => {
	const entityStore = useMemo(() => EntityStore.getInstance(), []);
	const [data, setData] = useState<any[]>([]);
	const [option, setOption] = useState('0');
	const [rankSep, setRankSep] = useState(70);
	const [nodeSep, setNodeSep] = useState(100);
	const [modifierMizimized, toggleMinimize] = useState(true);
	const [syncStarted, setSyncStarted] = useState<boolean>(false);
	const [ranker, setRanker] = useState('network-simplex');
	const [watcherStatus, setWatcherStatus] = useState<OperationPhase | undefined>();
	let zoomIn: () => void = () => {};
	let zoomOut: () => void = () => {};
	let reset: () => void = () => {};
	const nm = React.useContext(Context)?.notifications;
	const syncEnvironment = useApi(ArgoEnvironmentsService.syncEnvironment);
	const deleteEnvironment = useApi(ArgoEnvironmentsService.deleteEnvironment);
	const [environmentCondition, setEnvironmentCondition] = useState<any>(null);
	const expandedNodes: Set<EnvironmentComponentItem> = new Set<EnvironmentComponentItem>();

	const initZoomEventHandlers = (zoomInHandler: () => void, zoomOutHandler: () => void, resetHandler: () => void) => {
		zoomIn = zoomInHandler;
		zoomOut = zoomOutHandler;
		reset = resetHandler;
	};

	const routeToAppView = (name: string) => {
		const { protocol, host } = window.location;
		window.location.href = `${protocol}//${host}/applications/${name}/resource-view`;
	};

	// const addResourceNodesForApplication = (
	// 	resources: ResourceResult[],
	// 	parentRef: EnvironmentComponentItem,
	// 	generatedNodes: any
	// ) => {
	// 	generatedNodes.push(
	// 		...resources.map((r: ResourceResult) => {
	// 			return {
	// 				id: r.name,
	// 				name: r.name,
	// 				dependsOn: [parentRef.componentName],
	// 				syncFinishedAt: parentRef.syncFinishedAt,
	// 				shape: getShape(''),
	// 				icon: <AppIcon height={128} width={128} y="3" />,
	// 				syncStatus: r.status,
	// 				healthStatusIcon: null,
	// 				expandIcon:
	// 					r.kind === 'Application' ? (
	// 						<Expand
	// 							title="expand"
	// 							onClick={e => {
	// 								e.stopPropagation();
	// 								routeToAppView(r.name);
	// 							}}
	// 						/>
	// 					) : (
	// 						''
	// 					),
	// 				labels: {
	// 					component_type: 'argocd',
	// 				},
	// 				kind: r.kind,
	// 				onNodeClick,
	// 			};
	// 		})
	// 	);
	// 	return generatedNodes;
	// };

	const generateNodes = () => {
		const projectId = entityStore.getTeam((environmentItem as Environment).teamId)?.name;
		cleanDagNodeCache(environmentId);
		const generatedNodes: any[] = [
			{
				projectId,
				name: environmentId, // TODO check what name here goes from metadata
				id: 'root',
				shape: getShape('root'),
				icon: <LayersIcon />,
				syncStatus: environmentItem?.status,
				syncFinishedAt: environmentItem?.lastReconcileDatetime,
				componentStatus: environmentItem?.status,
				onNodeClick,
				estimatedCost: environmentItem?.estimatedCost,
			},
		];
		generatedNodes.push(
			...nodes.map((item: Component) => ({
				...item,
				onNodeClick,
				id: item.name,
				name: item.name,
				dependsOn: item.dependsOn?.length ? item.dependsOn : ['root'],
				icon: <ComputeIcon />,
					// item.labels?.component_type === 'argocd' ? (
						// <AppIcon height={128} width={128} y="4" />
					// ) : (
					// 	<ComputeIcon />
					// ),
				isSkipped: false,
				estimatedCost: item.estimatedCost,
				syncStatus: item.status || 'Unknown',
				componentStatus: item.status || 'Unknown',
				syncFinishedAt: item.lastReconcileDatetime,
				expandIcon: '',
			}))
		);

		// if (expandedNodes.size > 0) {
		// 	[...expandedNodes.values()].forEach(e =>
		// 		addResourceNodesForApplication(e.syncResult?.resources || [], e, generatedNodes)
		// 	);
		// }

		setData(generatedNodes);
	};

	useEffect(() => {
		if (!environmentItem || !nodes) return;
		generateNodes();
	}, [nodes, environmentItem]);

	// useEffect(() => {
	// 	if (!environmentItem) {
	// 		return;
	// 	}
	// 	const watcherSub = subscriberWatcher.subscribe(e => {
	// 		if (e?.application?.metadata?.name?.replace('-team-watcher', '') === environmentItem?.labels?.project_id) {
	// 			const status = e?.application?.status?.operationState?.phase;
	// 			setWatcherStatus(status);
	// 		}
	// 	});
	// 	if (environmentItem.conditions.length > 0) {
	// 		setEnvironmentCondition(getEnvironmentErrorCondition(environmentItem.conditions));
	// 	}
	// 	return () => watcherSub.unsubscribe();
	// }, [environmentItem]);

	useEffect(() => {
		if (!syncStarted) {
			return;
		}
		setTimeout(() => {
			setSyncStarted(false);
		}, 10000);
	}, [syncStarted]);

	const getSyncIconClass = (syncStatus: any) => {
		if (!syncStatus) {
			return;
		}
		if (syncStatus === ESyncStatus.OutOfSync) {
			return '--out-of-sync';
		} else if (syncStatus === ESyncStatus.Synced) {
			return '--in-sync';
		} else {
			return '--unknown';
		}
	};

	return (
		<div>
			<p className="node__title m-0 cursor-pointer">
				<span>* Costs are monthly estimates calculated at the time of last reconciliation</span>
				<div className="dag-controls">
					<button
						className="dag-controls-reconcile"
						onClick={async (e: any) => {
							e.stopPropagation();
							// if (environmentItem?.healthStatus !== 'Progressing')
								await syncMe(
									environmentItem as Environment,
									syncStarted,
									setSyncStarted,
									nm as NotificationsApi,
									watcherStatus
								);
						}}>
						<span
							className={`tooltip ${
								// environmentItem?.healthStatus !== 'Progressing' &&
								// !syncStarted &&
								// environmentCondition &&
								// 'error'
								''
							}`}>{`${
							// environmentItem?.healthStatus === 'Progressing' || syncStarted
							// 	? 'Reconciling...'
							// 	: environmentCondition || 'Reconcile Environment'
							'Reconcile Environment'
						}`}</span>
						<SyncIcon
							className={`large-health-icon-container__sync-button large-health-icon-container__sync-button${getSyncIconClass(
								// environmentItem?.syncStatus
								''
							)} large-health-icon-container__sync-button${
								// environmentItem?.healthStatus === 'Progressing' || syncStarted ? '--in-progress' : ''
								''
							}`}
							title="Reconcile Environment"
						/>
						Reconcile
					</button>
					<button className="dag-controls-zoom-button" onClick={() => zoomIn()}>
						<span className={`tooltip`}>Zoom In</span>+
					</button>
					<button className="dag-controls-zoom-button" onClick={() => zoomOut()}>
						-<span className={`tooltip`}>Zoom Out</span>
					</button>
					<div className={`tree-view-modifer ${modifierMizimized ? 'tree-view-modifer__minimized' : ''}`}>
						<div className="modifier">
							<label>Network Type:</label>
							<Select
								value={ranker}
								options={networkTypes}
								onChange={(res: { value: string }): void => {
									setRanker(res.value);
								}}
							/>
						</div>
						<div className="modifier">
							<label>Curve Type:</label>
							<Select
								value={option}
								options={curveTypes}
								onChange={(res: { value: string }): void => {
									setOption(res.value);
								}}
							/>
						</div>
						<div className="z-input modifier">
							<label>Parent-Child Spacing:</label>
							<input
								className="shadowy-input"
								type="number"
								value={rankSep}
								onChange={e => {
									if (e.target.value) {
										setRankSep(Number(e.target.value));
									} else {
										setRankSep(0);
									}
								}}
							/>
						</div>
						<div className="z-input modifier">
							<label>Sibling Spacing:</label>
							<input
								className="shadowy-input"
								type="number"
								value={nodeSep}
								onChange={e => {
									if (e.target.value) {
										setNodeSep(Number(e.target.value));
									} else {
										setNodeSep(0);
									}
								}}
							/>
						</div>
					</div>
					<button
						className="dag-controls-show-more"
						title="Show Modifiers"
						onClick={() => {
							toggleMinimize(!modifierMizimized);
						}}>
						<ArrowUp />
						<span className={`tooltip`}>More Options</span>
					</button>
				</div>
			</p>
			<Tree
				environmentId={environmentId}
				arrowType={option}
				data={data}
				nodeSep={nodeSep}
				rankSep={rankSep}
				ranker={ranker}
				onNodeClick={onNodeClick}
				zoomEventHandlers={initZoomEventHandlers}
				rankDir=""
				deploymentType=""
			/>
		</div>
	);
};

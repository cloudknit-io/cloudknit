import { ReactComponent as ArrowUp } from 'assets/images/icons/chevron-right.svg';
import { ReactComponent as ComputeIcon } from 'assets/images/icons/DAG-View/config.svg';
import { ReactComponent as LongArrow } from 'assets/images/icons/DAG/long-arrow-down.svg';
import { ReactComponent as LayersIcon } from 'assets/images/icons/DAG-View/environment-icon.svg';
import { ReactComponent as AppIcon } from 'assets/images/icons/DAG-View/Layers.svg';
import { getHealthStatusIcon, getSyncStatusIcon } from 'components/molecules/cards/renderFunctions';
import { useApi } from 'hooks/use-api/useApi';
import { EnvironmentComponentItem, EnvironmentItem } from 'models/projects.models';
import React, { FC, useEffect, useState } from 'react';
import { ArgoEnvironmentsService } from 'services/argo/ArgoEnvironments.service';
import { ReactComponent as SyncIcon } from 'assets/images/icons/sync-icon.svg';
import { ReactComponent as Expand } from 'assets/images/icons/expand.svg';

import Tree from './TreeView';
import { ESyncStatus, HealthStatuses, OperationPhase, ZSyncStatus } from 'models/argo.models';
import { getEnvironmentErrorCondition, syncMe } from 'pages/authorized/environments/helpers';
import { Context } from 'context/argo/ArgoUi';
import { subscriberResourceTree, subscriberWatcher } from 'utils/apiClient/EventClient';
import { breadcrumbObservable, pageHeaderObservable } from 'pages/authorized/contexts/EnvironmentHeaderContext';
import { LocalStorageKey } from 'models/localStorage';
import { ArgoComponentsService } from 'services/argo/ArgoComponents.service';
import { useParams } from 'react-router-dom';
import { ArgoStreamService } from 'services/argo/ArgoStream.service';
import ApiClient from 'utils/apiClient';
import { ZSidePanel } from 'components/molecules/side-panel/SidePanel';
import { ZLoaderCover } from 'components/atoms/loader/LoaderCover';
import { ConfigWorkflowViewApplication } from 'pages/authorized/environment-components/config-workflow-view/ConfigWorkflowViewApplication';

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
	environmentItem?: EnvironmentItem;
}

export const ComponentResourceTree: FC = () => {
	const { componentId } = useParams() as any;
	const [data, setData] = useState<any[]>([]);
	const [selectedNode, setSelectedNode] = useState<any>(null);
	const [showRightPanel, setShowRightPanel] = useState<boolean>(false);
	const [option, setOption] = useState('0');
	const [rankSep, setRankSep] = useState(70);
	const [nodeSep, setNodeSep] = useState(100);
	const [ranker, setRanker] = useState('network-simplex');
	let zoomIn: () => void = () => {};
	let zoomOut: () => void = () => {};
	let reset: () => void = () => {};
	const nm = React.useContext(Context)?.notifications;
	const syncEnvironment = useApi(ArgoEnvironmentsService.syncEnvironment);
	const deleteEnvironment = useApi(ArgoEnvironmentsService.deleteEnvironment);

	const initZoomEventHandlers = (zoomInHandler: () => void, zoomOutHandler: () => void, resetHandler: () => void) => {
		zoomIn = zoomInHandler;
		zoomOut = zoomOutHandler;
		reset = resetHandler;
	};

	const getResourceTreeData = () => {
		console.log(componentId);
		ArgoComponentsService.getApplicationResourceTree(componentId)
			.then(({ data }) => {
				const generatedNodes = data.nodes.map((e: any) => {
					return {
						id: e.uid,
						name: e.name,
						dependsOn: e.parentRefs ? e.parentRefs.map((p: any) => p.uid) : ['root'],
						syncFinishedAt: e.createdAt,
						shape: getShape(''),
						icon: <AppIcon height={128} width={128} y="3" />,
						syncStatusIcon: getSyncStatusIcon(ZSyncStatus.InSync),
						healthStatusIcon: e.health?.status ? getHealthStatusIcon(e.health?.status) : null,
						expandIcon: e.kind === 'Application' ? <Expand title="expand" /> : '',
						labels: {
							component_type: 'argocd',
						},
						kind: e.kind,
					};
				});

				console.log(data.nodes);
				generatedNodes.push({
					projectId: componentId,
					name: componentId, // TODO check what name here goes from metadata
					id: 'root',
					shape: getShape('root'),
					icon: <LayersIcon />,
					status: 'Synced',
					syncStatusIcon: getSyncStatusIcon(ZSyncStatus.InSync),
					healthStatusIcon: getHealthStatusIcon('Healthy'),
					syncFinishedAt: new Date(Date.now()).toISOString(),
					labels: {
						component_type: 'argocd',
					},
				});
				setData(generatedNodes);
				return generatedNodes;
			})
			.then(nodes => {
				ApiClient.get<any>(`/cd/api/v1/component/${componentId}`).then(({ data }) => {
					if (data) {
						updateStatus(data, nodes);
					}
				});
			});
	};

	const initializeResourceStream = () => {
		subscriberResourceTree.subscribe(appData => {
			if (appData && data.length > 0) {
				updateStatus(appData.application, data);
			}
		});
		ArgoStreamService.streamResourceTree(componentId);
	};

	const updateStatus = (appData: any, data: any) => {
		const root = data.find((e: any) => e.id === 'root');
		root.status = appData.status.sync.status;
		root.syncStatusIcon = getSyncStatusIcon(root.status);
		root.labels = appData.metadata.labels;
		appData.status.resources.forEach((r: any) => {
			const a = data.find((e: any) => e.name === r.name && e.kind === r.kind);
			if (a) {
				a.status = r.status;
				a.syncStatusIcon = getSyncStatusIcon(r.status);
			}
		});

		setData([...data]);
	};

	useEffect(() => {
		pageHeaderObservable.next({
			breadcrumbs: [],
			headerTabs: [],
			pageName: 'Resource View',
			filterTitle: '',
			onSearch: () => {},
			onViewChange: () => {},
			buttonText: '',
			checkBoxFilters: null,
		});
		breadcrumbObservable.next(false);
		getResourceTreeData();
	}, []);

	useEffect(() => {
		if (data?.length > 0 && subscriberResourceTree.observers?.length === 0) {
			initializeResourceStream();
			return () => subscriberResourceTree.unsubscribe();
		}
	}, [data]);

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
		<div style={{ position: 'relative' }}>
			<Tree
				environmentId={componentId}
				arrowType={option}
				data={data}
				nodeSep={nodeSep}
				rankSep={rankSep}
				ranker={ranker}
				rankDir={'LR'}
				onNodeClick={(configName: string, visualizationHandler?: Promise<any>) => {
					const node = data.find(e => e.id === configName);
					setSelectedNode({ ...node, id : node.name });
					setShowRightPanel(true);
				}}
				zoomEventHandlers={initZoomEventHandlers}
				deploymentType={''}
			/>
			{selectedNode && (
				<ZSidePanel isShown={showRightPanel} onClose={(): void => setShowRightPanel(false)}>
					<ZLoaderCover loading={false}>
						<ConfigWorkflowViewApplication
							key={selectedNode.id}
							projectId={data.find(e => e.id === 'root').labels.project_id}
							environmentId={data.find(e => e.id === 'root').labels.environment_id}
							config={selectedNode}
						/>
					</ZLoaderCover>
				</ZSidePanel>
			)}
		</div>
	);
};

import { getSyncStatusIcon } from 'components/molecules/cards/renderFunctions';
import { EnvironmentItem } from 'models/projects.models';
import React, { FC, useEffect, useState } from 'react';

import { ZLoaderCover } from 'components/atoms/loader/LoaderCover';
import { ZSidePanel } from 'components/molecules/side-panel/SidePanel';
import { Context } from 'context/argo/ArgoUi';
import { ESyncStatus } from 'models/argo.models';
import { breadcrumbObservable, pageHeaderObservable } from 'pages/authorized/contexts/EnvironmentHeaderContext';
import { useParams } from 'react-router-dom';
import { subscriberResourceTree } from 'utils/apiClient/EventClient';
import Tree from './TreeView';

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

	const initZoomEventHandlers = (zoomInHandler: () => void, zoomOutHandler: () => void, resetHandler: () => void) => {
		zoomIn = zoomInHandler;
		zoomOut = zoomOutHandler;
		reset = resetHandler;
	};

	const getResourceTreeData = () => {
		console.log(componentId);
	};

	const initializeResourceStream = () => {
		subscriberResourceTree.subscribe(appData => {
			if (appData && data.length > 0) {
				updateStatus(appData.application, data);
			}
		});
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
					setSelectedNode({ ...node, id: node.name });
					setShowRightPanel(true);
				}}
				zoomEventHandlers={initZoomEventHandlers}
				deploymentType={''}
			/>
			{selectedNode && (
				<ZSidePanel isShown={showRightPanel} onClose={(): void => setShowRightPanel(false)}>
					<ZLoaderCover loading={false}></ZLoaderCover>
				</ZSidePanel>
			)}
		</div>
	);
};

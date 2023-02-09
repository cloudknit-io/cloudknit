import { ReactComponent as ArrowUp } from 'assets/images/icons/chevron-right.svg';
import { ReactComponent as ComputeIcon } from 'assets/images/icons/DAG-View/config.svg';
import { ReactComponent as LayersIcon } from 'assets/images/icons/DAG-View/environment-icon.svg';
import { Select } from 'components/argo-core';
import React, { FC, useCallback, useEffect, useMemo, useState } from 'react';

import { AuditStatus } from 'models/argo.models';
import { EntityStore } from 'models/entity.store';
import { TreeReconcile } from 'pages/authorized/environments/helpers';
import { Reconciler } from 'pages/authorized/environments/Reconciler';
import { cleanDagNodeCache } from './node-figure-helper';
import Tree from './TreeView';
import { Component, Environment } from 'models/entity.type';

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
	let zoomIn: () => void = () => {};
	let zoomOut: () => void = () => {};
	let reset: () => void = () => {};
	const initZoomEventHandlers = (zoomInHandler: () => void, zoomOutHandler: () => void, resetHandler: () => void) => {
		zoomIn = zoomInHandler;
		zoomOut = zoomOutHandler;
		reset = resetHandler;
	};
	const renderTree = useCallback(() => {
		if (data.length === 0) return <></>;
		return (
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
		);
	}, [data, option, nodeSep, rankSep, ranker, onNodeClick]);
	

	const generateNodes = () => {
		const projectId = entityStore.getTeam((environmentItem as Environment).teamId)?.name;
		cleanDagNodeCache(environmentId);
		const generatedNodes: any[] = [
			{
				projectId,
				name: environmentId,
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
				isSkipped: [AuditStatus.SkippedProvision, AuditStatus.SkippedDestroy].includes(item.lastAuditStatus),
				estimatedCost: item.estimatedCost,
				syncStatus: item.status || 'Unknown',
				componentStatus: item.status || 'Unknown',
				syncFinishedAt: item.lastReconcileDatetime,
				expandIcon: '',
			}))
		);

		setData(generatedNodes);
	};

	useEffect(() => {
		if (!environmentItem || !nodes) return;
		generateNodes();
	}, [nodes, environmentItem]);

	useEffect(() => {
		if (!syncStarted) {
			return;
		}
		setTimeout(() => {
			setSyncStarted(false);
		}, 10000);
	}, [syncStarted]);

	return (
		<div>
			<p className="node__title m-0 cursor-pointer">
				<span>* Costs are monthly estimates calculated at the time of last reconciliation</span>
				<div className="dag-controls">
					{environmentItem && <Reconciler environment={environmentItem} template={TreeReconcile} />}
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
			{renderTree()}
		</div>
	);
};

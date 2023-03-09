import { AuditStatus, ESyncStatus, ZSyncStatus } from 'models/argo.models';
import { Component, DAG, Environment } from 'models/entity.type';
import { DagNode, DagProps } from './DagNode';
import { ReactComponent as ComputeIcon } from 'assets/images/icons/DAG-View/config.svg';
import { ReactComponent as LayersIcon } from 'assets/images/icons/DAG-View/environment-icon.svg';
import dagre from 'dagre';
import { ConnectionLineType, Edge, Handle, MarkerType, Node, Position } from 'reactflow';
import { EntityStore } from 'models/entity.store';

export const generateRootNode = (environment: Environment) => {
	const data = {
		icon: <LayersIcon />,
		...environment,
	};
	return (
		<>
			<Handle
				id="a"
				className="targetHandle"
				style={{ zIndex: 2, top: 0 }}
				position={Position.Top}
				type="source"
				isConnectable={true}
			/>
			<Handle
				id="b"
				className="targetHandle"
				style={{
					top: 0,
				}}
				position={Position.Top}
				type="target"
				isConnectable={true}
			/>
			<DagNode
				data={{
					cost: data.estimatedCost,
					name: data.name,
					icon: data.icon,
					status: data.status as ZSyncStatus,
					timestamp: data.lastReconcileDatetime,
					operation: 'Provision',
					isSkipped: false,
				}}
			/>
		</>
	);
};

export const generateComponentNode = (component: Component) => {
	const data = {
		...component,
		icon: <ComputeIcon />,
		isSkipped: [AuditStatus.SkippedProvision, AuditStatus.SkippedDestroy].includes(component.lastAuditStatus),
	};

	return (
		<>
			<Handle
				id="a"
				className="targetHandle"
				style={{ zIndex: 2, top: 0 }}
				position={Position.Top}
				type="source"
				isConnectable={true}
			/>
			<Handle
				id="b"
				className="targetHandle"
				style={{ top: 0 }}
				position={Position.Top}
				type="target"
				isConnectable={true}
			/>
			<DagNode
				data={{
					cost: data.estimatedCost,
					name: data.name,
					icon: data.icon,
					status: data.status as ZSyncStatus,
					timestamp: data.lastReconcileDatetime,
					operation: data.isDestroyed ? 'Destroy' : 'Provision',
					isSkipped: data.isSkipped,
				}}
			/>
		</>
	);
};

export const initializeLayout = () => {
	const dagreGraph = new dagre.graphlib.Graph();
	dagreGraph.setDefaultEdgeLabel(() => ({}));

	const nodeWidth = 250;
	const nodeHeight = 60;

	const getLayoutedElements = (nodes: any, edges: any, direction = 'TB') => {
		const isHorizontal = direction === 'LR';
		dagreGraph.setGraph({ rankdir: direction });

		nodes.forEach((node: any) => {
			dagreGraph.setNode(node.id, { width: nodeWidth, height: nodeHeight });
		});
		edges.forEach((edge: any) => {
			dagreGraph.setEdge(edge.source, edge.target);
		});

		dagre.layout(dagreGraph);

		nodes.forEach((node: any) => {
			const nodeWithPosition = dagreGraph.node(node.id);
			node.targetPosition = isHorizontal ? 'left' : 'top';
			node.sourcePosition = isHorizontal ? 'right' : 'bottom';

			// We are shifting the dagre node position (anchor=center center) to the top left
			// so it matches the React Flow node anchor point (top left).
			node.position = {
				x: nodeWithPosition.x - nodeWidth / 2,
				y: (nodeWithPosition.y - nodeHeight / 2) * 1.2,
			};

			return node;
		});

		return { nodes, edges };
	};

	return {
		getLayoutedElements,
	};
};

const getNode = (id: string, label: JSX.Element): Node => {
	return {
		id,
		data: { label },
		position: { x: 0, y: 0 },
		style: { padding: 0, border: 0 },
		draggable: false,
		selectable: false,
	};
};

const getEdge = (id: string, source: string, target: string): Edge => {
	return {
		id,
		source,
		target,
		type: ConnectionLineType.SimpleBezier,
		markerEnd: { type: MarkerType.ArrowClosed, width: 15, height: 15, color: '#333' },
		style: {
			strokeWidth: 1,
			stroke: '#333',
		},
		sourceHandle: 'a',
		targetHandle: 'b',
	};
};

export const generateNodesAndEdges = (environment: Environment) => {
	const nodes: Node[] = [];
	const edges: Edge[] = [];

	const envNode = getNode(environment.argoId.toString(), generateRootNode(environment));
	nodes.push(envNode);

	const components = EntityStore.getInstance().getComponentsByEnvId(environment.id);

	for (let e of environment.dag) {
		const item = components.find(c => c.name === e.name) as Component;
		if (!item) continue;

		nodes.push(getNode(item.argoId.toString(), generateComponentNode(item)));

		if (!e.dependsOn?.length) {
			edges.push(
				getEdge(`e${environment.argoId}-${item.argoId}`, environment.argoId.toString(), item.argoId.toString())
			);
		} else {
			e.dependsOn.forEach(d => {
				const dc = components.find(ic => ic.name === d);
				if (dc) {
					edges.push(getEdge(`e${dc.argoId}-${item.argoId}`, dc.argoId.toString(), item.argoId.toString()));
				}
			});
		}
	}

	return {
		nodes,
		edges,
	};
};

export const colorLegend = [
	{
		key: 'Succeeded',
		value: '#D3E4CD',
		order: 5,
	},
	{
		key: 'Failed',
		value: '#FDE9F2',
		order: 6,
	},
	{
		key: 'Waiting for approval',
		value: '#FFC85C',
		order: 2,
	},
	{
		key: 'In Progress',
		value: 'linear-gradient(to right,rgba(226, 194, 185, 0.5),rgba(226, 194, 185, 1),rgb(190, 148, 137))',
		order: 0,
	},
	{
		key: 'Destroyed/Not Provisioned',
		value: '#ddd',
		order: 5,
	},
];

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

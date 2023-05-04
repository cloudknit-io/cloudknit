import { ReactComponent as ComputeIcon } from 'assets/images/icons/DAG-View/config.svg';
import { ReactComponent as LayersIcon } from 'assets/images/icons/DAG-View/environment-icon.svg';
import dagre from 'dagre';
import { AuditStatus, ESyncStatus, ZSyncStatus } from 'models/argo.models';
import { EntityStore } from 'models/entity.store';
import { Component, Environment } from 'models/entity.type';
import { useCallback } from 'react';
import { Edge, Handle, MarkerType, Node, Position, useStore } from 'reactflow';
import { DagNode } from './DagNode';

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
				style={{ zIndex: 2, top: 30 }}
				position={Position.Bottom}
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
				style={{ zIndex: 2, top: 30 }}
				position={Position.Bottom}
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
		type: 'smart',
		markerEnd: { type: MarkerType.ArrowClosed, width: 15, height: 15, color: '#333' },
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
		case ZSyncStatus.ValidationFailed:
		case ZSyncStatus.ApplyFailed:
		case ZSyncStatus.ProvisionFailed:
		case ZSyncStatus.DestroyFailed:
		case AuditStatus.DestroyApplyFailed:
		case AuditStatus.DestroyPlanFailed:
		case AuditStatus.ProvisionApplyFailed:
		case AuditStatus.ProvisionPlanFailed:
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

function getNodeIntersection(intersectionNode: any, targetNode: any) {
	// https://math.stackexchange.com/questions/1724792/an-algorithm-for-finding-the-intersection-point-between-a-center-of-vision-and-a
	const {
		width: intersectionNodeWidth,
		height: intersectionNodeHeight,
		positionAbsolute: intersectionNodePosition,
	} = intersectionNode;
	const targetPosition = targetNode.positionAbsolute;

	const w = intersectionNodeWidth / 2;
	const h = intersectionNodeHeight / 2;

	const x2 = intersectionNodePosition.x + w;
	const y2 = intersectionNodePosition.y + h;
	const x1 = targetPosition.x + w;
	const y1 = targetPosition.y + h;

	const xx1 = (x1 - x2) / (2 * w) - (y1 - y2) / (2 * h);
	const yy1 = (x1 - x2) / (2 * w) + (y1 - y2) / (2 * h);
	const a = 1 / (Math.abs(xx1) + Math.abs(yy1));
	const xx3 = a * xx1;
	const yy3 = a * yy1;
	const x = w * (xx3 + yy3) + x2;
	const y = h * (-xx3 + yy3) + y2;

	return { x, y };
}

// returns the position (top,right,bottom or right) passed node compared to the intersection point
function getEdgePosition(node: any, intersectionPoint: any) {
	const n = { ...node.positionAbsolute, ...node };
	const nx = Math.round(n.x);
	const ny = Math.round(n.y);
	const px = Math.round(intersectionPoint.x);
	const py = Math.round(intersectionPoint.y);

	if (px <= nx + 1) {
		return Position.Left;
	}
	if (px >= nx + n.width - 1) {
		return Position.Right;
	}
	if (py <= ny + 1) {
		return Position.Top;
	}
	if (py >= n.y + n.height - 1) {
		return Position.Bottom;
	}

	return Position.Top;
}

// returns the parameters (sx, sy, tx, ty, sourcePos, targetPos) you need to create an edge
export function getEdgeParams(source: any, target: any) {
	const sourceIntersectionPoint = getNodeIntersection(source, target);
	const targetIntersectionPoint = getNodeIntersection(target, source);

	const sourcePos = getEdgePosition(source, sourceIntersectionPoint);
	const targetPos = getEdgePosition(target, targetIntersectionPoint);

	return {
		sx: sourceIntersectionPoint.x,
		sy: sourceIntersectionPoint.y,
		tx: targetIntersectionPoint.x,
		ty: targetIntersectionPoint.y,
		sourcePos,
		targetPos,
	};
}

export function createNodesAndEdges() {
	const nodes = [];
	const edges = [];
	const center = { x: window.innerWidth / 2, y: window.innerHeight / 2 };

	nodes.push({ id: 'target', data: { label: 'Target' }, position: center });

	for (let i = 0; i < 8; i++) {
		const degrees = i * (360 / 8);
		const radians = degrees * (Math.PI / 180);
		const x = 250 * Math.cos(radians) + center.x;
		const y = 250 * Math.sin(radians) + center.y;

		nodes.push({ id: `${i}`, data: { label: 'Source' }, position: { x, y } });

		edges.push({
			id: `edge-${i}`,
			target: 'target',
			source: `${i}`,
			type: 'floating',
			markerEnd: {
				type: MarkerType.Arrow,
			},
		});
	}

	return { nodes, edges };
}

export function FloatingEdge({ id, source, target, markerEnd, style }: any) {
	const sourceNode = useStore(useCallback(store => store.nodeInternals.get(source), [source]));
	const targetNode = useStore(useCallback(store => store.nodeInternals.get(target), [target]));

	if (!sourceNode || !targetNode) {
		return null;
	}

	const { sx, sy, tx, ty, sourcePos, targetPos } = getEdgeParams(sourceNode, targetNode);

	const gB = ({
		sourceX,
		sourceY,
		sourcePosition = Position.Bottom,
		targetX,
		targetY,
		targetPosition = Position.Top,
		curvature = 0.25,
	}: any) => {
		const [sourceControlX, sourceControlY]: any = getControlWithCurvature({
			pos: sourcePosition,
			x1: sourceX,
			y1: sourceY,
			x2: targetX,
			y2: targetY,
			c: curvature,
		});

		return [`M${sourceX},${sourceY} Q${sourceControlX},${sourceControlY} ${targetX},${targetY}`];
	};

	const [edgePath] = gB({
		sourceX: sx,
		sourceY: sy,
		sourcePosition: sourcePos,
		targetPosition: targetPos,
		targetX: tx,
		targetY: ty,
		curvature: 0,
	});

	return <path id={id} className="react-flow__edge-path" d={edgePath} markerEnd={markerEnd} style={style} />;
}

function getControlWithCurvature({ pos, x1, y1, x2, y2, c }: any) {
	switch (pos) {
		case Position.Left:
			return [x1 - calculateControlOffset(x1 - x2, c), y1];
		case Position.Right:
			return [x1 + calculateControlOffset(x2 - x1, c), y1];
		case Position.Top:
			return [x1, y1 - calculateControlOffset(y1 - y2, c)];
		case Position.Bottom:
			return [x2, y1 + calculateControlOffset(y2 - y1, c)];
	}
}

function calculateControlOffset(distance: any, curvature: any) {
	if (distance >= 0) {
		return 0.5 * distance;
	}
	return curvature * 25 * Math.sqrt(-distance);
}

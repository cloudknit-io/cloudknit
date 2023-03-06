import { EntityStore } from 'models/entity.store';
import { Component, Environment } from 'models/entity.type';
import React, { useCallback, useEffect, useMemo, useState } from 'react';
import ReactFlow, { addEdge, ConnectionLineType, useEdgesState, useNodesState, useReactFlow } from 'reactflow';
import 'reactflow/dist/style.css';
import './tree-view-new.scss';
import { generateNodesAndEdges, initializeLayout } from './tree-view.helper';
import { TreeViewControls } from './TreeViewControls';

interface Props {
	onNodeClick: any;
	environmentItem?: Environment;
}

export const TreeView: React.FC<Props> = ({ onNodeClick, environmentItem }) => {
	const [dagNodes, setDagNodes, onNodesChange] = useNodesState([]);
	const [dagEdges, setDagEdges, onEdgesChange] = useEdgesState([]);
	const [render, setRender] = useState<Symbol | null>(null);
	
	const getLayoutedElements = useMemo(() => {
		const { getLayoutedElements } = initializeLayout();
		return getLayoutedElements;
	}, []);

	useEffect(() => {
		if (!environmentItem) return;
		const sub = EntityStore.getInstance().emitterComp.subscribe((components: Component[]) => {
			if (components.length === 0 || components[0].envId !== environmentItem.id) return;
			setRender(Symbol('render'));
		});
		setRender(Symbol('render'));
		setTimeout(() => {
			const vp = document.querySelector('.react-flow__viewport') as HTMLElement;
			if (vp) {
				const tf = window.getComputedStyle(vp).transform;
				if (tf) {
					const matrix = tf.split(',');
					matrix[matrix.length - 1] = '0)';
					vp.style.transform = matrix.join(',');
					vp.style.opacity = '1';
				}
			}
		}, 300);
		return () => sub.unsubscribe();
	}, [environmentItem]);

	useEffect(() => {
		if (!render || !environmentItem) return;
		const { nodes, edges } = generateNodesAndEdges(environmentItem);
		const { nodes: layoutedNodes, edges: layoutedEdges } = getLayoutedElements(nodes, edges);
		setDagNodes(layoutedNodes);
		setDagEdges(layoutedEdges);
	}, [render]);

	const onConnect = useCallback(
		params => setDagEdges(eds => addEdge({ ...params, type: ConnectionLineType.SmoothStep, animated: true }, eds)),
		[]
	);

	const onLayout = useCallback(
		direction => {
			const { nodes: layoutedNodes, edges: layoutedEdges } = getLayoutedElements(dagNodes, dagEdges, direction);

			setDagNodes([...layoutedNodes]);
			setDagEdges([...layoutedEdges]);
		},
		[dagNodes, dagEdges]
	);

	return (
		<>
			<TreeViewControls environment={environmentItem} />
			<div
				className="layoutflow"
				style={{ height: '80vh' }}>
				<ReactFlow
					nodes={dagNodes}
					edges={dagEdges}
					onNodesChange={onNodesChange}
					onEdgesChange={onEdgesChange}
					onConnect={onConnect}
					connectionLineType={ConnectionLineType.SmoothStep}
					onNodeClick={e => onNodeClick(e.currentTarget.getAttribute('data-id'))}
					maxZoom={1}
					fitView
					
				/>
				{/* <div className="controls">
				<button onClick={() => onLayout('TB')}>vertical layout</button>
				<button onClick={() => onLayout('LR')}>horizontal layout</button>
			</div> */}
			</div>
		</>
	);
	// };
};

import './style.scss';

import { Edge } from 'components/organisms/dag-view/Edge';
import { Node } from 'components/organisms/dag-view/Node';
import * as d3 from 'd3';
import { hierarchy, HierarchyPointLink, HierarchyPointNode } from 'd3';
import { DagNode } from 'models/dag.models';
import React, { FC, useEffect, useRef, useState } from 'react';

type Props = {
	data: DagNode;
	onNodeClick: Function;
};

export const DagView: FC<Props> = ({ data, onNodeClick }: Props) => {
	const [nodes, setNodes] = useState<HierarchyPointNode<DagNode>[]>();
	const [links, setLinks] = useState<HierarchyPointLink<DagNode>[]>();
	const [containerWidth, setContainerWidth] = useState<number>(0);
	const [containerHeight, setContainerHeight] = useState<number>(100);
	const containerRef = useRef<HTMLDivElement>(null);

	const generateTree = (): { nodes: HierarchyPointNode<DagNode>[]; links: HierarchyPointLink<DagNode>[] } => {
		const tree: d3.TreeLayout<DagNode> = d3
			.tree<DagNode>()
			.nodeSize([40, 150])
			.separation(() => 6);
		const rootNode = tree(hierarchy(data, d => d?.children || []));
		const nodes = rootNode.descendants();
		const links = rootNode.links();

		return { nodes, links };
	};

	useEffect(() => {
		const data = generateTree();
		setNodes(data.nodes);
		setLinks(data.links);
	}, [data]);

	useEffect(() => {
		if (containerRef.current) {
			setContainerWidth(containerRef.current.getBoundingClientRect().width);
			const gElem = containerRef.current.querySelector('g');

			setImmediate(() => {
				gElem && setContainerHeight(gElem.getBoundingClientRect().height);
			});
		}
	}, [containerRef]);

	return (
		<div className={`rd3t-container`} ref={containerRef}>
			<svg
				className={`rd3t-svg`}
				height={containerHeight + 60}
				width={'100%'}
				onClick={() => console.log('loaded')}>
				<g transform={`translate(${containerWidth / 2}, 60)`}>
					{links &&
						links.map((linkData: HierarchyPointLink<DagNode>, i: number) => (
							<Edge key={'edge-' + i} linkData={linkData} />
						))}
					{nodes &&
						nodes.map(({ data, x, y, parent }: HierarchyPointNode<DagNode>, i: number) => {
							return (
								<Node
									key={'node-' + i}
									nodeDatum={data}
									position={{ x, y }}
									parent={parent}
									toggleNode={() => console.log('TOGGLE NODE')}
									onNodeClick={() => onNodeClick(data.name)}
									onNodeMouseOver={() => console.log('MOUSE OVER')}
									onNodeMouseOut={() => console.log('MOUSE OUT')}
								/>
							);
						})}
				</g>
			</svg>
		</div>
	);
};

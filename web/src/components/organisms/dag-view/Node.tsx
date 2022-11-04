import { HierarchyPointNode } from 'd3';
import { select } from 'd3-selection';
import { ZSyncStatus } from 'models/argo.models';
import { DagNode, Point } from 'models/dag.models';
import React, { FC, useEffect, useState } from 'react';

type Props = {
	nodeDatum: DagNode;
	position: Point;
	parent: HierarchyPointNode<DagNode> | null;
	toggleNode: Function;
	onNodeClick: Function;
	onNodeMouseOver: Function;
	onNodeMouseOut: Function;
};

const textLayout = {
	title: {
		textAnchor: 'middle',
		x: 0,
		y: 50,
	},
};

export const Node: FC<Props> = ({ nodeDatum, onNodeClick, onNodeMouseOver, onNodeMouseOut, position }: Props) => {
	const [nodeRef, setNodeRef] = useState<any>();

	useEffect(() => {
		select(nodeRef).attr('transform', `translate(${position.x},${position.y})`).style('opacity', 1);
	});

	const getClassName = (status: any): string => {
		switch (status) {
			case ZSyncStatus.InSync:
				return '--successful';
			case ZSyncStatus.Provisioned:
				return '--successful';
			case ZSyncStatus.OutOfSync:
				return '--failed';
			case ZSyncStatus.Initializing:
				return '--initializing';
			case ZSyncStatus.RunningPlan:
				return '--pending';
			default:
				return '--unknown';
		}
	};

	return (
		<g
			ref={n => {
				setNodeRef(n);
			}}
			transform={`translate(${position.x},${position.y})`}
			className="rd3t-node">
			<circle
				r={20}
				className={'rd3t-node__pod rd3t-node__pod' + getClassName(nodeDatum.status)}
				onMouseOver={(): void => onNodeMouseOver()}
				onMouseOut={(): void => onNodeMouseOut()}
				onClick={(): void => onNodeClick(nodeDatum.name)}
				z="1"
			/>
			<rect
				className={'rd3t-node__outline' + getClassName(nodeDatum.status)}
				x="-30"
				y="-30"
				width="60"
				height="60"
				rx="30"
				z="0"
				fill="none"
			/>
			<g className="rd3t-node__label">
				<text
					className="rd3t-node__label__title"
					{...textLayout.title}
					onClick={(): void => onNodeClick(nodeDatum.name)}>
					{nodeDatum.name}
				</text>
			</g>
		</g>
	);
};

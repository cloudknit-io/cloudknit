import './style.scss';

import { NodeIcon, NodeStatus, ZDiagramNode } from 'components/molecules/diagram-node/DiagramNode';
import { ZEditor } from 'components/molecules/editor/Editor';
import { ZAsyncRenderer } from 'components/molecules/zasync-renderer/ZAsyncRenderer';
import React, { FC } from 'react';
import { useEffect } from 'react';
import { useState } from 'react';

import { ZAccordion } from '../accordion/ZAccordion';
import { ZStreamRenderer } from 'components/molecules/zasync-renderer/ZStreamRenderer';
import { ZSyncStatus } from 'models/argo.models';
type Props = {
	approvedBy: string;
	nodes: { [name: string]: any };
};
const nodesOrder = [
	'provision plan',
	'destroy plan',
	'plan',
	'notify',
	'approve',
	'apply',
	'provision apply',
	'destroy apply',
];

export const ZWorkflowDiagram: FC<Props> = ({ nodes, approvedBy }: Props) => {
	const [nodeArray, setNodes] = useState<any>([]);
	const [mutated, setMutated] = useState<any>(null);
	const [summary, setSummary] = useState<Map<string, string>>(new Map());

	const getNodeStatus = (node: any): NodeStatus => {
		switch (node.phase) {
			case 'Failed':
				return 'Failed';
			case 'Error':
				return 'Failed';
			case 'Running':
				return node.displayName === 'apply' ? 'InProcess' : 'Pending';
			case 'Succeeded':
				return 'Succeeded' && mutated === node ? 'Mutated' : 'Succeeded';
			case 'Terminating':
				return 'Failed';
			default:
				return 'Disregarded';
		}
	};

	const getNodeIcon = (node: any): NodeIcon => {
		switch (node.phase) {
			case 'Failed':
				return 'Failed';
			case 'Error':
				return 'Failed';
			case 'Running':
				return node.displayName === 'apply' ? 'InProcess' : 'Pending';
			case 'Succeeded':
				return 'Synced' && mutated === node ? 'Mutated' : 'Synced';
			case 'Terminating':
				return 'Failed';
			case 'Skipped':
				return 'Skipped';
			default:
				return 'Failed';
		}
	};

	const accordionHeader = (node: any, name: string) => {
		return (
			<div
				className={`workflow-accordion-header workflow-accordion-header_phase--${getNodeStatus(
					node
				)} workflow-accordion-header workflow-accordion-header_phase--${
					node.displayName.toLowerCase().includes('teardown') ||
					node.displayName.toLowerCase().includes('skipped')
						? 'destroy'
						: ''
				}`}>
				<ZDiagramNode text="" icon={getNodeIcon(node)} status={getNodeStatus(node)} />
				{name} <SmallText data={node.displayName === 'approve' ? approvedBy : summary.get(node.displayName)} />
			</div>
		);
	};

	const accordionContent = (content: any) => {
		return <div className="workflow-accordion-content zscrollbar">{content}</div>;
	};

	const initSummary = async (node: any, data: any) => {
		try {
			if (!data) {
				return;
			}
			let sum =
				data.match(/^Plan:.*$/gm) ||
				data.match(/^Apply complete!.*$/gm) ||
				data.match(/^No changes\. Infrastructure is up-to-date\.$/gm);
			if (!sum) {
				return '';
			}
			setMutated(sum[0].includes('Plan:') ? node : null);
			summary.set(node.displayName, sum[0]);
			setSummary(new Map([...summary.entries()]));
		} catch (err) {
			return;
		}
	};

	const getContent = (node: any) => {
		if (node.phase === 'Succeeded' || node.phase === 'Failed') {
			return (
				<ZAsyncRenderer
					key={node.name}
					promise={node.getS3Logs().then((data: any) => {
						initSummary(node, data);
						return data;
					})}
					Component={ZEditor}
					componentProps={{
						readOnly: true
					}}
					defaultValue={'Loading info...'}
				/>
			);
		} else {
			return <ZStreamRenderer key={node.name} subject={node.getNodeLogs().listen()} Component={ZEditor} componentProps={{
				readOnly: true
			}} />;
		}
	};

	const resetAccordionHeader = (node: any) => {
		setNodes(
			nodeArray.map((n: any) => {
				if (n.nodeData === node) {
					return {
						...n,
						accordionHeader: accordionHeader(node, node.displayName),
					};
				}
				return n;
			})
		);
	};

	useEffect(() => {
		if (!mutated || !summary) return;
		resetAccordionHeader(mutated);
	}, [mutated]);

	useEffect(() => {
		if (summary.size === 0) {
			return;
		}
		summary.forEach((v: string, k: string) => resetAccordionHeader(nodes.find((n: any) => n.displayName === k)));
	}, [summary]);

	useEffect(() => {
		setSummary(new Map());
		setMutated(null);
		if (
			nodes.some(
				(e: any) => e.configStatus === ZSyncStatus.SkippedReconcile || e.configStatus === ZSyncStatus.Skipped
			)
		) {
			const displayName =
				nodes[0].configStatus === ZSyncStatus.Skipped ? 'Skipped Teardown' : 'Skipped Reconcile';
			setNodes([
				{
					nodeData: {},
					accordionHeader: accordionHeader(
						{
							phase: 'Skipped',
							displayName,
						},
						displayName
					),
					accordionContent: accordionContent(''),
				},
			]);
		} else {
			setNodes(
				nodes
					.sort(function (a: any, b: any) {
						return nodesOrder.indexOf(a.displayName) - nodesOrder.indexOf(b.displayName);
					})
					.map((node: any, index: number) => {
						return {
							nodeData: node,
							accordionHeader: accordionHeader(node, node.displayName),
							accordionContent: accordionContent(
								node.displayName === 'approve' ? node.getZFeedbackModal() : getContent(node)
							),
							collapsed: index < nodes.length - 1
						};
					})
			);
		}
	}, [nodes]);

	return (
		<div className="zlifecycle-workflow-diagram">
			<ZAccordion items={nodeArray} />
		</div>
	);
};

interface SmallProps {
	data: any;
}
export const SmallText: React.FC<SmallProps> = ({ data }: SmallProps) => {
	return <small style={{ marginLeft: '20px', color: 'gray', textTransform: 'none' }}>{data}</small>;
};

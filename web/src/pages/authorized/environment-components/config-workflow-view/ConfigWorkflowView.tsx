import './style.scss';

import { ZEditor } from 'components/molecules/editor/Editor';
import { AuditView } from 'components/organisms/audit_view/AuditView';
import { HierarchicalView } from 'components/organisms/hierarchical-view/hierarchical-view';
import { ZWorkflowDiagram } from 'components/organisms/workflow-diagram/WorkflowDiagram';
import { ZSyncStatus } from 'models/argo.models';
import { ZFeedbackModal } from 'pages/authorized/environment-components/config-workflow-view/FeedbackModal';
import React, { FC, useEffect, useState } from 'react';
import { AuditService } from 'services/audit/audit.service';

import { StateFileView } from 'components/organisms/state-file-view/StateFileView';
import { useRef } from 'react';
import { EventClientLogs } from 'utils/apiClient/EventClient';
import { auditColumns, getSeparatedConfigId, ViewType, ViewTypeTabName } from '../helpers';
import { ConfigWorkflowLeftView } from './ConfigWorkflowLeftView';
import { CompAuditData, Component } from 'models/entity.type';

type Props = {
	projectId: string;
	environmentId: string;
	config: Component;
	workflowData: any;
	logs: string | null;
	plans: string | null;
	auditData: CompAuditData[];
};

export type WorkflowNode = {
	configStatus: string;
	getZFeedbackModal: Function;
	getS3Logs: Function;
	getNodeLogs: Function;
	displayName: string;
	name: string;
	phase: string;
};

export const ConfigWorkflowView: FC<Props> = (props: Props) => {
	const { projectId, environmentId, config, logs, plans, workflowData, auditData } = props;
	const clientLogMap = new Map<string, EventClientLogs>();
	const nodesRef = useRef<Map<string, EventClientLogs>>(new Map<string, EventClientLogs>());
	const [filteredNodes, setFilteredNodes] = useState<WorkflowNode[]>([]);
	const [viewType, setViewType] = useState<any>(ViewType.Concise_Logs);
	const [ilRepo, setIlRepo] = useState<string>('');
	const separatedConfigId = config ? getSeparatedConfigId(config) : null;

	useEffect(() => {
		if (!workflowData) {
			setFilteredNodes([]);
			return;
		}
		const newNodes: any[] = Object.values(workflowData?.status?.nodes || {}).filter(
			(node: any) => node.type === 'Steps' || node.type === 'Pod' || node.type === 'Suspend'
		);
		const planNode = newNodes.find(e => e.displayName === 'plan');
		const teardown = planNode?.inputs?.parameters?.find((param: any) => param.name === 'is_destroy')?.value;
		const ilRepo = workflowData?.spec?.arguments?.parameters?.find((param: any) => param.name === 'il_repo')?.value;
		setIlRepo(ilRepo);
		const foundApprovedNode = newNodes.find((node: any) => node.displayName === 'approve') as any;
		const needsApproval = foundApprovedNode?.phase === 'Running' || false;
		const isApproved = foundApprovedNode?.phase !== 'Running' || false;
		setFilteredNodes(
			newNodes
				.filter((node: any) => node.displayName !== 'notify' && node.type !== 'Steps')
				.map((node: any) => {
					const podName = node.boundaryID + '-run' + node.id.replace(node.boundaryID, '');
					const status = getComponentStatus(config);
					if ((node.displayName === 'apply' || node.displayName === 'plan') && status) {
						node.displayName = `${status} ${node.displayName}`;
					}

					const clientLogUrl = `/wf/api/v1/stream/projects/${projectId}/environments/${environmentId}/config/${config.argoId}/${config.lastWorkflowRunId}/log/${podName}`;
					if (
						!clientLogMap.has(clientLogUrl) &&
						node.phase !== 'Succeeded' &&
						node.displayName !== 'approve' &&
						node.phase !== 'Failed'
					) {
						nodesRef.current.set(clientLogUrl, new EventClientLogs(clientLogUrl));
					}

					return {
						phase: node.phase,
						name: node.name,
						displayName: node.displayName,
						configStatus: config.status,
						getZFeedbackModal: () => getZFeedbackModal(isApproved, needsApproval),
						getS3Logs: async () =>
							AuditService.getInstance()
								[node.displayName.includes('apply') ? 'fetchApplyLogs' : 'fetchPlanLogs'](
									separatedConfigId?.team || '',
									separatedConfigId?.environment || '',
									separatedConfigId?.component || '',
									0,
									true
								)
								.then(({ data }) => {
									if (Array.isArray(data) && data.length > 0) {
										return data[0].body;
									}
									return data;
								}),
						getNodeLogs: () => nodesRef.current.get(clientLogUrl),
					};
				})
		);

		return () => {
			[...nodesRef.current.values()].forEach(e => e.close());
		};
	}, [workflowData]);

	const getComponentStatus = (config: Component) => {
		if (config.status === ZSyncStatus.Skipped) {
			return 'skipped destroy';
		} else if (config.status === ZSyncStatus.SkippedReconcile) {
			return 'skipped provision';
		} else {
			return config.isDestroyed ? 'destroy' : 'provision';
		}
	};

	const getZFeedbackModal = (isApproved: boolean, needsApproval: boolean) => {
		if (!needsApproval) {
			return <></>;
		}
		return (
			<ZFeedbackModal
				approved={isApproved}
				onApprove={async () => {
					if (!auditData || auditData.length === 0) return;
					const latestId = auditData.sort((d1, d2) => d2.reconcileId - d1.reconcileId)[0].reconcileId;
					await AuditService.getInstance().approve(latestId);
				}}
				onDecline={() => {}}
			/>
		);
	};

	const getView = () => {
		switch (viewType) {
			case ViewType.Concise_Logs:
				return <ZWorkflowDiagram nodes={filteredNodes} approvedBy={''} />;
			case ViewType.Detailed_Logs:
				return (
					<div>
						<ZEditor height="80vh" data={logs || ''} />
					</div>
				);
			case ViewType.Detailed_Cost_Breakdown:
				return (
					<div className="zlifecycle-config-workflow-view__diagram--detailed-cost-breakdown">
						{config.id && <HierarchicalView data={config} />}
					</div>
				);
			case ViewType.Audit_View:
				return (
					<AuditView
						auditId={config.id}
						auditColumns={auditColumns}
						config={config}
						auditData={auditData}
						fetchLogs={AuditService.getInstance().fetchLogs.bind(
							AuditService.getInstance(),
							separatedConfigId?.team || '',
							separatedConfigId?.environment || '',
							separatedConfigId?.component || ''
						)}
					/>
				);
			case ViewType.State_File:
				if (!workflowData || !ilRepo) {
					return <></>;
				}
				return (
					<StateFileView
						componentName={separatedConfigId?.component || ''}
						environmentName={separatedConfigId?.environment || ''}
						teamName={separatedConfigId?.team || ''}
						workflowName={workflowData.metadata.name}
						ilRepo={ilRepo}
					/>
				);
		}
	};

	const getTabs = () => {
		return (
			<nav className="nav-tab">
				<ul>
					{Object.values(ViewType)
						.sort()
						.map(tabId => (
							<li
								key={tabId}
								className={`nav-tab_item nav-tab_item${viewType === tabId ? '--active' : ''}`}>
								<a
									onClick={() => {
										setViewType(tabId as number);
									}}>
									{ViewTypeTabName.get(tabId as number)}
								</a>
							</li>
						))}
				</ul>
			</nav>
		);
	};

	return (
		<>
			<div className="zlifecycle-config-workflow-view zscrollbar">
				<div className="zlifecycle-config-workflow-view__config-info">
					<div className="zlifecycle-config-workflow-view__header">
						<p className="heading">
							<span>{config.name}</span>
						</p>
					</div>
					<ConfigWorkflowLeftView config={config} />
				</div>
				{
					<div className="zlifecycle-config-workflow-view__diagram">
						{getTabs()}{' '}
						<div style={{ overflowY: 'auto', height: 'calc(100vh - 110px)', paddingRight: '20px' }}>
							{getView()}
						</div>
					</div>
				}
			</div>
		</>
	);
};

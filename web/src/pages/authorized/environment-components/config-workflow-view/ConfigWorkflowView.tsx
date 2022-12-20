import './style.scss';

import { filterLabels } from 'components/molecules/cards/EnvironmentComponentCards';
import { renderHealthStatus, renderLabels, renderSyncedStatus } from 'components/molecules/cards/renderFunctions';
import { ZEditor } from 'components/molecules/editor/Editor';
import { ZStreamRenderer } from 'components/molecules/zasync-renderer/ZStreamRenderer';
import { AuditData, AuditView } from 'components/organisms/audit_view/AuditView';
import { HierarchicalView } from 'components/organisms/hierarchical-view/hierarchical-view';
import { ZWorkflowDiagram } from 'components/organisms/workflow-diagram/WorkflowDiagram';
import { Context } from 'context/argo/ArgoUi';
import { useApi } from 'hooks/use-api/useApi';
import { ZSyncStatus } from 'models/argo.models';
import { OptionItem } from 'models/general.models';
import { EnvironmentComponentItem } from 'models/projects.models';
import { ZFeedbackModal } from 'pages/authorized/environment-components/config-workflow-view/FeedbackModal';
import React, { FC, useContext, useEffect, useState } from 'react';
import { ArgoWorkflowsService } from 'services/argo/ArgoWorkflows.service';
import { AuditService } from 'services/audit/audit.service';
import { CostingService } from 'services/costing/costing.service';

import { auditColumns, getSeparatedConfigId, ViewType, ViewTypeTabName } from '../helpers';
import { EventClientLogs } from 'utils/apiClient/EventClient';
import { StateFileView } from 'components/organisms/state-file-view/StateFileView';
import { ArgoComponentsService } from 'services/argo/ArgoComponents.service';
import { ConfigWorkflowLeftView } from './ConfigWorkflowLeftView';

type Props = {
	projectId: string;
	environmentId: string;
	config: EnvironmentComponentItem;
	workflowData: any;
	logs: string | null;
	plans: string | null;
};

export const ConfigWorkflowView: FC<Props> = (props: Props) => {
	const { projectId, environmentId, config, logs, plans, workflowData } = props;
	const clientLogMap = new Map<string, EventClientLogs>();
	const tabs: OptionItem[] = [
		{
			id: 'plans',
			name: 'Plans',
		},
		{
			id: 'logs',
			name: 'Logs',
		},
	];
	const [filteredNodes, setFilteredNodes] = useState<any>([]);
	const [viewType, setViewType] = useState<any>(ViewType.Concise_Logs);
	const [isApproved, setIsApproved] = useState<boolean>(false);
	const [approvedBy, setApprovedBy] = useState<string>('');
	const [needsApproval, setNeedsApproval] = useState<boolean>(false);
	const { fetch: fetchApprove } = useApi(ArgoWorkflowsService.approveConfigWorkflow);
	const { fetch: fetchDecline } = useApi(ArgoWorkflowsService.declineConfigWorkflow);
	const [componentStatus, setComponentStatus] = useState<ZSyncStatus>(ZSyncStatus.Unknown);
	const [ilRepo, setIlRepo] = useState<string>('');
	const separatedConfigId = config ? getSeparatedConfigId(config) : null;

	const ctx = useContext(Context);

	useEffect(() => {
		const delayedStatus = [ZSyncStatus.Destroyed, ZSyncStatus.Provisioned, ZSyncStatus.InSync];
		if (delayedStatus.includes(config.componentStatus)) {
			workflowData?.status?.phase === 'Succeeded' && setComponentStatus(config.componentStatus);
		} else {
			setComponentStatus(config.componentStatus);
		}
	}, [config?.componentStatus, workflowData?.status?.phase]);

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
		setIsApproved(isApproved);
		setNeedsApproval(needsApproval);
		if (foundApprovedNode) {
			AuditService.getInstance()
				.getApprovedBy(config.displayValue, '-1')
				.then((auditData: any) => {
					setApprovedBy(auditData?.approved_by || '');
				});
		}
		setFilteredNodes(
			newNodes
				.filter((node: any) => node.displayName !== 'notify' && node.type !== 'Steps')
				.map((node: any) => {
					const podName = node.boundaryID + '-run' + node.id.replace(node.boundaryID, '');
					const status = getComponentStatus(config, teardown);
					if ((node.displayName === 'apply' || node.displayName === 'plan') && status) {
						node.displayName = `${status} ${node.displayName}`;
					}

					const paramSet = {
						environmentId: config?.labels?.environment_id || '',
						projectId: config?.labels?.project_id || '',
						configId: config?.id || '',
						workflowId: config?.labels?.last_workflow_run_id,
					};

					const clientLogUrl = `/wf/api/v1/stream/projects/${projectId}/environments/${environmentId}/config/${config.id}/${config.labels?.last_workflow_run_id}/log/${podName}`;
					
					if (
						!clientLogMap.has(clientLogUrl) &&
						node.phase !== 'Succeeded' &&
						node.displayName !== 'approve' &&
						node.phase !== 'Failed'
					) {
						clientLogMap.set(clientLogUrl, new EventClientLogs(clientLogUrl));
					}
					
					return {
						...node,
						configStatus: config.componentStatus,
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
						getNodeLogs: () => clientLogMap.get(clientLogUrl),
					};
				})
		);
		return () => {
			[...clientLogMap.values()].forEach(o => {
				o.close();
			});
			clientLogMap.clear();
		};
	}, [workflowData]);

	const getComponentStatus = (config: EnvironmentComponentItem, teardown: string) => {
		if (config.componentStatus === ZSyncStatus.Skipped) {
			return 'skipped destroy';
		} else if (config.componentStatus === ZSyncStatus.SkippedReconcile) {
			return 'skipped provision';
		} else {
			return teardown === 'true' ? 'destroy' : teardown === 'false' ? 'provision' : '';
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
					fetchApprove({
						projectId: projectId,
						environmentId: environmentId,
						configId: config.id,
						workflowId: workflowData.metadata.name,
						data: {
							name: workflowData.metadata.name,
							namespace: 'argocd',
						},
					}).then(resp => {
						console.log(resp.data);
						Promise.resolve(
							ArgoComponentsService.patchComponentStatus(
								config.displayValue,
								ZSyncStatus.InitializingApply
							)
						);
						Promise.resolve(AuditService.getInstance().patchApprovedBy(config.name));
					});
				}}
				onDecline={() => {
					fetchDecline({
						projectId: projectId,
						environmentId: environmentId,
						configId: config.id,
						workflowId: workflowData.metadata.name,
						data: {
							message: 'no message',
							name: workflowData.id,
							namespace: 'argocd',
						},
					}).then((d: any) => {
						console.log('----------------------------------> decline', d);
					});
				}}
			/>
		);
	};

	const getView = () => {
		switch (viewType) {
			case ViewType.Concise_Logs:
				return <ZWorkflowDiagram nodes={filteredNodes} approvedBy={approvedBy} />;
			case ViewType.Detailed_Logs:
				return (
					<div>
						<ZEditor height="80vh" data={logs || ''} />
					</div>
				);
			case ViewType.Detailed_Cost_Breakdown:
				return (
					<div className="zlifecycle-config-workflow-view__diagram--detailed-cost-breakdown">
						{config.id && (
							<HierarchicalView data={config} componentId={config.id} />
						)}
					</div>
				);
			case ViewType.Audit_View:
				return (
					<AuditView
						auditId={config.id || ''}
						auditColumns={auditColumns}
						config={config}
						fetch={AuditService.getInstance().getComponent}
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
							<span>{config.componentName}</span>
						</p>
					</div>
					<ConfigWorkflowLeftView config={config} configLabels={config.labels} />
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

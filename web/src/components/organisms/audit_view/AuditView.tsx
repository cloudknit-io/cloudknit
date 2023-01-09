import './styles.scss';

import { ReactComponent as Chevron } from 'assets/images/icons/chevron-right.svg';
import { Loader } from 'components/atoms/loader/Loader';
import { ZTable } from 'components/atoms/table/Table';
import { NodeIcon, NodeStatus, ZDiagramNode } from 'components/molecules/diagram-node/DiagramNode';
import { ZEditor } from 'components/molecules/editor/Editor';
import { isArray } from 'lodash';
import React from 'react';
import { useState } from 'react';
import { useEffect } from 'react';
import { FC } from 'react';
import { Subject } from 'rxjs';
import { AuditService } from 'services/audit/audit.service';

import { ZAccordion, ZAccordionItem } from '../accordion/ZAccordion';
import { AuditStatus, ZSyncStatus } from 'models/argo.models';
import { SmallText } from '../workflow-diagram/WorkflowDiagram';
import { EnvironmentComponentItem } from 'models/projects.models';
import { CompAuditData, Component, EnvAuditData } from 'models/entity.store';

type AuditData = EnvAuditData | CompAuditData;

type Props = {
	auditId: number;
	auditData: EnvAuditData[] | CompAuditData[];
	auditColumns: any[];
	fetchLogs?: (auditId: number) => Promise<any>;
	config?: Component;
};

export const AuditView: FC<Props> = ({ auditId, auditData, auditColumns, fetchLogs, config }: Props) => {
	const auditServiceInstance = AuditService.getInstance();
	const [selectedLog, setSelectedLog] = useState<AuditData | null>(null);
	const [logs, setLogs] = useState<ZAccordionItem[] | null>();
	const nodeOrder = ['plan', 'apply'];
	const [latestReconcileId, setLatestReconcileId] = useState(-1);
	const getOperation = (status: AuditStatus) => {
		switch (status) {
			case AuditStatus.Provisioned:
			case AuditStatus.Provisioning:
			case AuditStatus.ProvisionApplyFailed:
			case AuditStatus.ProvisionPlanFailed:
			case AuditStatus.Env_Provisioning:
			case AuditStatus.Env_Provision_Ended:
				return 'Provision';
			case AuditStatus.Destroyed:
			case AuditStatus.Destroying:
			case AuditStatus.DestroyApplyFailed:
			case AuditStatus.DestroyPlanFailed:
				return 'Destroy';
			case AuditStatus.Env_Destroying:
			case AuditStatus.Env_Destroy_Ended:
				return 'Teardown';
			case AuditStatus.Skipped:
			case AuditStatus.SkippedDestroy:
				return 'Skipped Destroy';
			case AuditStatus.SkippedProvision:
				return 'Skipped Provision';
			case AuditStatus.SkippedReconcile:
				return 'Skipped Reconcile';
		}
		return '--';
	};

	useEffect(() => {
		if (!auditData) return;
		setSelectedLog(null);
		let recId = -1;
		auditData.forEach((d: AuditData) => {
			if (recId < d.reconcileId) {
				recId = d.reconcileId;
			}
		});

		setLatestReconcileId(recId);
	}, [auditData]);

	// useEffect(() => {
	// 	if (isArray(auditData)) auditServiceInstance.setAuditCache(auditId, auditData);
	// }, [auditData]);

	const getNodeStatus = (status: string): NodeStatus => {
		switch (status) {
			case AuditStatus.Failed:
			case AuditStatus.ProvisionApplyFailed:
			case AuditStatus.ProvisionPlanFailed:
			case AuditStatus.DestroyApplyFailed:
			case AuditStatus.DestroyPlanFailed:
				return 'Failed';
			case AuditStatus.Initialising:
			case AuditStatus.Initializing:
			case AuditStatus.Provisioning:
			case AuditStatus.Destroying:
				return 'InProcess';
			case 'success':
			case AuditStatus.Provisioned:
			case AuditStatus.Destroyed:
			case AuditStatus.Success:
				return 'Succeeded';
			case 'Mutated':
				return 'Mutated';
		}
		return 'InProcess';
	};

	const getNodeIcon = (status: string): NodeIcon => {
		switch (status) {
			case AuditStatus.Failed:
			case AuditStatus.ProvisionApplyFailed:
			case AuditStatus.ProvisionPlanFailed:
			case AuditStatus.DestroyApplyFailed:
			case AuditStatus.DestroyPlanFailed:
				return 'Failed';
			case AuditStatus.Initialising:
			case AuditStatus.Initializing:
			case AuditStatus.Provisioning:
			case AuditStatus.Destroying:
				return 'InProcess';
			case 'success':
			case AuditStatus.Success:
			case AuditStatus.Provisioned:
			case AuditStatus.Destroyed:
				return 'Synced';
			case 'Mutated':
				return 'Mutated';
		}
		return 'InProcess';
	};

	const getSummaryAndMutatedData = (body: string) => {
		try {
			if (!body) {
				return ['', false];
			}
			let summary =
				body.match(/^Plan:.*$/gm) ||
				body.match(/^Apply complete!.*$/gm) ||
				body.match(/^No changes\. Infrastructure is up-to-date\.$/gm);
			if (!summary) {
				return ['', false];
			}
			return [summary[0], summary[0].includes('Plan:')];
		} catch (err) {
			return ['', false];
		}
	};

	const getStatus = (i: number, d: any, m: boolean, r: AuditStatus) => {
		const status: AuditStatus = i < d.length - 1 ? AuditStatus.Success : r;
		if (m && [AuditStatus.Provisioned, AuditStatus.Destroyed, AuditStatus.Success].includes(status)) {
			return 'Mutated';
		} else {
			return status;
		}
	};

	const renderItems = () => {
		if (!auditData) {
			return (
				<div style={{ display: 'flex', justifyContent: 'center' }}>
					<Loader height={16} width={16} />
				</div>
			);
		}

		if (selectedLog) {
			if (logs) {
				return (
					<div className={`zlifecycle-audit-logs`}>
						<div
							className="hide-logs"
							onClick={() => {
								setLogs(null);
								setSelectedLog(null);
							}}>
							<Chevron /> Back to runs
						</div>
						<ZAccordion items={logs} />
					</div>
				);
			} else {
				return <Loader height={16} width={16} />;
			}
		}

		return (
			<div className="zlifecycle-audit-table">
				<ZTable
					table={{
						columns: auditColumns,
						rows: auditData.sort((a: AuditData, b: AuditData) => b.reconcileId - a.reconcileId),
					}}
					onRowClick={(rowData: AuditData) => {
						if (
							!fetchLogs ||
							[
								AuditStatus.SkippedReconcile,
								AuditStatus.SkippedDestroy,
								AuditStatus.Skipped,
								AuditStatus.SkippedProvision,
							].includes(rowData.status.toLowerCase() as AuditStatus)
						) {
							return;
						}

						setSelectedLog(rowData);
						fetchLogs(rowData.reconcileId).then(({ data }) => {
							let zi: ZAccordionItem[] = [];
							if (data && data === 'No Object was found') {
								const item: ZAccordionItem = {
									accordionContent: 'No logs were found',
									accordionHeader: 'No logs were found',
									collapsed: true,
								};
								zi.push(item);
							} else if (data.length > 0) {
								const items: ZAccordionItem[] = data
									.map((item: any) => {
										const sm = getSummaryAndMutatedData(item.body);
										return {
											key: item.key.split('/').splice(-1)[0].replace('_output', ''),
											body: item.body,
											summary: sm[0],
											mutated: sm[1],
										};
									})
									.sort((a: any, b: any) => {
										return nodeOrder.indexOf(a.key) - nodeOrder.indexOf(b.key);
									})
									.map((item: any, index: number) => ({
										accordionHeader: (
											<div
												className={`zlifecycle-audit-logs-${getNodeStatus(
													getStatus(index, data, item.mutated, rowData.status)
												)}`}>
												<ZDiagramNode
													text={item.key}
													icon={getNodeIcon(
														getStatus(index, data, item.mutated, rowData.status)
													)}
													status={getNodeStatus(
														getStatus(index, data, item.mutated, rowData.status)
													)}
												/>
												<SmallText data={item.summary} />
											</div>
										),
										accordionContent: <ZEditor data={item.body} />,
										collapsed: true,
									}));
								zi = [...items];
							}
							setLogs(zi);
						});
					}}
					rowHeight={40}
					rowConditionalClass={(data: any) => {
						if (
							[
								AuditStatus.Initialising,
								AuditStatus.Initializing,
								AuditStatus.Provisioning,
								AuditStatus.Destroying,
							].includes(data?.status)
						) {
							return 'zlifecycle-audit-table-row zlifecycle-audit-table-row-initializing';
						}
						if (
							data?.status === 'Success' ||
							data?.status === 'success' ||
							data?.status === 'ended' ||
							data?.status === 'destroy_ended' ||
							[
								AuditStatus.Provisioned,
								AuditStatus.Destroyed,
								AuditStatus.Skipped,
								AuditStatus.SkippedReconcile,
							].includes(data?.status)
						) {
							return 'zlifecycle-audit-table-row zlifecycle-audit-table-row-success';
						}
						if (
							[
								AuditStatus.Failed,
								AuditStatus.DestroyApplyFailed,
								AuditStatus.DestroyPlanFailed,
								AuditStatus.ProvisionPlanFailed,
								AuditStatus.ProvisionApplyFailed,
							].includes(data?.status)
						) {
							return 'zlifecycle-audit-table-row zlifecycle-audit-table-row-failed';
						}
						return '';
					}}
				/>
			</div>
		);
	};

	return renderItems();
};

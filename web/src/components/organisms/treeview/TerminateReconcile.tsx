import { Environment } from 'models/entity.type';
import { ReactComponent as TerminateIcon } from 'assets/images/icons/card-status/sync/Not Sync.svg';
import { AuditService } from 'services/audit/audit.service';
import { useContext, useEffect, useState } from 'react';
import { Context } from 'context/argo/ArgoUi';
import { NotificationType } from 'components/argo-core';
import { ZSyncStatus } from 'models/argo.models';
import { ArgoStreamService } from 'services/argo/ArgoStream.service';

export type TerminateReconcileProps = {
	environment: Environment;
};
export const TerminateReconcile: React.FC<TerminateReconcileProps> = ({ environment }) => {
	const notifications = useContext(Context)?.notifications;
	const [reconciling, setReconciling] = useState<boolean>(false);

	useEffect(() => {
		const ecd = ArgoStreamService.streamEnvironment(environment.argoId);
		const watcherSub = ecd.listen().subscribe((e: any) => {
			const healthStatus = e?.data?.result?.application?.status?.health?.status;
			setReconciling(healthStatus === 'Progressing');
		});

		return () => {
			watcherSub.unsubscribe();
		}
	}, [environment]);
	return (
		<button
			disabled={!reconciling}
			className="dag-controls-terminate"
			onClick={async (e: any) => {
				if (!reconciling) return;
				e.stopPropagation();
				if (!window.confirm('Are you sure you want to cancel the reconcile?')) return;
				const response = await AuditService.getInstance().terminate(environment.latestEnvRecon.reconcileId);
				if (response.status === 201) {
					notifications?.show({
						content: 'Reconcile Cancelled',
						type: NotificationType.Success,
					});
				} else {
					notifications?.show({
						content: 'Reconcile Cancellation Failed',
						type: NotificationType.Error,
					});
				}
			}}>
			<TerminateIcon title="Cancel Reconcile" />
			Cancel
		</button>
	);
};

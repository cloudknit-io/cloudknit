import { Environment } from 'models/entity.type';
import { ReactComponent as TerminateIcon } from 'assets/images/icons/card-status/sync/Not Sync.svg';
import { AuditService } from 'services/audit/audit.service';
import { useContext } from 'react';
import { Context } from 'context/argo/ArgoUi';
import { NotificationType } from 'components/argo-core';
import { ZSyncStatus } from 'models/argo.models';

export type TerminateReconcileProps = {
	environment: Environment;
};
export const TerminateReconcile: React.FC<TerminateReconcileProps> = ({ environment }) => {
	const notifications = useContext(Context)?.notifications;
	return (
		<button
			disabled={![
				ZSyncStatus.Initializing,
				ZSyncStatus.Destroying,
				ZSyncStatus.Provisioning,
			].includes(environment.status as ZSyncStatus)}
			className="dag-controls-terminate"
			onClick={async (e: any) => {
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

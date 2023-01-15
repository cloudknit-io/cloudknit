import { NotificationType } from 'components/argo-core/notifications/notifications';
import { Context } from 'context/argo/ArgoUi';
import { OperationPhase, OperationPhases } from 'models/argo.models';
import { EntityStore, Environment, Team } from 'models/entity.store';
import React, { FC, useCallback, useEffect, useState } from 'react';
import { ArgoEnvironmentsService } from 'services/argo/ArgoEnvironments.service';
import { ArgoStreamService } from 'services/argo/ArgoStream.service';
import { subscriberWatcher } from 'utils/apiClient/EventClient';
import { hardSync } from './helpers';

export type ReconcilerProps = {
	environment: Environment;
	template: (environment: Environment, reconciling: boolean, triggerSync: () => Promise<any>) => React.ReactElement;
};

export const Reconciler: FC<ReconcilerProps> = ({ environment, template }) => {
	const nm = React.useContext(Context)?.notifications;
	const [watcherStatus, setWatcherStatus] = useState<OperationPhase>();
	const [reconciling, setReconciling] = useState<boolean>(false);
	const [syncStarted, setSyncStarted] = useState<boolean>(false);

    useEffect(() => {
		if (!syncStarted) {
			return;
		}
		setTimeout(() => {
			setSyncStarted(false);
		}, 10000);
	}, [syncStarted]);

	useEffect(() => {
		if (!environment) {
			return;
		}
        setReconciling(false);
        setSyncStarted(false);
		const ecd = ArgoStreamService.streamEnvironment(environment.argoId);
		const watcherSub = ecd.listen().subscribe((e: any) => {
            const healthStatus = e?.result?.application?.status?.health?.status;
            setReconciling(healthStatus === 'Progressing')
			// if (e?.application?.metadata?.name?.replace('-team-watcher', '') === environment.argoId) {
			// 	const status = e?.application?.status?.operationState?.phase;
			// 	setWatcherStatus(status);
			// }
		});
		return () => {
			watcherSub.unsubscribe();
			ecd.close();
		};
	}, [environment]);

	const syncMe = useCallback(async () => {
		if (syncStarted) {
			return;
		}
		setSyncStarted(true);
		try {
			nm?.show({
				content: `Reconciling ${environment.name}`,
				type: NotificationType.Success,
			});
			if (watcherStatus === OperationPhases.Failed) {
				await hardSync(
					EntityStore.getInstance().getTeam(environment.teamId)?.name || '',
					environment.name || ''
				);
				return;
			}
			await ArgoEnvironmentsService.deleteEnvironment(environment.argoId as string);
			setTimeout(async () => {
				await ArgoEnvironmentsService.syncEnvironment(environment.argoId as string);
			}, 1500);
		} catch (err) {
			// if ((err as any)?.message?.includes('not found') && environment.syncStatus === 'OutOfSync') {
			//     await ArgoEnvironmentsService.syncEnvironment(environment.id as string);
			// }
			// console.log(err);
		}
	}, [environment]);

	return template(environment, reconciling || syncStarted, syncMe);
};

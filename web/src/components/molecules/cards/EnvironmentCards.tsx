import './style.scss';

import { ReactComponent as AWSIcon } from 'assets/images/icons/AWS.svg';
import { ReactComponent as SyncIcon } from 'assets/images/icons/sync-icon.svg';
import { renderEnvSyncedStatus } from 'components/molecules/cards/renderFunctions';
import { ZGridDisplayListWithLabel } from 'components/molecules/grid-display-list/GridDisplayList';
import { ESyncStatus, OperationPhase, ZSyncStatus } from 'models/argo.models';
import { ListItem } from 'models/general.models';
import { EnvironmentItem, EnvironmentsList } from 'models/projects.models';
import { getEnvironmentErrorCondition, renderCost, syncMe } from 'pages/authorized/environments/helpers';
import React, { FC, useMemo } from 'react';
import { useEffect } from 'react';
import { useState } from 'react';
import { useHistory } from 'react-router-dom';
import { ArgoMapper } from 'services/argo/ArgoMapper';
import { subscriber, subscriberWatcher } from 'utils/apiClient/EventClient';
import { Context } from 'context/argo/ArgoUi';
import { NotificationsApi } from 'components/argo-core/notifications/notification-manager';
import { FeatureKeys, featureToggled } from 'pages/authorized/feature_toggle';
import { ZSidePanel } from '../side-panel/SidePanel';
import { ErrorView } from 'components/organisms/error-view/ErrorView';
import { ErrorStateService } from 'services/error/error-state.service';
import { eventErrorColumns } from 'models/error.model';
import { EntityStore, Environment } from 'models/entity.store';

type Props = {
	environments: Environment[];
	compareEnabled?: any;
};

type PropsEnvironmentItem = {
	environment: Environment;
	notificationManager?: NotificationsApi;
	compareEnabled?: any;
};

const environmentTeam = (environment: EnvironmentItem): string => {
	return environment.labels ? environment.labels.project_id : 'Unknown';
};

export const environmentName = (environment: EnvironmentItem | undefined): string => {
	if (environment) return (environment.labels || {}).env_name || '';
	return '';
};

const environmentDescription = (environment: EnvironmentItem): string => {
	if ((environment.name || '').match(/dev/) !== null) {
		return 'Development environment for testing new features';
	}

	if ((environment.name || '').match(/prod/) !== null) {
		return 'Production environment deployed to client sites';
	}
	return 'Local use only';
};

export const EnvironmentCards: FC<Props> = ({ environments, compareEnabled }: Props) => {
	const contextApi = React.useContext(Context);
	// const hide = (env: EnvironmentItem) => {
	// 	return environments.some(ev => ev.id === `${environmentTeam(env)}-${environmentName(env)}` && ev !== env);
	// };
	return (
		<div className="bottom-offset">
			<div className="com-cards">
				{environments.map(
					(environment: Environment, _i) => (
						// environment.labels?.failed_environment ? (
						// 	!hide(environment) && <FailedEnvironmentCard key={`card-${_i}`} environment={environment} />
						// ) : (
						<EnvironmentCard
							key={`card-${_i}`}
							environment={environment}
							notificationManager={contextApi?.notifications}
							compareEnabled={compareEnabled}
						/>
					)
					// )
				)}
			</div>
		</div>
	);
};

// export const FailedEnvironmentCards: FC<Props> = ({ environments }: Props) => {
// 	const contextApi = React.useContext(Context);
// 	return (
// 		<div className="bottom-offset" id="failed-environments">
// 			<div className="com-cards">
// 				{environments.map((environment: EnvironmentItem, _i) => (
// 					<FailedEnvironmentCard key={`card-${_i}`} environment={environment} />
// 				))}
// 			</div>
// 		</div>
// 	);
// };

export const EnvironmentCard: FC<PropsEnvironmentItem> = ({
	environment,
	notificationManager,
	compareEnabled,
}: PropsEnvironmentItem) => {
	const entityStore = useMemo(() => EntityStore.getInstance(), []);
	const ref = React.createRef<any>();
	const [env, setEnv] = useState<Environment>(environment);
	const [gridItems, setGridItems] = useState<ListItem[]>([]);
	const [syncStarted, setSyncStarted] = useState<boolean>(false);
	const [watcherStatus, setWatcherStatus] = useState<OperationPhase | undefined>();
	const [environmentCondition, setEnvironmentCondition] = useState<any>(null);
	const [selected, setSelected] = useState<boolean>(false);

	useEffect(() => {
		setEnv(environment);
	}, [environment]);

	useEffect(() => {
		const gi = mapGridItems(env);
		setGridItems(gi);
	}, [env]);

	const mapGridItems = (environment: Environment): ListItem[] => {
		const gridItems = [
			{
				label: 'Team',
				value: entityStore.getTeam(env.teamId)?.name,
			},
			{
				label: 'Name',
				value: environment.name,
			},
			{
				label: 'Cost',
				value: -1,
			},
			{
				label: 'Cloud',
				value: <>{<AWSIcon />}</>,
			},
		];

		// if (environment.labels?.env_status) {
		// 	gridItems.splice(2, 0, {
		// 		label: 'Status',
		// 		value: (
		// 			<>
		// 				{renderEnvSyncedStatus(
		// 					environment.labels.env_status as ZSyncStatus,
		// 					'',
		// 					'',
		// 					env.syncFinishedAt
		// 				)}
		// 			</>
		// 		),
		// 	});
		// }

		return gridItems;
	};
	const getRoutePath = (environment: Environment) => {
		const state: any = history.location.state;
		if (state && 'type' in state) {
			return `/${EntityStore.getInstance().getTeam(environment.teamId)?.name}/${environment.name}/${state.type}`;
		}
		return `/${EntityStore.getInstance().getTeam(environment.teamId)?.name}/${environment.name}`;
	};

	// TODO: Sync Status
	const getSyncIconClass = (environment: any) => {
		if (environment.syncStatus === ESyncStatus.OutOfSync) {
			return '--out-of-sync';
		} else if (environment.syncStatus === ESyncStatus.Synced) {
			return '--in-sync';
		} else {
			return '--unknown';
		}
	};

	// TODO: Watcher for environemts
	// useEffect(() => {
	// 	const watcherSub = subscriberWatcher.subscribe(e => {
	// 		if (e?.application?.metadata?.name?.replace('-team-watcher', '') === environment?.labels?.project_id) {
	// 			const status = e?.application?.status?.operationState?.phase;
	// 			setWatcherStatus(status);
	// 		}
	// 	});
	// 	if (environment.conditions?.length > 0) {
	// 		setEnvironmentCondition(getEnvironmentErrorCondition(environment.conditions));
	// 	}
	// 	return () => watcherSub.unsubscribe();
	// }, [environment]);

	useEffect(() => {
		if (!syncStarted) {
			return;
		}
		setTimeout(() => {
			setSyncStarted(false);
		}, 10000);
	}, [syncStarted]);

	useEffect(() => {
		selected ? ref.current.classList.add('compare-selected') : ref.current.classList.remove('compare-selected');
	}, [selected]);

	const history = useHistory();
	return (
		<div
			ref={ref}
			className={`environment-card com-card com-card--with-header environment-card--${
				'Unkown'// TODO: env status env.labels?.env_status
			}`}
			onClick={(e): void => {
				history.push(getRoutePath(env));
			}}>
			<div className="environment-card__header com-card__header">
				<div className="com-card__header__title">
					<div>
						<h4>
							{EntityStore.getInstance().getTeam(environment.teamId)?.name}: {env.name}
						</h4>
					</div>
				</div>
				<div className="large-health-icon-container">
					{
						<SyncIcon
							className={`large-health-icon-container__sync-button large-health-icon-container__sync-button${getSyncIconClass(
								env
							)} large-health-icon-container__sync-button${
								'Progressing'
								// env.healthStatus === 'Progressing' || syncStarted ? '--in-progress' : ''
							}`}
							title={environmentCondition || 'Reconcile Environment'}
							onClick={async e => {
								e.stopPropagation();
								//TODO: Syncing env
								// if (env.healthStatus !== 'Progressing')
								// 	await syncMe(
								// 		env,
								// 		syncStarted,
								// 		setSyncStarted,
								// 		notificationManager as NotificationsApi,
								// 		watcherStatus
								// 	);
							}}
						/>
					}
					{featureToggled(FeatureKeys.DIFF_CHECKER, true) && (
						<input
							type="checkbox"
							className="select-compare"
							checked={selected}
							onClick={e => {
								e.stopPropagation();
								if (!selected) {
									const canSet = compareEnabled.setSelectedEnvironmentToCompare(env, true);
									setSelected(canSet);
								} else {
									compareEnabled.setSelectedEnvironmentToCompare(env, false);
									setSelected(false);
								}
							}}
						/>
					)}
				</div>
			</div>
			<div className="com-card__body">
				<ZGridDisplayListWithLabel items={gridItems} />
			</div>
		</div>
	);
};

// export const FailedEnvironmentCard: FC<PropsEnvironmentItem> = ({ environment }: PropsEnvironmentItem) => {
// 	const ref = React.createRef<any>();
// 	const [env, setEnv] = useState<EnvironmentItem>(environment);
// 	const [gridItems, setGridItems] = useState<ListItem[]>([]);
// 	const [sidePanel, setSidePanel] = useState<boolean>(false);

// 	useEffect(() => {
// 		const gi = mapGridItems(env);
// 		setGridItems(gi);
// 	}, [env]);

// 	const mapGridItems = (environment: EnvironmentItem): ListItem[] => {
// 		const gridItems = [
// 			{
// 				label: 'Team',
// 				value: environmentTeam(env),
// 			},
// 			{
// 				label: 'Name',
// 				value: env.name,
// 			},
// 			{
// 				label: 'Cost',
// 				value: 'N/A',
// 			},
// 			{
// 				label: 'Cloud',
// 				value: <>{<AWSIcon />}</>,
// 			},
// 		];

// 		gridItems.splice(2, 0, {
// 			label: 'Status',
// 			value: <>{renderEnvSyncedStatus(ZSyncStatus.ProvisionFailed, '', '', env.syncFinishedAt)}</>,
// 		});

// 		return gridItems;
// 	};

// 	const getErrorTitle = () => {
// 		return `${environmentTeam(env) + ':'}${environmentName(env)}`;
// 	};

// 	return (
// 		<>
// 			<ZSidePanel
// 				isShown={sidePanel}
// 				onClose={() => {
// 					setSidePanel(false);
// 				}}>
// 				<ErrorView
// 					id={env.name}
// 					columns={eventErrorColumns}
// 					dataRows={ErrorStateService.getInstance().errorsInEnvironment(env.labels?.env_name || '')}
// 				/>
// 			</ZSidePanel>
// 			<div
// 				ref={ref}
// 				className={`environment-card com-card com-card--with-header environment-card`}
// 				onClick={(e): void => {
// 					setSidePanel(true);
// 				}}>
// 				<div className="environment-card__header com-card__header">
// 					<div className="com-card__header__title">
// 						<div>
// 							<h4>
// 								{environmentTeam(env)}: {environmentName(env)}
// 							</h4>
// 						</div>
// 					</div>
// 					<div className="large-health-icon-container"></div>
// 				</div>
// 				<div className="com-card__body">
// 					<ZGridDisplayListWithLabel items={gridItems} />
// 				</div>
// 			</div>
// 		</>
// 	);
// };

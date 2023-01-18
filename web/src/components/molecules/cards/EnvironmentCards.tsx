import './style.scss';

import { ReactComponent as AWSIcon } from 'assets/images/icons/AWS.svg';
import { CostRenderer, renderEnvSyncedStatus } from 'components/molecules/cards/renderFunctions';
import { ZGridDisplayListWithLabel } from 'components/molecules/grid-display-list/GridDisplayList';
import { OperationPhase, ZEnvSyncStatus, ZSyncStatus } from 'models/argo.models';
import { ListItem } from 'models/general.models';
import { EnvironmentItem } from 'models/projects.models';
import { EnvCardReconcile } from 'pages/authorized/environments/helpers';
import React, { FC, useMemo } from 'react';
import { useEffect } from 'react';
import { useState } from 'react';
import { useHistory } from 'react-router-dom';
import { Context } from 'context/argo/ArgoUi';
import { NotificationsApi } from 'components/argo-core/notifications/notification-manager';
import { FeatureKeys, featureToggled } from 'pages/authorized/feature_toggle';
import { EntityStore } from 'models/entity.store';
import { Reconciler } from 'pages/authorized/environments/Reconciler';
import { Environment } from 'models/entity.type';

type Props = {
	environments: Environment[];
	compareEnabled?: any;
};

type PropsEnvironmentItem = {
	environment: Environment;
	notificationManager?: NotificationsApi;
	compareEnabled?: any;
};

export const environmentName = (environment: EnvironmentItem | undefined): string => {
	if (environment) return (environment.labels || {}).env_name || '';
	return '';
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
	const reconciling = [
		ZEnvSyncStatus.Provisioning,
		ZEnvSyncStatus.Destroying,
		ZEnvSyncStatus.Initializing,
	].includes(environment?.status as ZEnvSyncStatus) || syncStarted;

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
				value: <CostRenderer data={environment.estimatedCost} />,
			},
			{
				label: 'Cloud',
				value: <>{<AWSIcon />}</>,
			},
		];

		if (environment.status) {
			gridItems.splice(2, 0, {
				label: 'Status',
				value: (
					<>
						{renderEnvSyncedStatus(
							environment.status as ZSyncStatus,
							'',
							'',
							env.lastReconcileDatetime.toString()
						)}
					</>
				),
			});
		}

		return gridItems;
	};
	const getRoutePath = (environment: Environment) => {
		const state: any = history.location.state;
		if (state && 'type' in state) {
			return `/${EntityStore.getInstance().getTeam(environment.teamId)?.name}/${environment.name}/${state.type}`;
		}
		return `/${EntityStore.getInstance().getTeam(environment.teamId)?.name}/${environment.name}`;
	};

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
						environment && <Reconciler key={environment.argoId} environment={environment} template={EnvCardReconcile}/>
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
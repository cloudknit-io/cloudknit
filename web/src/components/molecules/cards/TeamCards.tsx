import './style.scss';

import { ZText } from 'components/atoms/text/Text';
import { renderEnvSyncedStatus } from 'components/molecules/cards/renderFunctions';
import { ZGridDisplayListWithLabel } from 'components/molecules/grid-display-list/GridDisplayList';
import { ListItem } from 'models/general.models';
import { EnvironmentItem, EnvironmentsList, TeamItem, TeamsList } from 'models/projects.models';
import { renderCost } from 'pages/authorized/teams/helpers';
import React, { FC, useEffect, useState } from 'react';
import { useHistory } from 'react-router-dom';
import { subscriber } from 'utils/apiClient/EventClient';
import { useApi } from 'hooks/use-api/useApi';
import { ArgoEnvironmentsService } from 'services/argo/ArgoEnvironments.service';
import { ApplicationWatchEvent, ZSyncStatus } from 'models/argo.models';
import { streamMapper } from 'helpers/streamMapper';
import { ArgoMapper } from 'services/argo/ArgoMapper';
import { Loader } from 'components/atoms/loader/Loader';
import AuthStore from 'auth/AuthStore';

type Props = {
	teams: TeamsList;
};

type TeamItemProps = {
	team: TeamItem;
};

const getEnvironmentStatus = (team: any, children: any[]) => {
	const envLimit = 3;
	return (
		<div className="environment-status">
			{team.resources?.length && children.length > 0 ? (
				<>
					<div className="environment-status__preview">
						{children?.slice(0, envLimit).map(resource => (
							<div>
								<b title={resource.name.replace(`${team.id}-`, '')}>
									{resource.name.replace(`${team.id}-`, '')}
								</b>
								{renderEnvSyncedStatus((resource.labels?.env_status as ZSyncStatus) || 'Unknown')}
							</div>
						))}
						{children.length > envLimit && <div>More...</div>}
					</div>
					{children.length > envLimit && (
						<div className="environment-status__tooltip">
							{children.map(resource => (
								<div>
									<b title={resource.name.replace(`${team.id}-`, '')}>
										{resource.name.replace(`${team.id}-`, '')}
									</b>
									{renderEnvSyncedStatus((resource.labels?.env_status as ZSyncStatus) || 'Unknown')}
								</div>
							))}
						</div>
					)}
				</>
			) : (
				<>{team.resources?.length ? <Loader height={16} width={16} /> : <></>}</>
			)}
		</div>
	);
};
const items = (team: TeamItem, children: EnvironmentItem[]): ListItem[] => {
	return [
		{
			label: 'Cost',
			value: renderCost(team.id),
		},
		{
			label: `Envs (${team.resources?.length || 0})`,
			value: getEnvironmentStatus(team, children),
		},
	];
};

export const TeamCards: FC<Props> = ({ teams }: Props) => {
	return (
		<div className="team com-cards">
			{teams.map((team: TeamItem, _i) => (
				<TeamCard team={team} key={_i} />
			))}
		</div>
	);
};

export const TeamCard: FC<TeamItemProps> = ({ team }: TeamItemProps) => {
	const history = useHistory();
	const { fetch } = useApi(ArgoEnvironmentsService.getEnvironments);
	const [streamData, setStreamData] = useState<ApplicationWatchEvent | null>(null);
	const [environments, setEnvironments] = useState<EnvironmentsList>([]);
	const teamId = (team.id || '').replace(AuthStore.getOrganization()?.name + '-', '');

	useEffect(() => {
		const $subscription = subscriber.subscribe(response => {
			setStreamData(response);
		});

		return (): void => $subscription.unsubscribe();
	}, []);

	useEffect(() => {
		const newItems = streamMapper<EnvironmentItem>(
			streamData,
			environments,
			ArgoMapper.parseEnvironment,
			'environment'
		);
		setEnvironments(newItems);
	}, [streamData, environments]);

	useEffect(() => {
		fetch(teamId).then(({ data }) => {
			if (data) {
				setEnvironments(data);
			}
		});
	}, [team]);

	return (
		<div
			className="com-card com-card--with-header"
			onClick={(): void => {
				history.push('/' + teamId);
			}}>
			<div className="com-card__header">
				<ZText.Body>{team.name}</ZText.Body>
			</div>
			<div className="com-card__body">
				<ZGridDisplayListWithLabel items={items(team, environments)} />
			</div>
		</div>
	);
};

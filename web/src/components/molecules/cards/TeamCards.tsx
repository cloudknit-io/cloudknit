import './style.scss';

import { Loader } from 'components/atoms/loader/Loader';
import { ZText } from 'components/atoms/text/Text';
import { CostRenderer, renderEnvSyncedStatus } from 'components/molecules/cards/renderFunctions';
import { ZGridDisplayListWithLabel } from 'components/molecules/grid-display-list/GridDisplayList';
import { ZSyncStatus } from 'models/argo.models';
import { ListItem } from 'models/general.models';
import React, { FC, useEffect, useState } from 'react';
import { useHistory } from 'react-router-dom';
import { Environment, Team } from 'models/entity.type';

type Props = {
	teams: Team[];
};

type TeamItemProps = {
	team: Team;
};

const getEnvironmentStatus = (envs: Environment[]) => {
	const envLimit = 3;
	return (
		<div className="environment-status">
			{envs?.length > 0 ? (
				<>
					<div className="environment-status__preview">
						{envs?.slice(0, envLimit).map(resource => (
							<div>
								<b title={resource.name}>
									{resource.name}
								</b>
								{renderEnvSyncedStatus(resource.status as ZSyncStatus)}
							</div>
						))}
						{envs.length > envLimit && <div>More...</div>}
					</div>
					{envs.length > envLimit && (
						<div className="environment-status__tooltip">
							{envs.map(resource => (
								<div>
									<b title={resource.name}>
										{resource.name}
									</b>
									{renderEnvSyncedStatus(resource.status as ZSyncStatus)}
								</div>
							))}
						</div>
					)}
				</>
			) : (
				<>{envs?.length ? <Loader height={16} width={16} /> : <></>}</>
			)}
		</div>
	);
};

const items = (team: Team): ListItem[] => {
	return [
		{
			label: 'Cost',
			value: <CostRenderer data={team.estimatedCost} />,
		},
		{
			label: `Envs (${team.environments.length})`,
			value: getEnvironmentStatus(team.environments),
		},
	];
};

export const TeamCards: FC<Props> = ({ teams }: Props) => {
	return (
		<div className="team com-cards">
			{teams.map((team: Team, _i) => (
				<TeamCard team={team} key={_i} />
			))}
		</div>
	);
};

export const TeamCard: FC<TeamItemProps> = ({ team }: TeamItemProps) => {
	const history = useHistory();
	const [teamItem, setTeam] = useState<Team>(team);
	const teamId = team.name;

	useEffect(() => {
		setTeam(teamItem);	
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
				<ZGridDisplayListWithLabel items={items(team)} />
			</div>
		</div>
	);
};

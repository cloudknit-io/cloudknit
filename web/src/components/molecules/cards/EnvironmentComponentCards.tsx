import './style.scss';

import { ZText } from 'components/atoms/text/Text';
import { ZModelCard } from 'components/molecules/cards/Card';
import { renderHealthStatus, renderLabels, renderSyncedStatus } from 'components/molecules/cards/renderFunctions';
import { ZSyncStatus } from 'models/argo.models';
import { EnvironmentComponentItem, EnvironmentComponentsList } from 'models/projects.models';
import { renderCost as renderCostEnv } from 'pages/authorized/environments/helpers';
import React, { FC, useEffect, useState } from 'react';
import { Component } from 'models/entity.store';

type Props = {
	components: Component[];
	projectId: string;
	envName: string;
	onClick: Function;
	showAll?: boolean;
	selectedConfig?: Component;
	workflowPhase?: string;
};

type EnvironmentComponentItemProps = {
	config: Component;
	showAll?: boolean;
	onClick: Function;
	isSelected?: boolean;
	workflowPhase?: string;
};

const getFullName = (teamId = '', environmentId = '', componentName = '') => {
	return `${teamId} : ${getEnvironmentName(teamId, environmentId)} : ${componentName}`;
};

const getEnvironmentName = (teamId = '', environmentId = '') => {
	return environmentId.replace(`${teamId}-`, '');
};

export const filterLabels = (config: Component): { [name: string]: string } => {
	const labels: any = {};
	labels.team_id = labels.project_id;
	labels.dependsOn = config.dependsOn.toString();
	return labels;
};

const totalCost = (components: EnvironmentComponentsList): string => {
	let cost = 0;

	{
		components.map((config: EnvironmentComponentItem) => (cost += parseFloat(config.componentCost)));
	}

	return cost.toFixed(3).toString();
};

const mapGridItems = (component: Component, componentStatus: ZSyncStatus) => {
	return (
		<>
			{renderSyncedStatus(
				componentStatus
			)}
		</>
	);
};

const getLabels = (component: Component): any => {
	return <></>;
};

export const EnvironmentComponentCards: FC<Props> = ({
	components,
	projectId,
	envName,
	onClick,
	showAll,
	selectedConfig,
	workflowPhase,
}: Props) => {
	return (
		<div className="bottom-offset">
			<div className="page-offset display-flex">
				{showAll ? null : (
					<>
						<div>
							<ZText.Body className="color-gray" size="20" lineHeight="18" weight="bold">
								Environment
							</ZText.Body>
							<h5 className="color-gray">{envName}</h5>
						</div>
						<div>
							<ZText.Body className="color-gray" size="20" lineHeight="18" weight="bold">
								Est. Monthly Cost
							</ZText.Body>
							<h5 className="color-gray">{renderCostEnv(projectId, envName)}</h5>
						</div>
					</>
				)}
			</div>
			<div className="com-cards border">
				{components.map((config: Component) => (
					<ConfigCard
						key={config.name}
						config={config}
						showAll={showAll}
						onClick={onClick}
						isSelected={selectedConfig === config}
						workflowPhase={workflowPhase}
					/>
				))}
			</div>
		</div>
	);
};

export const ConfigCard: FC<EnvironmentComponentItemProps> = ({
	config,
	showAll,
	onClick,
	isSelected,
	workflowPhase,
}: EnvironmentComponentItemProps) => {
	const [componentStatus, setComponentStatus] = useState<ZSyncStatus>(ZSyncStatus.Unknown);
	// useEffect(() => {
	// 	const delayedStatus = [ZSyncStatus.Destroyed, ZSyncStatus.Provisioned, ZSyncStatus.InSync];
	// 	if (delayedStatus.includes(config.componentStatus) && isSelected) {
	// 		workflowPhase === 'Succeeded' && setComponentStatus(config.componentStatus);
	// 	} else {
	// 		setComponentStatus(config.componentStatus);
	// 	}
	// }, [config?.componentStatus, workflowPhase]);

	return (
		<ZModelCard
			classNames={`component-card ${
				// config.componentStatus === ZSyncStatus.Destroyed ? 'destroyed' : ''
				''
			}`}
			key={config.name}
			model="Environment Component"
			teamName={''}
			envName={''}
			estimatedCost={"-1"}
			title={config.name}
			items={mapGridItems(config, componentStatus)}
			labels={getLabels(config)}
			onClick={(): void => onClick(config)}
		/>
	);
};

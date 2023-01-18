import './style.scss';

import { ZText } from 'components/atoms/text/Text';
import { ZModelCard } from 'components/molecules/cards/Card';
import { CostRenderer, renderHealthStatus, renderLabels, renderSyncedStatus } from 'components/molecules/cards/renderFunctions';
import { ZSyncStatus } from 'models/argo.models';
import { EnvironmentComponentItem, EnvironmentComponentsList } from 'models/projects.models';
import React, { FC, useEffect, useState } from 'react';
import { EntityStore } from 'models/entity.store';
import { Component, Environment } from 'models/entity.type';

type Props = {
	components: Component[];
	projectId: string;
	env?: Environment;
	onClick: Function;
	showAll?: boolean;
	selectedConfig?: Component;
	workflowPhase?: string;
};

type EnvironmentComponentItemProps = {
	config: Component;
	onClick: Function;
	showFullName?: boolean | undefined;
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

const mapGridItems = (component: Component) => {
	return (
		<>
			{renderSyncedStatus(
				component.status as ZSyncStatus
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
	env,
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
							<h5 className="color-gray">{env?.argoId}</h5>
						</div>
						<div>
							<ZText.Body className="color-gray" size="20" lineHeight="18" weight="bold">
								Est. Monthly Cost
							</ZText.Body>
							<h5 className="color-gray">{<CostRenderer data={env?.estimatedCost} />}</h5>
						</div>
					</>
				)}
			</div>
			<div className="com-cards border">
				{components.map((config: Component) => (
					<ConfigCard
						key={config.argoId}
						config={config}
						onClick={onClick}
						showFullName={showAll}
					/>
				))}
			</div>
		</div>
	);
};

export const ConfigCard: FC<EnvironmentComponentItemProps> = ({
	config,
	onClick,
	showFullName = false,
}: EnvironmentComponentItemProps) => {
	const env = EntityStore.getInstance().getEnvironmentById(config.envId);
	const team = EntityStore.getInstance().getTeam(env?.teamId || -1)

	return (
		<ZModelCard
			classNames={`component-card ${
				config.status === ZSyncStatus.Destroyed ? 'destroyed' : ''
			}`}
			key={config.argoId}
			model="Environment Component"
			teamName={team?.name || ''}
			envName={env?.name || ''}
			estimatedCost={<CostRenderer data={config.estimatedCost} />}
			title={showFullName ? config.argoId : config.name}
			items={mapGridItems(config)}
			labels={getLabels(config)}
			onClick={(): void => onClick(config)}
		/>
	);
};

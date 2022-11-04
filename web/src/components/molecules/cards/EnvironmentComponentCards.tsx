import './style.scss';

import { ZText } from 'components/atoms/text/Text';
import { ZModelCard } from 'components/molecules/cards/Card';
import { renderHealthStatus, renderLabels, renderSyncedStatus } from 'components/molecules/cards/renderFunctions';
import { ZSyncStatus } from 'models/argo.models';
import { EnvironmentComponentItem, EnvironmentComponentsList } from 'models/projects.models';
import { renderCost } from 'pages/authorized/environment-components/helpers';
import { renderCost as renderCostEnv } from 'pages/authorized/environments/helpers';
import React, { FC, useEffect, useState } from 'react';

type Props = {
	components: EnvironmentComponentsList;
	projectId: string;
	envName: string;
	onClick: Function;
	showAll?: boolean;
	selectedConfig?: EnvironmentComponentItem;
	workflowPhase?: string;
};

type EnvironmentComponentItemProps = {
	config: EnvironmentComponentItem;
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

export const filterLabels = (config: any): { [name: string]: string } => {
	const { labels } = config; 
	labels.team_id = labels.project_id;
	labels.dependsOn = config.dependsOn.toString();
	const HIDDEN_KEYS = new Set(['component_status', 'component_name', 'component_cost', 'is_destroy', 'audit_status']);
	return Object.keys(labels).reduce((obj: any, key) => {
		if (labels[key] && !HIDDEN_KEYS.has(key) && !key.startsWith('depends_on_')) {
			obj[key] = labels[key];
		}
		return obj;
	}, {});
};

const totalCost = (components: EnvironmentComponentsList): string => {
	let cost = 0;

	{
		components.map((config: EnvironmentComponentItem) => (cost += parseFloat(config.componentCost)));
	}

	return cost.toFixed(3).toString();
};

const mapGridItems = (component: EnvironmentComponentItem, componentStatus: ZSyncStatus) => {
	return (
		<>
			{renderSyncedStatus(
				componentStatus,
				component.operationPhase,
				component.runningStatus,
				component.syncFinishedAt,
				component
			)}
		</>
	);
};

const getLabels = (component: EnvironmentComponentItem): any => {
	return component.labels ? renderLabels(filterLabels(component)) : <></>;
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
				{components.map((config: EnvironmentComponentItem) => (
					<ConfigCard
						key={getFullName(
							config.labels?.project_id,
							config.labels?.environment_id,
							config.componentName
						)}
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
	useEffect(() => {
		const delayedStatus = [ZSyncStatus.Destroyed, ZSyncStatus.Provisioned, ZSyncStatus.InSync];
		if (delayedStatus.includes(config.componentStatus) && isSelected) {
			workflowPhase === 'Succeeded' && setComponentStatus(config.componentStatus);
		} else {
			setComponentStatus(config.componentStatus);
		}
	}, [config?.componentStatus, workflowPhase]);

	return (
		<ZModelCard
			classNames={`component-card ${config.componentStatus === ZSyncStatus.Destroyed ? 'destroyed' : ''}`}
			key={config.id}
			model="Environment Component"
			teamName={config.labels?.project_id || ''}
			envName={getEnvironmentName(config.labels?.project_id, config.labels?.environment_id)}
			estimatedCost={renderCost(config.id)}
			title={
				showAll
					? getFullName(config.labels?.project_id, config.labels?.environment_id, config.componentName)
					: config.componentName
			}
			items={mapGridItems(config, componentStatus)}
			labels={getLabels(config)}
			onClick={(): void => onClick(config)}
		/>
	);
};

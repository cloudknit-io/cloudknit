import { filterLabels } from 'components/molecules/cards/EnvironmentComponentCards';
import { CostRenderer, renderLabels, renderSyncedStatus } from 'components/molecules/cards/renderFunctions';
import { ZSyncStatus } from 'models/argo.models';
import { Component } from 'models/entity.store';
import { EnvironmentComponentItem } from 'models/projects.models';
import React, { ReactNode, useEffect, useState } from 'react';
import { getSeparatedConfigId } from '../helpers';

export type Props = {
	config: Component;
};
export const ConfigWorkflowLeftView: React.FC<Props> = ({ config }: Props) => {
	const [cost, setCost] = useState<JSX.Element>(<></>);
	const [envName, setEnvName] = useState<string>('');
	const [teamName, setTeamName] = useState<string>('');
	const [operation, setOperation] = useState<string>('');

	useEffect(() => {
		if (!config) return;
		const names = getSeparatedConfigId(config);
		setTeamName(names.team || '');
		setEnvName(names.environment);
		setCost(<CostRenderer data={config.estimatedCost} />);
		setOperation(config.isDestroyed ? 'destroy' : 'provision');
	}, [config]);

	return (
		<div className="labels">
			<div className="config-info">
				<div>
					{<span>Team:</span>} {teamName}
				</div>
				<div>
					{<span>Environment:</span>} {envName}
				</div>
				<div>
					{<span>Est. Monthly Cost:</span>} {cost}
				</div>
				<div>
					{<span>Operation:</span>} <span className="capitalize-text">{operation}</span>
				</div>
				<div>{renderSyncedStatus(config.status as ZSyncStatus, '', '', '', config)}</div>
			</div>
			<div>
				{renderLabels({
					teamName,
					envName,
					dependsOn: (config.dependsOn?.length ? config.dependsOn : [envName]).join(','),
				})}
			</div>
		</div>
	);
};

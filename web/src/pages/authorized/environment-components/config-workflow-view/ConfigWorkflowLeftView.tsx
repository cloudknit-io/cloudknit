import { filterLabels } from 'components/molecules/cards/EnvironmentComponentCards';
import { CostRenderer, renderLabels, renderSyncedStatus } from 'components/molecules/cards/renderFunctions';
import { ZSyncStatus } from 'models/argo.models';
import { EnvironmentComponentItem } from 'models/projects.models';
import React, { ReactNode, useEffect, useState } from 'react';

export type Props = {
    configLabels?: any;
    config: EnvironmentComponentItem;
};
export const ConfigWorkflowLeftView: React.FC<Props> = ({ configLabels, config }: Props) => {
	const [cost, setCost] = useState<JSX.Element>(<></>);
	const [labels, setLabels] = useState<ReactNode[]>([]);
	const [envName, setEnvName] = useState<string>('');
	const [teamName, setTeamName] = useState<string>('');
	const [operation, setOperation] = useState<string>('');

	useEffect(() => {
		if (!config || !configLabels) return;
		setTeamName(configLabels.project_id || '');
		setEnvName(configLabels.environment_id?.replace(configLabels.project_id + '-', ''));
		setCost(<CostRenderer data={config.componentCost}/>);
		setOperation(config.isDestroy ? 'destroy' : 'provision');
		setLabels(renderLabels(filterLabels(config)));
	}, [configLabels, config]);

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
				<div>{renderSyncedStatus(config.componentStatus, '', '', '', config)}</div>
			</div>
			<div>{labels}</div>
		</div>
	);
};

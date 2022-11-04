import './style.scss';

import { filterLabels } from 'components/molecules/cards/EnvironmentComponentCards';
import { renderHealthStatus, renderLabels, renderSyncedStatus } from 'components/molecules/cards/renderFunctions';
import { ZSyncStatus } from 'models/argo.models';
import { EnvironmentComponentItem } from 'models/projects.models';
import React, { FC, useContext, useEffect, useState } from 'react';
import { ArgoComponentsService } from 'services/argo/ArgoComponents.service';
import { eventColumns } from '../helpers';
import { ZTable } from 'components/atoms/table/Table';
import { ConfigWorkflowLeftView } from './ConfigWorkflowLeftView';

type Props = {
	projectId: string;
	environmentId: string;
	config: EnvironmentComponentItem;
};

export const ConfigWorkflowViewApplication: FC<Props> = (props: Props) => {
	const { projectId, environmentId, config } = props;
	const [events, setEvents] = useState<[]>([]);

	useEffect(() => {
		if (!config.id) return;
		ArgoComponentsService.getApplicationEvents(config.id).then(({ data }) => {
			if (data.items?.length > 0) {
				setEvents(data.items);
			}
		});
	}, [config]);

	const getTabs = () => {
		return (
			<nav className="nav-tab">
				<ul>
					<li key={'events'} className={`nav-tab_item nav-tab_item$--active`}>
						<a onClick={() => {}}>Events</a>
					</li>
				</ul>
			</nav>
		);
	};

	const getView = () => {
		return (
			<ZTable
				table={{
					columns: eventColumns,
					rows: events,
				}}
				rowHeight={40}
				rowConditionalClass={(data: any) => {
					if (data?.type === 'Normal') {
						return 'zlifecycle-event-table-row zlifecycle-event-table-row-success';
					}

					if (data) {
						return 'zlifecycle-event-table-row zlifecycle-event-table-row-failed';
					}
					return '';
				}}
			/>
		);
	};

	return (
		<>
			<div className="zlifecycle-config-workflow-view zscrollbar">
				<div className="zlifecycle-config-workflow-view__config-info">
					<div className="zlifecycle-config-workflow-view__header">
						<p className="heading">
							<span>{config.componentName || config.name}</span>
						</p>
					</div>
					<ConfigWorkflowLeftView config={config} key={config.id} configLabels={config.labels}/>
				</div>
				{
					<div className="zlifecycle-config-workflow-view__diagram">
						{getTabs()}{' '}
						<div style={{ overflowY: 'auto', height: 'calc(100vh - 110px)', paddingRight: '20px' }}>
							{getView()}
						</div>
					</div>
				}
			</div>
		</>
	);
};

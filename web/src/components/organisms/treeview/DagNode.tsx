import { CostRenderer, getSyncStatusIcon, renderSyncedStatus } from 'components/molecules/cards/renderFunctions';
import { ZSyncStatus } from 'models/argo.models';
import { getTime } from 'pages/authorized/environment-components/helpers';
import React from 'react';
import './tree-view-new.scss';
import { getClassName } from './tree-view.helper';

export type DagProps = {
	data: {
		icon: JSX.Element;
		name: string;
		cost: number;
		status: ZSyncStatus;
		timestamp: Date;
		operation: 'Provision' | 'Destroy',
		isSkipped: boolean;
	};
};

export const DagNode: React.FC<DagProps> = ({ data }) => {
	const { icon, name, cost, status, timestamp, operation, isSkipped } = data;
	return (
	<div className={`dag-node pod${getClassName(status || '')} ${isSkipped ? 'striped' : ''}`}>
			<div className="dag-node__icon">{icon}</div>
			<div className="dag-node__info">
				<div className="dag-node__info--name">{name}</div>
				<div className="dag-node__info--cost"><CostRenderer data={cost} /></div>
				<div className="dag-node__info--status">
					<div className="dag-node__info--status--icon">{getSyncStatusIcon(status, operation)}</div>
					<div className="dag-node__info--status--timestamp">&nbsp;|&nbsp;{getTime(timestamp.toString())}</div>
				</div>
			</div>
			<div className='dag-node__tooltip'>
				{renderSyncedStatus(status)}
			</div>
		</div>
	);
};

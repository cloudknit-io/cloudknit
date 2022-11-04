import './style.scss';

import { ReactComponent as InProcessIcon } from 'assets/images/icons/DAG/In-process.svg';
import { ReactComponent as PendingIcon } from 'assets/images/icons/DAG/Pending.svg';
import { ReactComponent as StartIcon } from 'assets/images/icons/DAG/Start.svg';
import { ReactComponent as SyncedIcon } from 'assets/images/icons/DAG/Synced.svg';
import { ReactComponent as ClearIcon } from 'assets/images/icons/field/clear.svg';
import classNames from 'classnames';
import { ZText } from 'components/atoms/text/Text';
import React, { FC, ReactNode } from 'react';

// TODO (T.P) operation phase
export type NodeStatus = 'Pending' | 'Disregarded' | 'Succeeded' | 'Failed' | 'InProcess' | 'Mutated' | 'Skipped';
export type NodeIcon = 'Start' | 'Pending' | 'InProcess' | 'Synced' | 'Failed' | 'Mutated' | 'Skipped';

type Props = {
	status: NodeStatus;
	icon: NodeIcon;
	text: string;
};

const getNodeClass = (status: NodeStatus): string => {
	switch (status) {
		case 'Disregarded':
			return 'zlifecycle-diagram-node__node--disregarded';
		case 'Failed':
			return 'zlifecycle-diagram-node__node--failed';
		case 'Pending':
			return 'zlifecycle-diagram-node__node--pending';
		case 'Succeeded':
			return 'zlifecycle-diagram-node__node--successful';
		case 'InProcess':
			return 'zlifecycle-diagram-node__node--in-process';
		case 'Mutated':
			return 'zlifecycle-diagram-node__node--mutated';
		default:
			return '';
	}
};

const getNodeIcon = (status: NodeIcon): ReactNode => {
	switch (status) {
		case 'Start':
			return <StartIcon />;
		case 'Failed':
			return <ClearIcon />;
		case 'Pending':
			return <PendingIcon />;
		case 'Synced':
		case 'Skipped':
			return <SyncedIcon />;
		case 'InProcess':
			return <InProcessIcon />;
		case 'Mutated':
			return <SyncedIcon className="mutated" />;
		default:
			return '';
	}
};

export const ZDiagramNode: FC<Props> = ({ status, icon, text }: Props) => {
	return (
		<div className="zlifecycle-diagram-node">
			<div className={classNames('zlifecycle-diagram-node__node', getNodeClass(status))}>{getNodeIcon(icon)}</div>
			<ZText.Body className="zlifecycle-diagram-node__text" size="12" lineHeight="16">
				{text}
			</ZText.Body>
		</div>
	);
};

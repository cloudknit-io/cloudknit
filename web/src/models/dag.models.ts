import { ReactElement } from 'react';

import { SyncStatusCode, ZSyncStatus } from './argo.models';

export interface DagNode {
	name: string;
	icon: ReactElement<SVGSVGElement>;
	status: ZSyncStatus;
	children?: DagNode[];
}

export interface Point {
	x: number;
	y: number;
}

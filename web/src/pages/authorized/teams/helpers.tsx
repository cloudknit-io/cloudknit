import { TableColumn } from 'components/atoms/table/Table';
import { CostRenderer, renderHealthStatus, renderSyncedStatus } from 'components/molecules/cards/renderFunctions';
import { ZStreamRenderer } from 'components/molecules/zasync-renderer/ZStreamRenderer';
import React from 'react';
import { CostingService } from 'services/costing/costing.service';

export const teamTableColumns: TableColumn[] = [
	{
		id: 'name',
		name: 'Name',
		// width: 250,
	},
	{
		id: 'resources',
		name: 'Environments',
		render: (data: any) => data?.length,
	},
	{
		id: 'repoUrl',
		name: 'Repo Url',
	}
];

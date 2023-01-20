import { TableColumn } from 'components/atoms/table/Table';

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

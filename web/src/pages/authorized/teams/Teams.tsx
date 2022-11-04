import { ZLoaderCover } from 'components/atoms/loader/LoaderCover';
import { ZTable } from 'components/atoms/table/Table';
import { TeamCards } from 'components/molecules/cards/TeamCards';
import { useApi } from 'hooks/use-api/useApi';
import { HealthStatusCode, SyncStatusCode, SyncStatuses } from 'models/argo.models';
import { TeamItem, TeamsList } from 'models/projects.models';
import React, { useEffect, useState } from 'react';
import { ArgoTeamsService } from 'services/argo/ArgoProjects.service';

import { usePageHeader } from '../contexts/EnvironmentHeaderContext';
import { getCheckBoxFilters, renderHealthStatusItems, renderSyncStatusItems } from '../environments/helpers';
import { teamTableColumns } from './helpers';

export const Teams: React.FC = () => {
	const { pageHeaderObservable, breadcrumbObservable } = usePageHeader();
	const [loading, setLoading] = useState<boolean>(true);
	const [query, setQuery] = useState<string>('');
	const [syncStatusFilter, setSyncStatusFilter] = useState<Set<SyncStatusCode>>(new Set<SyncStatusCode>());
	const [healthStatusFilter, setHealthStatusFilter] = useState<Set<HealthStatusCode>>(new Set<HealthStatusCode>());
	const [filterDropDownOpen, toggleFilterDropDown] = useState(false);
	const [checkBoxFilters, setCheckBoxFilters] = useState<JSX.Element>(<></>);
	const fetchTeams = useApi(ArgoTeamsService.getProjects).fetch;
	const [teams, setTeams] = useState<TeamsList>([]);
	const [viewType, setViewType] = useState<string>('');
	const [filterItems, setFilterItems] = useState<Array<() => JSX.Element>>([]);

	useEffect(() => {
		fetchTeams().then(({ data }) => {
			if (data) {
				setTeams(data);
			}
			setLoading(false);
		});
	}, [fetchTeams]);

	const setQueryValue = (queryLoc: string): void => {
		setQuery(queryLoc.toLowerCase());
	};

	const setToggleFilterDropDownValue = (val: boolean) => {
		toggleFilterDropDown(val);
	};

	const syncStatusMatch = (item: TeamItem): boolean => {
		return syncStatusFilter.has(item.syncStatus);
	};

	const healthStatusMatch = (item: TeamItem): boolean => {
		return healthStatusFilter.has(item.healthStatus);
	};

	const getFilteredData = (): TeamsList => {
		let filteredItems = [...teams];
		if (syncStatusFilter.size > 0) {
			filteredItems = [...filteredItems.filter(syncStatusMatch)];
		}

		if (healthStatusFilter.size > 0) {
			filteredItems = [...filteredItems.filter(healthStatusMatch)];
		}

		return filteredItems.filter(item => {
			return item.name.toLowerCase().includes(query);
		});
	};

	useEffect(() => {
		setFilterItems([
			renderSyncStatusItems
				.bind(null, SyncStatuses, syncStatusFilter, setSyncStatusFilter, 'Team Status')
				.bind(null, (status: string) => teams.filter(e => e.syncStatus === status).length),
		]);
	}, [teams, syncStatusFilter, healthStatusFilter]);

	useEffect(() => {
		setCheckBoxFilters(getCheckBoxFilters(filterDropDownOpen, setToggleFilterDropDownValue, filterItems));
	}, [filterItems, filterDropDownOpen, filterDropDownOpen]);

	useEffect(() => {
		pageHeaderObservable.next({
			breadcrumbs: [],
			headerTabs: [],
			pageName: 'Teams',
			filterTitle: '',
			onSearch: setQueryValue,
			onViewChange: setViewType,
			buttonText: '',
			checkBoxFilters: checkBoxFilters,
		});

		breadcrumbObservable.next(false);
	}, [teams, checkBoxFilters]);

	const renderView = () => {
		switch (viewType) {
			case 'list':
				return (
					<div className="zlifecycle-table">
						<ZTable table={{ columns: teamTableColumns, rows: getFilteredData() }} />
					</div>
				);
			default:
				return <TeamCards teams={getFilteredData() || teams} />;
		}
	};

	return (
		<div className="zlifecycle-page">
			<ZLoaderCover loading={loading}>
				<section className="dashboard-content">{renderView()}</section>
			</ZLoaderCover>
		</div>
	);
};

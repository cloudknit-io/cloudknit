import { ZLoaderCover } from 'components/atoms/loader/LoaderCover';
import { ZTable } from 'components/atoms/table/Table';
import { TeamCards } from 'components/molecules/cards/TeamCards';
import { HealthStatusCode, SyncStatusCode } from 'models/argo.models';
import { EntityStore } from 'models/entity.store';
import { Team } from 'models/entity.type';
import { TeamItem } from 'models/projects.models';
import React, { useEffect, useState } from 'react';
import { useMemo } from 'react';

import { usePageHeader } from '../contexts/EnvironmentHeaderContext';
import { getCheckBoxFilters } from '../environments/helpers';
import { teamTableColumns } from './helpers';

export const Teams: React.FC = () => {
	const entityStore = useMemo(() => EntityStore.getInstance(), []);
	const { pageHeaderObservable, breadcrumbObservable } = usePageHeader();
	const [loading, setLoading] = useState<boolean>(true);
	const [query, setQuery] = useState<string>('');
	const [syncStatusFilter, setSyncStatusFilter] = useState<Set<SyncStatusCode>>(new Set<SyncStatusCode>());
	const [healthStatusFilter, setHealthStatusFilter] = useState<Set<HealthStatusCode>>(new Set<HealthStatusCode>());
	const [filterDropDownOpen, toggleFilterDropDown] = useState(false);
	const [checkBoxFilters, setCheckBoxFilters] = useState<JSX.Element>(<></>);
	const [teams, setTeams] = useState<Team[]>([]);
	const [viewType, setViewType] = useState<string>('');
	const [filterItems, setFilterItems] = useState<Array<() => JSX.Element>>([]);

	useEffect(() => {
		const subscription = entityStore.emitter.subscribe((update) => {
			const teams = update.teams;
			if (teams.length === 0) return;
			setTeams(teams);
			setLoading(false);
		});

		return () => {
			subscription.unsubscribe();
		};
	}, []);

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

	const getFilteredData = (): Team[] => {
		return [];
		// let filteredItems = [...teams];
		// if (syncStatusFilter.size > 0) {
		// 	filteredItems = [...filteredItems.filter(syncStatusMatch)];
		// }

		// if (healthStatusFilter.size > 0) {
		// 	filteredItems = [...filteredItems.filter(healthStatusMatch)];
		// }

		// return filteredItems.filter(item => {
		// 	return item.name.toLowerCase().includes(query);
		// });
	};

	// useEffect(() => {
	// 	setFilterItems([
	// 		renderSyncStatusItems
	// 			.bind(null, SyncStatuses, syncStatusFilter, setSyncStatusFilter, 'Team Status')
	// 			.bind(null, (status: string) => teams.filter(e => e.syncStatus === status).length),
	// 	]);
	// }, [teams, syncStatusFilter, healthStatusFilter]);

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
				return <TeamCards teams={teams} />;
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

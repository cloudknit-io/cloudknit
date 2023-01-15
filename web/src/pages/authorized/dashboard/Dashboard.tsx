import './dashboard.styles.scss';

import { ReactComponent as OpenInFull } from 'assets/images/icons/open_in_full.svg';
import { ZLoaderCover } from 'components/atoms/loader/LoaderCover';
import { streamMapper } from 'helpers/streamMapper';
import { useApi } from 'hooks/use-api/useApi';
import { ApplicationWatchEvent } from 'models/argo.models';
import { PageHeaderTabs, TeamItem, TeamsList } from 'models/projects.models';
import React, { FC, useCallback, useEffect, useState } from 'react';
import { Subscription } from 'rxjs';
import { ArgoTeamsService } from 'services/argo/ArgoProjects.service';
import { subscriber } from 'utils/apiClient/EventClient';

import { usePageHeader } from '../contexts/EnvironmentHeaderContext';
import { d3Charts } from './helpers';

export const Dashboard: FC = () => {
	const [loading, setLoading] = useState<boolean>(true);
	const [projects, setTeams] = useState<TeamsList>([]);
	const [hierarchicalData, setHierarchicalData] = useState<any>([]);
	const [streamData, setStreamData] = useState<ApplicationWatchEvent | null>(null);
	const [query, setQuery] = useState<string>('');
	const [componentData, setComponentData] = useState<any>(null);
	const headerTabs: PageHeaderTabs = [];
	const [viewType, setViewType] = useState<string>('');
	const { pageHeaderObservable, breadcrumbObservable } = usePageHeader();
	const breadcrumbItems = [
		{
			path: '/dashboard',
			name: 'Dashboard',
			active: true,
		},
	];

	useEffect(() => {
		const $subscription: Subscription = subscriber.subscribe(response => {
			setStreamData(response);
		});

		return (): void => $subscription.unsubscribe();
	}, []);

	// useEffect(() => {
	// 	const newItems = streamMapper<TeamItem>(streamData, projects, ArgoMapper.parseTeam, 'project');
	// 	setTeams(newItems);
	// }, [streamData]);

	// useEffect(() => {
	// 	const promises = [fetch(), getEnvironments(null), getAllComponents()];
	// 	Promise.all(promises).then((data: any) => {
	// 		if (data[0].data && data[1].data && data[2].data) {
	// 			const comps = data[2].data.map((e: any) => ({ ...e, value: 1 }));
	// 			const env = data[1].data.map((e: any) => ({
	// 				...e,
	// 				children: comps.filter((c: any) => (c.labels ? c.labels['environment_id'] === e.id : false)),
	// 			}));
	// 			setHierarchicalData(
	// 				data[0].data.map((p: any) => ({
	// 					...p,
	// 					children: env.filter((e: any) => (e.labels ? e.labels['project_id'] === p.id : false)),
	// 				}))
	// 			);
	// 		}
	// 	});
	// }, [fetch, getEnvironments]);

	const setQueryValue = (queryLoc: string): void => setQuery(queryLoc);

	const filteredProjects = useCallback((): TeamItem[] => {
		return projects.filter(item => item?.name?.toLowerCase().includes(query));
	}, [query, projects]);

	// useEffect(() => {
	// 	getAllComponents().then(({ data }) => {
	// 		if (data) {
	// 			setComponentData(data);
	// 			setLoading(false);
	// 		}
	// 	});
	// }, [getAllComponents]);

	useEffect(() => {
		pageHeaderObservable.next({
			breadcrumbs: breadcrumbItems,
			headerTabs,
			pageName: 'Dashboard',
			filterTitle: 'Filter by team:',
			onSearch: setQueryValue,
			onViewChange: setViewType,
			buttonText: 'New environment',
			checkBoxFilters: <></>,
		});

		breadcrumbObservable.next(false);
	}, []);

	return (
		<div className="zlifecycle-page">
			<ZLoaderCover loading={loading}>
				<div
					className="d3-container"
					onClick={e => {
						const target = (e.target as any).closest('button.view-toggle');
						if (target) {
							const container = e.currentTarget;
							const graph = target.closest('div.d3-graph');
							if (container && graph) {
								const elements = Array.from(container.querySelectorAll('div.d3-graph'));
								if (container.classList.contains('maximize')) {
									elements.forEach(e => {
										e.classList.remove('hide-d3');
									});
								} else {
									elements.forEach(e => {
										if (graph !== e) e.classList.add('hide-d3');
									});
								}
							}
							e.currentTarget.classList.toggle('maximize');
						}
					}}>
					{d3Charts(hierarchicalData, componentData).map(e => (
						<div className="d3-graph" id={e.id}>
							<div>
								<label>{e.label}</label>
								<button className="view-toggle">
									<OpenInFull />
								</button>
							</div>
							{e.jsx}
						</div>
					))}
				</div>
			</ZLoaderCover>
		</div>
	);
};

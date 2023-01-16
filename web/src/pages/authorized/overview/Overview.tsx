import { EntityStore } from 'models/entity.store';
import { Update } from 'models/entity.type';
import React, { useEffect, useMemo, useState } from 'react';
import { breadcrumbObservable, pageHeaderObservable } from '../contexts/EnvironmentHeaderContext';
import { CircularClusterPacking } from '../dashboard/CircularClusterPacking';
import '../dashboard/dashboard.styles.scss';

export const Overview: React.FC = () => {
	const entityStore = useMemo(() => EntityStore.getInstance(), []);
	const [hierarchicalData, setHierarchicalData] = useState<any>([]);

	useEffect(() => {
		entityStore.emitter.subscribe((val: Update) => {
			const envs = val.environments;
			const teams = val.teams;
			const data = teams.map(t => ({
				name: t.name,
				data: t,
				children: envs
					.filter(e => e.teamId === t.id)
					.map(e => ({
						name: e.argoId,
						data: e,
						children: e.dag.map(d => ({
							name: `${e.argoId}-${d.name}`,
							// asyncData: PromisentityStore.getComponents
							value: 1,
						})),
					})),
			}));
			setHierarchicalData(data);
		});
	}, []);

	useEffect(() => {
		pageHeaderObservable.next({
			breadcrumbs: [],
			headerTabs: [],
			pageName: null,
			filterTitle: '',
			onSearch: () => {},
			buttonText: '',
			onViewChange: () => {},
		});
	});

	useEffect(() => {
		breadcrumbObservable.next(false);
	}, [breadcrumbObservable]);
	return (
		<div id="cluster" className="graph-container">
			<CircularClusterPacking data={hierarchicalData} />
		</div>
	);
};

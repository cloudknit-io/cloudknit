import { ZLoaderCover } from 'components/atoms/loader/LoaderCover';
import { ZText } from 'components/atoms/text/Text';
import { EntityStore } from 'models/entity.store';
import { Update } from 'models/entity.type';
import React, { useEffect, useMemo, useState } from 'react';
import { breadcrumbObservable, pageHeaderObservable } from '../contexts/EnvironmentHeaderContext';
import { CircularClusterPacking } from '../dashboard/CircularClusterPacking';
import '../dashboard/dashboard.styles.scss';

export const Overview: React.FC = () => {
	const entityStore = useMemo(() => EntityStore.getInstance(), []);
	const [hierarchicalData, setHierarchicalData] = useState<any>([]);
	const [loading, setLoading] = useState<boolean>(true);

	useEffect(() => {
		entityStore.emitter.subscribe((val: Update) => {
			if (!entityStore.AllDataFetched) return;
			const envs = val.environments;
			const teams = val.teams;
			const components = val.components;
			const data = teams.map(t => ({
				name: t.name,
				data: t,
				children: envs
					.filter(e => e.teamId === t.id)
					.map(e => ({
						name: e.argoId,
						data: e,
						children: components
							.filter(c => c.envId === e.id)
							.map(c => ({
								name: c.argoId,
								data: c,
								value: 1,
							})),
					})),
			}));
			setHierarchicalData(data);
			setLoading(false);
		});
		Promise.resolve(entityStore.getTeams(true));
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
		<>
			<ZText.Body className="page-offset" size="36" weight="bold">
				Environments Overview
			</ZText.Body>
			<div id="cluster" className="graph-container">
				<ZLoaderCover loading={loading}>
					<CircularClusterPacking data={hierarchicalData} />
				</ZLoaderCover>
			</div>
		</>
	);
};

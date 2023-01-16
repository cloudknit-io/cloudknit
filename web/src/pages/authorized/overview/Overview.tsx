import { EntityStore, Update } from "models/entity.store";
import React, { useEffect, useMemo, useState } from "react";
import { breadcrumbObservable, pageHeaderObservable } from "../contexts/EnvironmentHeaderContext";
import { CircularClusterPacking } from "../dashboard/CircularClusterPacking";
import '../dashboard/dashboard.styles.scss';

export const Overview: React.FC = () => {
	const entityStore =  useMemo(() => EntityStore.getInstance(), []);
    const [hierarchicalData, setHierarchicalData] = useState<any>([]);

    useEffect(() => {
		entityStore.emitter.subscribe((val: Update) => {
			const envs = val.environments;
			const teams = val.teams;
			const components = val.components;
			const data = teams.map(t => ({
				name: t.name,
				data: t,
				children: envs.filter(e => e.teamId === t.id).map(e => ({
					name: e.argoId,
					data: e,
					children: components.filter(c => c.envId === e.id).map(c => ({
						name: c.argoId,
						data: c,
						value: 1
					}))
				}))
			}))
			setHierarchicalData(data);
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
    return <div id="cluster" className="graph-container"><CircularClusterPacking data={hierarchicalData}  /></div>;
}
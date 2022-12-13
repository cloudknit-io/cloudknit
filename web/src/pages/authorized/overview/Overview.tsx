import { useApi } from "hooks/use-api/useApi";
import React, { useEffect, useState } from "react";
import { ArgoComponentsService } from "services/argo/ArgoComponents.service";
import { ArgoEnvironmentsService } from "services/argo/ArgoEnvironments.service";
import { ArgoTeamsService } from "services/argo/ArgoProjects.service";
import { breadcrumbObservable, pageHeaderObservable } from "../contexts/EnvironmentHeaderContext";
import { CircularClusterPacking } from "../dashboard/CircularClusterPacking";
import '../dashboard/dashboard.styles.scss';

export const Overview: React.FC = () => {
    const { fetch } = useApi(ArgoTeamsService.getProjects);
	const { fetch: getEnvironments } = useApi(ArgoEnvironmentsService.getEnvironments);
	const { fetch: getAllComponents } = useApi(ArgoComponentsService.getAllComponents);
    const [hierarchicalData, setHierarchicalData] = useState<any>([]);

    useEffect(() => {
		const promises = [fetch(), getEnvironments(null), getAllComponents()];
		Promise.all(promises).then((data: any) => {
			if (data[0].data && data[1].data && data[2].data) {
				const comps = data[2].data.map((e: any) => ({ ...e, value: 1 }));
				const env = data[1].data.map((e: any) => ({
					...e,
					children: comps.filter((c: any) => (c.labels ? c.labels['environment_id'] === e.id : false)),
				}));
				setHierarchicalData(
					data[0].data.map((p: any) => ({
						...p,
						children: env.filter((e: any) => (e.labels ? e.labels['project_id'] === p.id : false)),
					}))
				);
			}
		});
	}, [fetch, getEnvironments]);

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
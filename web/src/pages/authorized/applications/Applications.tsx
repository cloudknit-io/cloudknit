import { ZLoaderCover } from 'components/atoms/loader/LoaderCover';
import { ZModelCard } from 'components/molecules/cards/Card';
import { environmentName } from 'components/molecules/cards/EnvironmentCards';
import { renderSyncedStatus } from 'components/molecules/cards/renderFunctions';
import { useApi } from 'hooks/use-api/useApi';
import { ZSyncStatus } from 'models/argo.models';
import { ListItem } from 'models/general.models';
import { EnvironmentsList, PageHeaderTabs } from 'models/projects.models';
import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { ArgoComponentsService } from 'services/argo/ArgoComponents.service';
import { ArgoEnvironmentsService } from 'services/argo/ArgoEnvironments.service';

import { usePageHeader } from '../contexts/EnvironmentHeaderContext';
import { allApplications as allApps, applications as apps, ZApplication } from './dummyApps';

export const Applications: React.FC = () => {
	const { fetch } = useApi(ArgoComponentsService.getComponents);
	const fetchEnvironments = useApi(ArgoEnvironmentsService.getEnvironments).fetch;
	const [environments, setEnvironments] = useState<EnvironmentsList>([]);
	const { projectId, environmentId } = useParams();
	const [query, setQuery] = useState<string>('');
	const showAll = environmentId === 'all' && projectId === 'all';
	const applications = showAll ? allApps : apps;

	const headerTabs: PageHeaderTabs = [
		...environments.map(environment => {
			const name: string = environmentName(environment);
			return {
				active: environmentId === environment.id,
				name: name.charAt(0).toUpperCase() + name.slice(1),
				path: `/${environment.id}/apps`,
			};
		}),
	];

	const { pageHeaderObservable, breadcrumbObservable } = usePageHeader();
	useEffect(() => {
		breadcrumbObservable.next(true);
	}, [breadcrumbObservable]);

	useEffect(() => {
		pageHeaderObservable.next({
			breadcrumbs: breadcrumbItems,
			headerTabs,
			pageName: 'Applications',
			filterTitle: 'Filter by environment',
			onSearch: setQueryValue,
			buttonText: '',
			onViewChange: () => {},
		});
	});
	const [loading, setLoading] = useState<boolean>(true);
	const breadcrumbItems = [
		{
			path: '/dashboard',
			name: 'All Environments',
			active: false,
		},
		{
			path: `/${projectId}`,
			name: projectId,
			active: false,
		},
		{
			path: `/${projectId}/${environmentId}`,
			name: environmentId,
			active: true,
		},
	];

	const mapGridItems = (app: ZApplication) => {
		return <>{renderSyncedStatus(ZSyncStatus.InSync, 'foo', 'bar', '2021-04-15T19:59:30Z')}</>;
	};
	const getFilteredData = (): ZApplication[] => {
		return applications.filter(app => {
			return app.name.toLowerCase().includes(query);
		});
	};

	useEffect(() => {
		fetchEnvironments(projectId).then(({ data }) => {
			if (data) {
				setEnvironments(data);
			}
			setLoading(false);
		});
	}, [projectId, fetch]);

	const setQueryValue = (queryLoc: string): void => {
		setQuery(queryLoc.toLowerCase());
	};

	return (
		<div className="zlifecycle-page">
			<ZLoaderCover loading={loading}>
				<section className="dashboard-content">
					<div className="bottom-offset">
						<div className="com-cards border">
							{getFilteredData().map((app: ZApplication) => (
								<ZModelCard
									classNames=""
									key={app.name}
									model="Application"
									teamName={projectId}
									envName={environmentName(environments.find(e => e.id === environmentId))}
									estimatedCost="0"
									title={app.name}
									items={mapGridItems(app)}
									onClick={(): void => {}}
								/>
							))}
						</div>
					</div>
				</section>
			</ZLoaderCover>
		</div>
	);
};

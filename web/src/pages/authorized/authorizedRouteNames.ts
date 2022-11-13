import React from 'react';
import { Dashboard } from './dashboard/Dashboard';
import { EnvironmentBuilder } from './environment-builder/EnvironmentBuilder';
import { EnvironmentComponents } from './environment-components/EnvironmentComponents';
import { Environments } from './environments/Environments';
import { Profile } from './profile/Profile';
import { Teams } from './teams/Teams';
import { FeatureKeys, FeatureRoutes } from './feature_toggle';
import { ComponentResourceTree } from 'components/organisms/tree-view/ComponentResourceTree';
import { QuickStart } from 'pages/authorized/quick-start/QuickStart';

export const PROJECTS_URL = '/dashboard';
const DASHBOARD_URL = '/demo-dashboard';
const PROFILE_URL = '/settings';
const TEAMS_URL = '/teams';
const ENVIRONMENT_BUILDER_URL = '/builder';
const ENVIRONMENTS_URL = '/:projectId';
const INFRA_URL = '/:projectId/:environmentId/infra';
const APPS_URL = '/:projectId/:environmentId/apps';
const RESOURCE_VIEW_URL = '/applications/:componentId/resource-view';
const QUICK_START_URL = '/quick-start';

const urls = [
	{ key: 'QUICK_START_URL', value: QUICK_START_URL },
	{ key: 'ENVIRONMENT_BUILDER_URL', value: ENVIRONMENT_BUILDER_URL },
	{ key: 'TEAMS_URL', value: TEAMS_URL },
	{ key: 'PROFILE_URL', value: PROFILE_URL },
	{ key: 'DASHBOARD_URL', value: DASHBOARD_URL },
	{ key: 'PROJECTS_URL', value: PROJECTS_URL },
	{ key: 'ENVIRONMENTS_URL', value: ENVIRONMENTS_URL },
	{ key: 'INFRA_URL', value: INFRA_URL },
	{ key: 'APPS_URL', value: APPS_URL },
	{ key: 'RESOURCE_VIEW_URL', value: RESOURCE_VIEW_URL },
];

Reflect.ownKeys(FeatureRoutes).forEach(key => {
	if (Reflect.get(FeatureRoutes, key) === false) {
		switch (key) {
			// case FeatureKeys.DASHBOARD:
			// 	{
			// 		const i = urls.findIndex(e => e.key === 'DASHBOARD_URL');
			// 		urls.splice(i, 1);
			// 	}
			// 	break;
			// case FeatureKeys.BUILDER:
			// 	{
			// 		const i = urls.findIndex(e => e.key === 'ENVIRONMENT_BUILDER_URL');
			// 		urls.splice(i, 1);
			// 	}
			// 	break;
			case FeatureKeys.APPLICATIONS:
				{
					const i = urls.findIndex(e => e.key === 'APPS_URL');
					urls.splice(i, 1);
				}
				break;
			case FeatureKeys.QUICK_START:
				{
					const i = urls.findIndex(e => e.key === 'QUICK_START_URL');
					urls.splice(i, 1);
				}
				break;
		}
	}
});

export const routes = urls;
export const privateRouteMap: { [key: string]: React.FC } = {
	QUICK_START_URL: QuickStart,
	ENVIRONMENT_BUILDER_URL: EnvironmentBuilder,
	TEAMS_URL: Teams,
	PROFILE_URL: Profile,
	DASHBOARD_URL: Dashboard,
	PROJECTS_URL: Environments,
	ENVIRONMENTS_URL: Environments,
	INFRA_URL: EnvironmentComponents,
	RESOURCE_VIEW_URL: ComponentResourceTree,
};

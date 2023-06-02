import React from 'react';
// import { Dashboard } from './dashboard/Dashboard';
// import { EnvironmentBuilder } from './environment-builder/EnvironmentBuilder';
import { QuickStart } from 'pages/authorized/quick-start/QuickStart';
import { Dashboard } from './dashboard/Dashboard';
import { EnvironmentBuilder } from './environment-builder/EnvironmentBuilder';
import { EnvironmentComponents } from './environment-components/EnvironmentComponents';
import { Environments } from './environments/Environments';
import { BradAdarshFeatureVisible, FeatureRoutes, playgroundFeatureVisible } from './feature_toggle';
import { Overview } from './overview/Overview';
import { Profile } from './profile/Profile';
import { Teams } from './teams/Teams';
import { TermsAndConditions } from './terms-and-conditions/TermsAndConditons';
import { ENVIRONMENT_VARIABLES } from 'utils/environmentVariables';

export const PROJECTS_URL = '/dashboard';
const DASHBOARD_URL = '/demo-dashboard';
const PROFILE_URL = '/settings';
const TEAMS_URL = '/teams';
const ENVIRONMENT_BUILDER_URL = '/builder';
const ENVIRONMENTS_URL = '/:projectId';
const INFRA_URL = '/:projectId/:environmentName';
const RESOURCE_VIEW_URL = '/applications/:componentId/resource-view';
export const QUICK_START_URL = '/quick-start';
const OVERVIEW_URL = '/overview';
export const ORG_REGISTRATION = '/org-registration';

const urls = [
	{ key: 'ORG_REGISTRATION', value: ORG_REGISTRATION},
  { key: 'OVERVIEW_URL', value: OVERVIEW_URL},
	{ key: 'QUICK_START_URL', value: QUICK_START_URL },
	{ key: 'ENVIRONMENT_BUILDER_URL', value: ENVIRONMENT_BUILDER_URL },
	{ key: 'TEAMS_URL', value: TEAMS_URL },
	{ key: 'PROFILE_URL', value: PROFILE_URL },
	{ key: 'DASHBOARD_URL', value: DASHBOARD_URL },
	{ key: 'PROJECTS_URL', value: PROJECTS_URL },
	{ key: 'ENVIRONMENTS_URL', value: ENVIRONMENTS_URL },
	{ key: 'INFRA_URL', value: INFRA_URL },
	{ key: 'RESOURCE_VIEW_URL', value: RESOURCE_VIEW_URL },
];

if (ENVIRONMENT_VARIABLES.PLAYGROUND_APP && !BradAdarshFeatureVisible()) {
    ['ORG_REGISTRATION', 'OVERVIEW_URL', 'QUICK_START_URL', 'ENVIRONMENT_BUILDER_URL', 'PROFILE_URL', 'RESOURCE_VIEW_URL'].forEach(e => {
        urls.splice(urls.findIndex(u => e === u.key), 1);
    })
}

Reflect.ownKeys(FeatureRoutes).forEach(key => {
	if (Reflect.get(FeatureRoutes, key) === false) {
		switch (key) {
			// Add a case and splice that route if feature flagged.
			// case FeatureKeys.QUICK_START:
		// switch (key) {
			// case FeatureKeys.DASHBOARD:
			// 	{
			// 		const i = urls.findIndex(e => e.key === 'QUICK_START_URL');
			// 		urls.splice(i, 1s);
			// 	}
			// 	break;
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
	OVERVIEW_URL: Overview,
	ORG_REGISTRATION: TermsAndConditions,
};

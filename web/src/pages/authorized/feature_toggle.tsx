import AuthStore from "auth/AuthStore";
import { ENVIRONMENT_VARIABLES } from "utils/environmentVariables";

const showFeatures = (process.env.REACT_APP_ENABLED_FEATURE_FLAGS || '')
	.toString()
	.split('|')
	.map(e => e.replace('show_', '').toUpperCase());

const Features: { [key: string]: boolean } = {
	PAGE_DASHBOARD: false,
	PAGE_BUILDER: false,
	PAGE_APPLICATIONS: false,
	TAB_DETAILED_LOGS: false,
	TAB_STATE_FILE: false,
	HARD_SYNC: false,
	DIFF_CHECKER: false,
	VISUALIZATION: false,
	BLUE_GREEN_DEPLOYMENT: false,
	TERM_AGREEMENT: false,
	QUICK_START: true,
};
showFeatures.forEach(feature => {
	Features[feature] = true;
});

export const FeatureRoutes: { [key: string]: boolean } = {
	PAGE_DASHBOARD: Features.PAGE_DASHBOARD,
	PAGE_BUILDER: Features.PAGE_BUILDER,
	PAGE_APPLICATIONS: Features.PAGE_APPLICATIONS,
	QUICK_START: Features.QUICK_START,
};

export const VisibleFeatures = Features;
export const FeatureKeys = {
	DASHBOARD: 'PAGE_DASHBOARD',
	BUILDER: 'PAGE_BUILDER',
	APPLICATIONS: 'PAGE_APPLICATIONS',
	DETAILED_LOGS: 'TAB_DETAILED_LOGS',
	STATE_FILE: 'TAB_STATE_FILE',
	HARD_SYNC: 'HARD_SYNC',
	DIFF_CHECKER: 'DIFF_CHECKER',
	VISUALIZATION: 'VISUALIZATION',
	BLUE_GREEN_DEPLOYMENT: 'BLUE_GREEN_DEPLOYMENT',
	TERM_AGREEMENT: 'TERM_AGREEMENT',
	QUICK_START: 'QUICK_START',
};

export const featureToggled = (featureKey: string, userBased: boolean = false) => {
	if (userBased) {
		return BradAdarshFeatureVisible() && VisibleFeatures[featureKey];
	}
	return VisibleFeatures[featureKey];
}

export function BradAdarshFeatureVisible() : boolean {
	if (ENVIRONMENT_VARIABLES.PLAYGROUND_APP) {
		return false;
	}
	const user = AuthStore.getUser();

	// sometimes life hands you lemons...
	return ['shahadarsh', 'bradj', 'shashank-cloudknit-io'].includes(user?.username || '');
}

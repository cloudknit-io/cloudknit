import { NotificationsApi } from 'components/argo-core';
import * as H from 'history';

export interface AppContext {
	router: {
		history: H.History;
		route: {
			location: H.Location;
		};
	};
	apis: {
		notifications: NotificationsApi;
	};
	history: H.History;
}

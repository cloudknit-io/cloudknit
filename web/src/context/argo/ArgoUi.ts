import { NavigationApi, NotificationsApi } from 'components/argo-core';
import { AppContext as ArgoAppContext } from 'context/argo/ArgoCore';
import { History } from 'history';
import { TeamItem } from 'models/projects.models';
import * as React from 'react';
import { BehaviorSubject } from 'rxjs';

export type AppContext = ArgoAppContext & {
	apis: {
		notifications: NotificationsApi;
		navigation: NavigationApi;
		baseHref: string;
		projects: TeamItem[];
	};
};

export interface ContextApis {
	notifications?: NotificationsApi;
	navigation?: NavigationApi;
	baseHref: string;
	projects: TeamItem[];
	failedEnvironments: BehaviorSubject<Map<string, any>>;
}
export const Context = React.createContext<(ContextApis & { history: History }) | null>(null);
export const { Provider, Consumer } = Context;

import AuthStore from 'auth/AuthStore';
import { NavigationManager, Notifications, NotificationsManager } from 'components/argo-core';
import { ZLoaderCover } from 'components/atoms/loader/LoaderCover';
import { Provider } from 'context/argo/ArgoUi';
import { createBrowserHistory } from 'history';
import { TeamItem } from 'models/projects.models';
import React, { Suspense, useEffect, useState } from 'react';
import { BrowserRouter as Router } from 'react-router-dom';
import { BehaviorSubject, from } from 'rxjs';
import { EntityService } from 'services/entity/entity.service';
const Routes = React.lazy(() => import('router/Routes'));

const FailedEnvironments = new BehaviorSubject<Map<string, any>>(new Map());

const App: React.FC = () => {
	const [loading, setLoading] = useState(true);
	const [projects, setTeams] = useState<TeamItem[]>([]);
	const bases = document.getElementsByTagName('base');
	const base = bases.length > 0 ? bases[0].getAttribute('href') || '/' : '/';
	const history = createBrowserHistory({ basename: base });
	const notificationsManager: NotificationsManager = new NotificationsManager();

	useEffect(() => {
		const $subscription = from(AuthStore.refresh()).subscribe(profile => {
			setLoading(false);
			if (!profile.selectedOrg) {
				return;
			}
			EntityService.getInstance();
		});

		return () => $subscription.unsubscribe();
	}, [fetch]);

	return (
		<>
			<Provider
				value={{
					history,
					baseHref: base,
					navigation: history && new NavigationManager(history),
					notifications: notificationsManager,
					projects: projects,
					failedEnvironments: FailedEnvironments,
				}}>
				<Router>
					<Suspense fallback={<ZLoaderCover loading={loading} fullScreen />}>
						<Routes />
					</Suspense>
				</Router>
			</Provider>
			<Notifications notifications={notificationsManager.notifications} />
		</>
	);
};

export default App;

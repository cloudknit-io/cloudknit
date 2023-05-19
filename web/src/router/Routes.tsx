import Anonymous from 'pages/anonymous';
import { LOGIN_URL } from 'pages/anonymous/anonymousRouteNames';
import { Login } from 'pages/anonymous/login/Login';
import { NotFound } from 'pages/anonymous/not-found/NotFound';
import Authorized from 'pages/authorized';
import { PROJECTS_URL, privateRouteMap, routes } from 'pages/authorized/authorizedRouteNames';
import { FC } from 'react';
import { BrowserRouter, Redirect, Route, Switch } from 'react-router-dom';
import PrivateRoute from 'router/PrivateRoute';
import PublicRoute from 'router/PublicRoute';

const Routes: FC = () => {
	return (
		<BrowserRouter>
			<Switch>
				<Redirect exact path="/" to={PROJECTS_URL} />
				<Route exact path={[LOGIN_URL]}>
					<Anonymous>
						<Switch>
							<PublicRoute restricted={true} component={Login} exact path={LOGIN_URL} />
						</Switch>
					</Anonymous>
				</Route>
				<Route exact path={routes.map(k => k.value)}>
					<Authorized>
						<Switch>
							{routes.map(route => (
								<PrivateRoute
									key={route.key}
									component={Reflect.get(privateRouteMap, route.key)}
									path={route.value}
									exact
								/>
							))}
						</Switch>
					</Authorized>
				</Route>
				<Route component={NotFound} />
			</Switch>
		</BrowserRouter>
	);
};

export default Routes;

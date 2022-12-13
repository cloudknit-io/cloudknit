import AuthStore from 'auth/AuthStore';
import { cleanDagNodeCache } from 'components/organisms/tree-view/node-figure-helper';
import { LOGIN_URL } from 'pages/anonymous/anonymousRouteNames';
import { NotFound } from 'pages/anonymous/not-found/NotFound';
import { QUICK_START_URL } from 'pages/authorized/authorizedRouteNames';
import { QuickStart } from 'pages/authorized/quick-start/QuickStart';
import { TermsAndConditions } from 'pages/authorized/terms-and-conditions/TermsAndConditons';
import React, { ElementType, FC, ReactNode, useEffect } from 'react';
import { Redirect, Route, RouteComponentProps, RouteProps } from 'react-router-dom';

interface PrivateRouteProps extends Omit<RouteProps, 'component'> {
	component: ElementType;
}

const PrivateRoute: FC<PrivateRouteProps> = ({ component: Component, ...rest }: PrivateRouteProps) => {
	useEffect(() => {
		if (rest.location?.pathname) {
			cleanDagNodeCache('temp');
		}
	}, [rest.location?.pathname]);

	return (
		// Show the component only when the user is logged in
		// Otherwise, redirect the user to /login page
		<Route
			{...rest}
			render={(props: RouteComponentProps<ReactNode>): ReactNode => {
				const user = AuthStore.getUser();
				if (user) {
					if (user.role !== 'Admin' && rest.location?.pathname?.includes('settings')) {
						return <NotFound />
					}

					if ((user.selectedOrg && !user.selectedOrg.githubRepo) || rest.location?.pathname === QUICK_START_URL) {
						return <QuickStart/>;
					}

					if (user.organizations.length === 0 || user.selectedOrg?.provisioned !== true) {
						return <TermsAndConditions />;
					}
					
					return <Component {...props} />;
				}

				return <Redirect to={LOGIN_URL} />;
			}}
		/>
	);
};

export default PrivateRoute;

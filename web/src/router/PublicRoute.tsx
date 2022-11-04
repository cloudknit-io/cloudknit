import AuthStore from 'auth/AuthStore';
import React, { ElementType, FC, ReactNode } from 'react';
import { Redirect, Route, RouteComponentProps, RouteProps } from 'react-router-dom';

interface PublicRouteProps extends Omit<RouteProps, 'component'> {
	component: ElementType;
	restricted: boolean;
}

const PublicRoute: FC<PublicRouteProps> = ({ component: Component, restricted, ...rest }: PublicRouteProps) => {
	return (
		// restricted = false meaning public route
		// restricted = true meaning restricted route
		<Route
			{...rest}
			render={(props: RouteComponentProps<ReactNode>): ReactNode =>
				AuthStore.getUser() && restricted ? <Redirect to="/dashboard" /> : <Component {...props} />
			}
		/>
	);
};

export default PublicRoute;

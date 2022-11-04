import './nav-bar.scss';

import classNames from 'classnames';
import PropTypes from 'prop-types';
import React, { FC } from 'react';

export interface NavBarProps extends React.Props<any> {
	items: Array<{ path: string; title: string }>;
	history: any;
	version?: () => React.ReactElement;
}

export function isActiveRoute(locationPath: string, path: string) {
	return locationPath === path || locationPath.startsWith(`${path}/`);
}

export const NavBar: FC<NavBarProps> = ({ items = [], history }: NavBarProps) => {
	return (
		<div className="nav-bar">
			{items.map(item => (
				<div
					className={classNames('nav-item', { active: isActiveRoute(history.location.pathname, item.path) })}
					onClick={() => history.push(item.path)}>
					{item.title}
				</div>
			))}
		</div>
	);
};

NavBar.contextTypes = {
	router: PropTypes.object,
};

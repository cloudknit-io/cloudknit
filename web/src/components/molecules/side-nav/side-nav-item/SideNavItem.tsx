import { ReactComponent as ArrowUp } from 'assets/images/icons/arrow_drop_up.svg';
import classNames from 'classnames';
import { NavItem } from 'models/nav-item.models';
import React, { FC, useState } from 'react';
import { useParams } from 'react-router';

import { SideNav } from '../SideNav';

export interface SideNavItemProps extends React.Props<any> {
	item: NavItem;
	history: any;
	version?: () => React.ReactElement;
	collapseNav: () => void;
}

export function isActiveRoute(locationPath: string, path: string) {
	return path ? locationPath === path || locationPath.startsWith(`${path}/`) : false;
}

// These need to be added and imported
// const icon: { [key: string]: JSX.Element } = {
//     Environments: <Settings className="nav-icon" />,
//     // Applications: <Settings className="nav-icon" />,
//     Reports: <Reports className="nav-icon" />,
//     Settings: <Settings className="nav-icon" />,
// };

export const SideNavItem: FC<SideNavItemProps> = ({ item, history, collapseNav }: SideNavItemProps) => {
	const getSubNav = (children: NavItem[] = []) => (
		<SideNav items={children} history={history} isSubMenu={true} collapseNav={collapseNav} />
	);
	const { projectId, environmentId } = useParams() as any;
	const pushRoute = () => {
		collapseNav();
		if (projectId && environmentId && item.optionalData) {
			history.push(`/${projectId}/${environmentId}/${item.optionalData.type}`, item.optionalData || {});
			return;
		}
		history.push(item.path, item.optionalData || {});
	};

	const getSecondaryPath = (item: NavItem) => {
		if (projectId && environmentId && item.optionalData) {
			return `/${projectId}/${environmentId}/${item.optionalData.type}`;
		}
		return '';
	};

	const isChildActiveRoute = (items: NavItem[] = []) => {
		return items.some(item =>
			isActiveRoute(history.location.pathname, item.optionalData ? getSecondaryPath(item) : item.path)
		);
	};

	const isActive = isActiveRoute(history.location.pathname, item.path) || isChildActiveRoute(item.children);

	const [getDropDownToggleValue, toggleDropDown] = useState(!isActive);

	return (
		<>
			<a
				className={classNames('side-nav-item', {
					active:
						isActiveRoute(
							history.location.pathname,
							item.optionalData ? getSecondaryPath(item) : item.path
						) || isChildActiveRoute(item.children),
				})}
				onClick={() => pushRoute()}>
				<span>
					{/* To Add Icons {icon[item.title]} */}
					{item.title}
				</span>
				{item.children ? (
					<ArrowUp
						className={`nav-drop-down ${getDropDownToggleValue ? 'collapsed' : ''}`}
						onClick={e => {
							e.stopPropagation();
							toggleDropDown(!getDropDownToggleValue);
						}}
					/>
				) : null}
			</a>
			{item.children ? (
				<div className={getDropDownToggleValue ? 'collapsed' : ''}>{getSubNav(item.children)}</div>
			) : null}
		</>
	);
};

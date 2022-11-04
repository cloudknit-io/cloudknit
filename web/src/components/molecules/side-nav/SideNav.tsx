import './style.scss';

import { NavItem } from 'models/nav-item.models';
import React, { FC } from 'react';

import { SideNavItem } from './side-nav-item/SideNavItem';

export interface SideNavProps extends React.Props<any> {
	items: NavItem[];
	history: any;
	version?: () => React.ReactElement;
	isSubMenu: boolean;
	collapseNav: () => void;
}

export const SideNav: FC<SideNavProps> = ({ items = [], history, isSubMenu, collapseNav }: SideNavProps) => {
	const getNavTree = (items: any[]): JSX.Element[] => {
		return items.map((item, _i) => (
			<SideNavItem item={item} history={history} key={`side-nav-item-${_i}`} collapseNav={collapseNav} />
		));
	};

	return <nav className={`side-nav-item-container ${isSubMenu ? 'sub-menu' : ''}`}>{getNavTree(items)}</nav>;
};

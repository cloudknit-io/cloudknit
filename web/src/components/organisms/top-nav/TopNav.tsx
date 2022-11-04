import './styles.scss';

import { ReactComponent as Chevron } from 'assets/images/icons/chevron-right.svg';
import { NavItem } from 'models/nav-item.models';
import React, { FC, useState } from 'react';
import { useHistory, useParams } from 'react-router';

type TopNavProps = {
	items: NavItem[];
	className?: string;
};

type TopNavItemProps = {
	item: NavItem;
	isActive: boolean;
};

export const TopNav: FC<TopNavProps> = ({ items, className }: TopNavProps) => {
	const history = useHistory();
	const isActiveRoute = (item: NavItem) => {
		if (history.location.pathname === item.path) {
			return true;
		}
		if (item.children?.some(child => child.path === history.location.pathname)) {
			return true;
		}
		return false;
	};
	return (
		<ul className={`Ztop-nav__menu ${className}`}>
			{items
				.filter(item => item.visible?.call(null) !== false)
				.map(item => (
					<TopNavItem key={item.title} item={item} isActive={isActiveRoute(item)} />
				))}
		</ul>
	);
};

const TopNavItem: FC<TopNavItemProps> = ({ item, isActive }: TopNavItemProps) => {
	const history = useHistory();
	const children = item.children?.filter(child => child.visible?.call(null) !== false) || [];
	return (
		<li className={`Ztop-nav__menu-item Ztop-nav__menu-item--${isActive ? 'active' : ''}`}>
			<a
				onClick={() => {
					history.push(item.path);
				}}>
				<span>{item.title}</span>
				{children.length > 0 ? (
					<span className="Ztop-nav__menu-item--icon-container">
						<Chevron className="arrow" />
					</span>
				) : null}
			</a>
			{children.length > 0 && <TopNav items={children} className="sub-nav" />}
		</li>
	);
};

import { ReactComponent as Logo } from 'assets/images/icons/logo.svg';
import { ZText } from 'components/atoms/text/Text';
import React, { useState } from 'react';
import AuthStore from 'auth/AuthStore';
import { ZDropdownMenuJSX } from 'components/molecules/dropdown-menu/DropdownMenu';
import { NavItem } from 'models/nav-item.models';
import { TopNav } from 'components/organisms/top-nav/TopNav';
import { BradAdarshFeatureVisible, FeatureKeys, FeatureRoutes } from 'pages/authorized/feature_toggle';

require('./top-bar.scss');

export interface TopBarFilter<T> {
	items: Array<{
		label?: string;
		value?: T;
		content?: (changeSelection: (selectedValues: T[]) => any) => React.ReactNode;
	}>;
	selectedValues: T[];
	selectionChanged: (selectedValues: T[]) => any;
}

export interface ActionMenu {
	className?: string;
	items: {
		action: () => any;
		title: string | React.ReactElement;
		iconClassName?: string;
		qeId?: string;
		disabled?: boolean;
	}[];
}

export interface Toolbar {
	filter?: TopBarFilter<any>;
	breadcrumbs?: { title: string; path?: string }[];
	tools?: React.ReactNode;
	actionMenu?: ActionMenu;
}

export interface TopBarProps {
	title: string;
	toolbar?: Toolbar;
}

const navItems: NavItem[] = [
	{ title: 'Teams', path: '/Teams' },
	{
		title: 'Environments',
		path: '/dashboard',
		children: [
			{
				title: 'Applications',
				path: '/all/all/apps',
				visible: () => Reflect.get(FeatureRoutes, FeatureKeys.APPLICATIONS),
			},
		],
	},
	{ title: 'Infra Components', path: '/all/all/infra' },
	{ title: 'Overview', path: '/overview'},
	{ title: 'Dashboard', path: '/demo-dashboard', visible: () => BradAdarshFeatureVisible() },
	{ title: 'Builder', path: '/builder', visible: () => BradAdarshFeatureVisible() },
	{ title: 'Settings', path: '/settings', visible: () => AuthStore.getUser()?.role === 'Admin' },
	{ title: 'Quick Start', path: '/quick-start'},
];

export const TopBar = ({ title }: TopBarProps) => {
	const [showDropdown, setShowDropDown] = useState<boolean>(false);
	const currentUser = AuthStore.getUser();

	return (
		<div className="top-bar" key="top-bar">
			<div className="top-bar__flex">
				<div className="top-bar__logo-container">
					<Logo style={{ width: '80px', marginRight: '30px' }} className="top-bar__logo" />
				</div>
			</div>
			<div className="top-bar__flex">
				<div className="top-bar__Ztop-nav">
					<nav>
						<TopNav items={navItems} />
					</nav>
				</div>
				<div className="top-bar__avatar">
					<img
						src={currentUser?.picture}
						height="36"
						width="36"
						onClick={() => setShowDropDown(!showDropdown)}
					/>
					<ZDropdownMenuJSX
						className="top-bar__avatar__dropdown"
						isOpened={showDropdown}
						items={[
							...(AuthStore.getUser()?.organizations || []).map(org => ({
								text: org.name || '',
								action: async () => {
									await AuthStore.selectOrganization(org.name);
								},
								selected: AuthStore.getUser()?.selectedOrg.name === org.name
							})),
							{ text: '', jsx: <a href={AuthStore.logoutUrl()}>Log&nbsp;Out</a>, action: () => true },
						]}
					/>
				</div>
			</div>
		</div>
	);
};

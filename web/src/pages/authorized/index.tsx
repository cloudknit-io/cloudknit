import { ReactComponent as ChevronRight } from 'assets/images/icons/chevron-right.svg';
import { Layout, TopBar } from 'components/argo-core';
import { SideNav } from 'components/molecules/side-nav/SideNav';
import { NavItem } from 'models/nav-item.models';
import React, { ReactElement, useState } from 'react';
import { useEffect } from 'react';
import { useHistory } from 'react-router-dom';
import { CostingService } from 'services/costing/costing.service';

import {
	breadcrumbObservable,
	EnvironmentPageHeaderCtx,
	pageHeaderObservable,
} from './contexts/EnvironmentHeaderContext';
import { EnvironmentHeader } from './environments/EnvironmentHeader';

const Authorized: React.FC = ({ children }) => {
	const history = useHistory();
	const [sideNavCollapsed, setSideNavCollapseState] = useState(true);
	const sideNavToggleCollapse = () => setSideNavCollapseState(!sideNavCollapsed);
	useEffect(() => {
		CostingService.getInstance().streamNotification();
	}, []);

	return (
		<Layout version={(): ReactElement => <div>v0.1.0</div>}>
			<div className="zlifecycle-header">
				<TopBar title="zLifecycle" />
			</div>
			<main className="zlifecycle-main-content">
				<div
					className={`dark-overlay ${sideNavCollapsed ? 'collapsed' : ''}`}
					onClick={() => sideNavToggleCollapse()}
				/>
				{/* <div className={`zlifecycle-nav family-dm ${sideNavCollapsed ? 'collapsed' : ''}`}>
					<button className="toggle-side-bar" onClick={() => sideNavToggleCollapse()}>
						<ChevronRight className="chevron-right" />
					</button>
					<SideNav items={navItems} history={history} isSubMenu={false} collapseNav={sideNavToggleCollapse} />
				</div> */}
				<div className="page zlifecycle-main">
					<EnvironmentPageHeaderCtx.Provider
						value={{
							breadcrumbObservable,
							pageHeaderObservable,
						}}>
						<EnvironmentHeader />
						{children}
					</EnvironmentPageHeaderCtx.Provider>
				</div>
			</main>
		</Layout>
	);
};

export default Authorized;

import { Layout, TopBar } from 'components/argo-core';
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

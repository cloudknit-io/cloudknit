import AuthStore from 'auth/AuthStore';
import { Layout, TopBar } from 'components/argo-core';
import { EntityStore } from 'models/entity.store';
import React, { ReactElement, useEffect } from 'react';

import {
	breadcrumbObservable,
	EnvironmentPageHeaderCtx,
	pageHeaderObservable
} from './contexts/EnvironmentHeaderContext';
import { EnvironmentHeader } from './environments/EnvironmentHeader';

const Authorized: React.FC = ({ children }) => {
	useEffect(() => {
		if (!AuthStore.getOrganization()) {
			return;
		}
		EntityStore.getInstance();
	}, []);

	return (
		<Layout version={(): ReactElement => <div>v0.1.0</div>}>
			<div className="zlifecycle-header">
				<TopBar title="zLifecycle" />
			</div>
			<main className="zlifecycle-main-content">
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

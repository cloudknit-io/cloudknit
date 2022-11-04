import AuthStore from 'auth/AuthStore';
import React from 'react';
import { FC } from 'react';

export const OnBoarding: FC = () => {
	const user = AuthStore.getUser();

	if (!user) {
		AuthStore.logout();
	}

	return (
		<section className="d-flex align-center justify-center" style={{ height: '100vh' }}>
			<em style={{textAlign: 'center', fontSize: '1.2em'}}>
				The infrastructure for{' '}
				<b style={{ color: 'teal' }}>{user?.selectedOrg.name}</b> is currently
				being provisioned.<br/> You will receive an email once it's completed. Until then, please read through{' '}
				<a style={{color: 'teal'}} href="https://docs.zlifecycle.com/getting_started/setup/" target="_blank">
					these docs
				</a>{' '}
				to continue your onboarding.
			</em>
		</section>
	);
};

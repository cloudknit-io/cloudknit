import AuthStore from 'auth/AuthStore';
import { ZLoaderCover } from 'components/atoms/loader/LoaderCover';
import { AuthContext } from 'context/auth-provider/AuthContext';
import { AuthState } from 'models/auth.models';
import { LocalStorageKey } from 'models/localStorage';
import { Profile, User } from 'models/user.models';
import { OrganizationSelection } from 'pages/authorized/dashboard/OrganizationSelection';
import { QuickStart, QuickStartContext } from 'pages/authorized/quick-start/QuickStart';
import { TermsAndConditions } from 'pages/authorized/terms-and-conditions/TermsAndConditons';
import React, { Children, useContext, useEffect, useState } from 'react';
import { AuthService } from 'services/auth/AuthService';
import ApiClient from 'utils/apiClient';
import { LocalStorage } from 'utils/localStorage/localStorage';

export const AuthProvider: React.FC = ({ children }) => {
	const [loading, setLoading] = useState<boolean>(false);
	const [authState, setAuthState] = useState<User>();

	useEffect(() => {
		ApiClient.setAuthState = (authStateModified: User | undefined): void => setAuthState(authStateModified);
		ApiClient.init();
		setLoading(true);
		AuthService.me()
			.then(async response => {;
				const user: User = response.data;
				let selectedOrg = null;
				if (user.organizations.length === 1) {
					selectedOrg = user.organizations[0];
				} else {
					selectedOrg =  AuthStore.getOrganization();
				}
				if (selectedOrg) {
					await AuthStore.selectOrganization(selectedOrg.name, false);	
					user.selectedOrg = selectedOrg;
				}
				setAuthState(user);
				setLoading(false);
			})
			.catch(() => {
				LocalStorage.removeItem(LocalStorageKey.USER);
				setLoading(false);
			});
	}, []);

	const value = {
		user: authState,
		setAuthState,
	} as AuthState;

	const getLandingScreen = () => {
		if (loading) {
			return <ZLoaderCover loading />;
		}

		if (authState) {
			if (Number(authState?.organizations?.length) > 1 && !authState.selectedOrg) {
				return <OrganizationSelection />
			}
		}

		return Children.only(children);
	}

	return (
		<div style={{ width: '100vw', height: '100vh' }}>
			<AuthContext.Provider value={value}>
				{getLandingScreen()}
			</AuthContext.Provider>
		</div>
	);
};
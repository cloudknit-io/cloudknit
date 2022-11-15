import './login.scss';

import github from 'assets/images/github.png';
import { ReactComponent as Logo } from 'assets/images/ck-logo.svg';
import AuthStore from 'auth/AuthStore';
import { Button } from 'components/atoms/button/Button';
import React, { useEffect, useState } from 'react';

export const Login: React.FC = () => {
	const [loading, setLoading] = useState<ActiveLoginType>(ActiveLoginType.NONE);

	const [error, setError] = useState<any>(null);
	const [acctCreateUrl, setAcctCreateUrl] = useState<any>(null);

	const login = (): void => {
		setLoading(ActiveLoginType.GIT);
		AuthStore.login();
	};

	useEffect(() => {
		const urlParams = new URLSearchParams(window.location.search);
		const error_keys = ['sso_error', 'auth_error'];
		
		const has_error = error_keys.map(ek => {
			return urlParams.has(ek) ? urlParams.get(ek) : null;
		}).find(ek => ek);

		if (!has_error) {
			return;
		}

		const createKey = "account_create_url";
		const url = urlParams.has(createKey) ? urlParams.get(createKey) : null;
		
		setAcctCreateUrl(url);
		setError(has_error);
	}, [window.location.search]);

	return (
		<div className="login-layout">
			<div className="login-panel">
				<div className="login-panel--content">
					<Logo style={{ maxWidth: '400px', width: '100%' }} className="top-bar__logo" />
				</div>
			</div>
			<div className="login-form-container">
				<div className="login-form">
					{error && <div style={{ color: 'red' }}>{error}</div>}
					{acctCreateUrl && <a href={acctCreateUrl} style={{ color: 'green' }}>Create Account</a>}
					<Button
						onClick={login}
						className="login-form--button login-form--button-github"
						isLoading={loading === ActiveLoginType.GIT}>
						<img className="login-form--button-icon" src={github} alt="" width="24px" />
						Log in with GitHub
					</Button>
				</div>
			</div>
		</div>
	);
};

enum ActiveLoginType {
	NONE,
	GIT,
	CUSTOM,
}

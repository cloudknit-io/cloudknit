import AuthStore from 'auth/AuthStore';
import React, { useContext, useEffect, useMemo, useState } from 'react';
import ApiClient from 'utils/apiClient';

import { ReactComponent as Copy } from 'assets/images/icons/copy.svg';
import { NotificationType } from 'components/argo-core';
import { Loader } from 'components/atoms/loader/Loader';
import { Context } from 'context/argo/ArgoUi';

type AccessToken = {
	access_token: string;
	expires_in: number;
	token_type: string;
};

export const AccessToken: React.FC = () => {
	const org = useMemo(() => AuthStore.getOrganization(), []);
	const nm = useContext(Context)?.notifications;
	const [accessToken, setAccessToken] = useState<string | null | undefined>(undefined);
	useEffect(() => {
		if (!org) return;
		ApiClient.get<AccessToken>('/auth/access-token')
			.then(resp => {
				setAccessToken(resp.data.access_token);
			})
			.catch(e => {
				nm?.show({
					content: 'Failed to retrieve access-token',
					type: NotificationType.Error,
					toastOptions: {
						autoClose: 3000,
					},
				});
				console.error(e);
				setAccessToken(null);
			});
	}, [org]);

	const renderUI = () => {
		if (!org) return <></>;
		if (accessToken === undefined)
			return (
				<>
					<Loader height={16} />
					Retrieving access token...
				</>
			);
		if (accessToken === null) return 'Failed to retrieve access token';

		return (
			<>
				<div className="secret-pair">
					<label className="secret-pair__name">Access Token</label>
					<span className="secret-pair__dummy-value">******</span>
				</div>
				<button
					className="copy-to-clipboard"
					onClick={e => {
						navigator.clipboard.writeText(accessToken);
						nm?.show({
							content: 'Copied successfully',
							type: NotificationType.Success,
							toastOptions: {
								autoClose: 1000,
							},
						});
					}}>
					<Copy />
				</button>
			</>
		);
	};

	return (
		<div style={{ height: '100px' }}>
			<div className="secret-info secret-info-ssm">{renderUI()}</div>
		</div>
	);
};

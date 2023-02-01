import AuthStore from 'auth/AuthStore';
import React, { useCallback, useContext, useEffect, useMemo, useState } from 'react';
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
    const [wait, setWait] = useState<boolean>(false);

	const generateToken = useCallback(() => {
        setWait(true);
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
			}).finally(() => {
                setWait(false);
            });
	}, []);

	const renderUI = () => {
		if (!org) return <></>;
		if (!accessToken)
			return (
				<div className="secret-pair">
					<label className="secret-pair__name">Generate an access token for personal access.</label>
					<button disabled={wait} className="secret-pair__button secret-pair__save shadowy-input" onClick={generateToken}>
						{wait ? <Loader height={16} /> : 'Generate Token'}
					</button>
				</div>
			);

		return (
			<>
				<div className="secret-pair">
					<label className="secret-pair__name secret-pair__name__warning">
						Please copy this token, right away! You will not be able to see it again
					</label>
					<span className="secret-pair__dummy-value secret-pair__dummy-value--with-border">{accessToken}</span>
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

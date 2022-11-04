import React, { useEffect, useState } from 'react';
import { SecretsService } from 'services/secrets/secrets.service';
import { AWSSSMSecret } from './aws-ssm-secrets';
import { Loader } from 'components/atoms/loader/Loader';
import { ReactComponent as Add } from 'assets/images/icons/add.svg';
import { getTime } from '../environment-components/helpers';

export interface Props {
	heading: string;
	secretKey: string;
	newSecret?: boolean;
	refreshEnvList?: (envName: string) => void;
}

type SecretResponseType = {
	key: string;
	value: string;
	lastModifiedDate: Date;
};

export const SecretList: React.FC<Props> = ({ heading, secretKey, newSecret, refreshEnvList }: Props) => {
	const [selectedSSMSecret, setSSMSelectedSecret] = useState<any>(null);
	const [selectedSecretKey, setSelectedSecretKey] = useState<any>(null);
	const secretsService = SecretsService.getInstance();
	const [existingSecrets, setExistingSecrets] = useState<Map<string, SecretResponseType[]>>(
		new Map<string, SecretResponseType[]>()
	);
	const [updateSecretList, setUpdateSecretList] = useState<boolean>(true);
	const [isLoading, setLoading] = useState<boolean>(false);

	useEffect(() => {
		if (!updateSecretList || !secretKey) {
			return;
		}
		if (newSecret) {
			return;
		}
		updateSecretsList();
		setLoading(true);
	}, [secretsService, updateSecretList, secretKey]);

	const updateSecretsList = () => {
		let path;
		
		if (secretKey.indexOf(':') > 0) {
			// secretKey = teamName:envName
			// we only care about envName
			path = secretKey.replaceAll(':', '/');
		}

		secretsService.getSsmSecrets(false, path).then(({ data }) => {			
			const resp = data as any;
			existingSecrets.clear();
			
			if (resp?.length > 0) {
				(resp as any).forEach((d: SecretResponseType) => {
					if (!existingSecrets.has(d.key)) {
						existingSecrets.set(d.key, []);
					}
					existingSecrets.get(d.key)?.push(d);
				});
			}

			setExistingSecrets(new Map([...existingSecrets.entries()]));
			setLoading(false);
		});
	};

	const saveCallback = (key: string) => {
		updateSecretsList();
		setSelectedSecretKey(key);
	};

	const closeCallback = () => {
		setSSMSelectedSecret(null);
	};

	return (
		<>
			{isLoading ? (
				<span className="d-flex align-center" style={{ paddingTop: '40px' }}>
					<Loader height={20} width={20} /> Checking for Existing Secrets
				</span>
			) : (
				<>
					{newSecret && (
						<span
							key={secretKey}
							onClick={e => {
								setSSMSelectedSecret(secretKey.replaceAll(':', '/'));
								setSelectedSecretKey(null);
							}}>
							+
						</span>
					)}

					<ul className="secrets-list secrets-list__1">
						<li
							key={`${secretKey}-new`}
							className="secrets-list__item"
							onClick={() => {
								setSSMSelectedSecret(secretKey.replaceAll(':', '/'));
								setSelectedSecretKey(null);
							}}>
							<span className="d-flex align-center">
								New <Add style={{ marginLeft: '5px' }} />
							</span>
						</li>
						{(existingSecrets.get(secretKey) || []).map(e => (
							<li
								key={`${secretKey}-${e.value}`}
								className="secrets-list__item d-flex"
								style={{ flexDirection: 'column' }}
								onClick={() => {
									setSSMSelectedSecret(secretKey.replaceAll(':', '/'));
									setSelectedSecretKey(e.value);
								}}>
								{e.value}{' '}
								<small style={{ marginLeft: '5px' }}>
									Last Updated, ({getTime(e.lastModifiedDate.toString() || '')})
								</small>
							</li>
						))}
					</ul>
					<div
						className={`secret-info secret-info-ssm secret-info-ssm-abs ${
							selectedSSMSecret === null ? 'hidden' : ''
						}`}
						onClick={e => {
							if (e.currentTarget === e.target) {
								closeCallback();
							}
						}}>
						<div className="secret-container">
							{
								<AWSSSMSecret
									secretScope={selectedSSMSecret}
									secretKey={selectedSecretKey}
									saveCallback={saveCallback}
									closeCallback={closeCallback}
									scopeEditable={newSecret}
								/>
							}
						</div>
					</div>
				</>
			)}
		</>
	);
};

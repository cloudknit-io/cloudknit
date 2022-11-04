import { NotificationsApi, NotificationType } from 'components/argo-core';
import { Loader } from 'components/atoms/loader/Loader';
import { Context } from 'context/argo/ArgoUi';
import React, { FormEvent } from 'react';
import { useEffect } from 'react';
import { useState } from 'react';
import { SecretsService } from 'services/secrets/secrets.service';
import { ReactComponent as Copy } from 'assets/images/icons/copy.svg';
import { getTime } from '../environment-components/helpers';

interface Props {
	secretScope: string;
	closeCallback: (id?: any) => void;
	newEnvironmentField?: boolean;
	existCallback?: (exists: boolean) => void
}

type AWSSecretType = {
	key: string;
	exists: boolean;
	lastModifiedDate?: Date;
};

export const AWSSecret: React.FC<Props> = ({ secretScope, newEnvironmentField, closeCallback, existCallback }: Props) => {
	const [secretMap, setSecretMap] = useState<Map<string, AWSSecretType>>(new Map<string, AWSSecretType>());
	const [showForm, setShowForm] = useState<boolean>(false);
	const [updating, setUpdating] = useState<boolean>(false);
	const [existing, setExisting] = useState<boolean | null>(null);
	const [localScope, setLocalScope] = useState<string>(secretScope);
	const nm = React.useContext(Context)?.notifications as NotificationsApi;
	const secretService = SecretsService.getInstance();
	const getFormPair = ({ awsKey, label }: any, required: boolean) => {
		return (
			<>
				<label className="secret-pair__name">{label}</label>
				<input required={required} name={awsKey} className="secret-pair__input shadowy-input" type="text" />
			</>
		);
	};

	const getFormTextAreaPair = ({ awsKey, label }: any) => {
		return (
			<>
				<label className="secret-pair__name">{label}</label>
				<textarea
					required={existing === false}
					name={awsKey}
					className="secret-pair__input shadowy-input"
					rows={3}
					style={{ outline: 'none', resize: 'vertical' }}></textarea>
			</>
		);
	};

	const secretKeys = [
		{
			awsKey: 'aws_access_key_id',
			label: 'AWS Access Key Id',
			getFormPair: function () {
				return getFormPair(this, existing === false);
			},
			getDummyPair: function () {
				return getDummyPair(this);
			},
		},
		{
			awsKey: 'aws_secret_access_key',
			label: 'AWS Secret Access Key',
			getFormPair: function () {
				return getFormTextAreaPair(this);
			},
			getDummyPair: function () {
				return getDummyPair(this);
			},
		},
		{
			awsKey: 'aws_session_token',
			label: 'AWS Session Token (optional)',
			getFormPair: function () {
				return getFormPair(this, false);
			},
			getDummyPair: function () {
				return getDummyPair(this);
			},
		},
	];

	useEffect(() => {
		setExisting(null);
		if (!secretScope) {
			return;
		}

		setLocalScope(secretScope);
		if (secretMap.get(secretScope)) {
			return;
		}

		if (newEnvironmentField) {
			setExisting(false);
			return;
		}

		secretService
			.secretsExists(
				secretKeys.map(e => e.awsKey),
				secretScope
			)
			.then(({ data }) => {
				if (data === false) {
					setExisting(false);
					setSecretMap(new Map());
				} else if (Array.isArray(data)) {
					const exists = data.some(e => e.exists);
					setExisting(exists);
					existCallback?.call(null, exists);
					setSecretMap(new Map(data.map((e: AWSSecretType) => [e.key, e])));
				}
			})
			.catch(err => {
				setExisting(false);
				nm.show({
					type: NotificationType.Error,
					content: 'Could not verify existing secrets',
				});
			});
	}, [secretScope]);

	const updateSecret = (e: FormEvent) => {
		e.preventDefault();
		const formData = new FormData(e.target as HTMLFormElement);
		const awsAccesskeyId = formData.get('aws_access_key_id')?.toString().trim();
		const awsSecretAccessKey = formData.get('aws_secret_access_key')?.toString().trim();
		const awsSessionToken = formData.get('aws_session_token')?.toString().trim();
		const secrets: any = [];
		if (awsAccesskeyId) {
			secrets.push({
				key: 'aws_access_key_id',
				value: awsAccesskeyId,
			});
		}
		if (awsSecretAccessKey) {
			secrets.push({
				key: 'aws_secret_access_key',
				value: awsSecretAccessKey,
			});
		}
		if (awsSessionToken) {
			secrets.push({
				key: 'aws_session_token',
				value: awsSessionToken,
			});
		}
		if (secrets.length > 0) {
			setUpdating(true);
			secretService
				.updateAWSSecret(secrets, localScope)
				.then(({ data }) => {
					if (data === false) {
						throw 'Not Updated';
					}
					if (data === true) {
						closeCallback(localScope);
						setUpdating(false);
						setShowForm(false);
						setExisting(true);
						secrets.forEach((s: any) => {
							secretMap.set(s.key, {
								exists: true,
								key: s.key,
								lastModifiedDate: new Date(),
							});
						});
						setSecretMap(new Map([...secretMap.entries()]));
						nm.show({
							type: NotificationType.Success,
							content: `Creds updated successfully!`,
						});
					}
				})
				.catch(err => {
					setUpdating(false);
					nm.show({
						type: NotificationType.Error,
						content: 'There was an error while updating secrets!',
					});
				});
		} else {
			nm.show({
				type: NotificationType.Warning,
				content: 'No Data to update.',
			});
		}
	};

	const getInputFieldForEnvironmentName = () => {
		if (!localScope) {
			return <></>;
		}
		const tokens = localScope.split('/');
		return (
			<span>
				{tokens[0]}
				{'/'}
				{tokens[1]}
				{'/'}
				<input
					required
					placeholder={'environment...'}
					name="secret-scope-environment"
					className="scope-edit shadowy-input"
					type="text"
					defaultValue={''}
					onChange={e => {
						setLocalScope(`${tokens[0]}/${tokens[1]}/${e.target.value}`);
					}}
				/>
			</span>
		);
	};

	const getDummyPair = ({ label, awsKey }: any) => {
		return (
			<>
				<label className="secret-pair__name">
					{label}
					<button
						className="copy-to-clipboard"
						onClick={e => {
							copySecretPath(`/${secretScope}/${awsKey}`);
						}}>
						<Copy />
					</button>
				</label>
				<span className="secret-pair__dummy-value">
					{secretMap.get(awsKey)?.exists ? (
						<span className='d-flex justify-between'>
							<b>{'******'}</b>
							<small>
								Last updated, {getTime(secretMap.get(awsKey)?.lastModifiedDate?.toString() || '')}
							</small>
						</span>
					) : (
						'Does not exist.'
					)}
				</span>
			</>
		);
	};

	const copySecretPath = (text: string) => {
		navigator.clipboard.writeText(text);
		nm.show({
			content: 'path copied successfully',
			type: NotificationType.Success,
			toastOptions: {
				autoClose: 2000,
			},
		});
	};

	return (
		<>
			{existing === null ? (
				<span style={{ display: 'flex', alignItems: 'center' }}>
					<Loader height={20} width={20} /> Checking for Existing Secrets
				</span>
			) : (
				<div>
					{showForm || existing === false ? (
						<form className="secret-pair form" onSubmit={e => updateSecret(e)}>
							<h5 className="secrets-container__heading">
								{newEnvironmentField ? getInputFieldForEnvironmentName() : localScope}
							</h5>
							{secretKeys.map(e => e.getFormPair())}
							<div>
								<button
									type="submit"
									disabled={updating ? true : false}
									className="secret-pair__button secret-pair__save shadowy-input">
									{updating ? 'Updating...' : 'Save'}
								</button>
								{existing && (
									<button
										onClick={() => {
											setShowForm(false);
										}}
										type="button"
										disabled={updating ? true : false}
										className="secret-pair__button secret-pair__cancel shadowy-input">
										{'Back'}
									</button>
								)}
							</div>
						</form>
					) : (
						<>
							<h5 className="secrets-container__heading">
								<span>{localScope}</span>
							</h5>
							<p className="secret-pair">
								{secretKeys.map(e => e.getDummyPair())}

								<div>
									<button
										onClick={() => {
											setShowForm(true);
										}}
										className="secret-pair__button secret-pair__update shadowy-input">
										Update
									</button>
									{existing && (
										<button
											onClick={() => {
												if (!window.confirm('Are you sure you want to clear the secret?')) {
													return;
												}
												setUpdating(true);
												Promise.all([
													secretKeys
														.filter(e => secretMap.get(e.awsKey)?.exists)
														.map(e =>
															secretService.deleteSsmSecret(
																`/${secretScope.toLowerCase()}/${e.awsKey}`
															)
														),
												])
													.then(() => {
														setUpdating(false);
														secretKeys.forEach(e =>
															secretMap.set(e.awsKey, { key: e.awsKey, exists: false })
														);
														setSecretMap(new Map([...secretMap.entries()]));
														nm.show({
															content: 'Cleared Credentials.',
															type: NotificationType.Success,
														});
														setExisting(false);
													})
													.catch(() => {
														setUpdating(false);
														nm.show({
															content: 'There was an error while deleting.',
															type: NotificationType.Error,
														});
													});
											}}
											className="secret-pair__button secret-pair__cancel shadowy-input">
											Clear
										</button>
									)}
									{secretMap.get('aws_session_token')?.exists && (
										<button
											disabled={updating ? true : false}
											onClick={() => {
												setUpdating(true);
												secretService
													.deleteSsmSecret(`/${secretScope.toLowerCase()}/aws_session_token`)
													.then(() => {
														setUpdating(false);
														secretMap.set('aws_session_token', {
															key: 'aws_session_token',
															exists: false,
														});
														setSecretMap(new Map([...secretMap.entries()]));
														nm.show({
															content: 'Deleted Successfully',
															type: NotificationType.Success,
														});
													})
													.catch(() => {
														setUpdating(false);
														nm.show({
															content: 'There was an error while deleting.',
															type: NotificationType.Error,
														});
													});
											}}
											className="secret-pair__button secret-pair__cancel shadowy-input">
											{updating ? 'Deleting...' : 'Delete Session Token'}
										</button>
									)}
								</div>
							</p>
						</>
					)}
				</div>
			)}
		</>
	);
};

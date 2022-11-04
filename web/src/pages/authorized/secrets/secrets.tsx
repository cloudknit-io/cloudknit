import { NotificationsApi, NotificationType } from 'components/argo-core';
import { Loader } from 'components/atoms/loader/Loader';
import { Context } from 'context/argo/ArgoUi';
import React, { FormEvent, useMemo } from 'react';
import { useEffect } from 'react';
import { useState } from 'react';
import { SecretsService } from 'services/secrets/secrets.service';
import { ReactComponent as Copy } from 'assets/images/icons/copy.svg';
import { SecretModel } from './secret-model';
import { getTime } from '../environment-components/helpers';

interface Props {
	secretScope: string;
	closeCallback: (id?: any) => void;
	newEnvironmentField?: boolean;
	secretModels: SecretModel[];
}

type SecretResponseType = {
	key: string;
	exists: boolean;
	lastModifiedDate?: Date;
};

const getFormPair = ({ key, name, multiline, notRequired, hide }: SecretModel) => {
	return (
		<>
			<label className="secret-pair__name">{name}</label>
			{multiline ? (
				<textarea
					required={!notRequired}
					name={key}
					className="secret-pair__input shadowy-input"
					rows={3}
					style={{ outline: 'none', resize: 'vertical' }}></textarea>
			) : (
				<input required={!notRequired} name={key} className="secret-pair__input shadowy-input" type="text" />
			)}
		</>
	);
};

const getDummyPair = (
	{ name, key }: SecretModel,
	scope: string,
	secret: SecretResponseType | undefined,
	copyCallback: (key: string) => void
) => {
	return (
		<>
			<label className="secret-pair__name">
				{name}
				<button
					className="copy-to-clipboard"
					onClick={e => {
						copyCallback(`/${scope}/${key}`);
					}}>
					<Copy />
				</button>
			</label>
			<span className="secret-pair__dummy-value">
				{secret?.exists ? (
					<span className="d-flex justify-between">
						<b>{'******'}</b>
						<small>Last Updated, {getTime(secret.lastModifiedDate?.toString() || '')}</small>
					</span>
				) : (
					'Does not exist.'
				)}
			</span>
		</>
	);
};

export const Secrets: React.FC<Props> = ({ secretScope, closeCallback, secretModels, newEnvironmentField }: Props) => {
	const [secretMap, setSecretMap] = useState<Map<string, SecretResponseType>>(new Map<string, SecretResponseType>());
	const [showForm, setShowForm] = useState<boolean>(false);
	const [updating, setUpdating] = useState<boolean>(false);
	const [existing, setExisting] = useState<boolean | null>(null);
	const [localScope, setLocalScope] = useState<string>(secretScope);
	const nm = React.useContext(Context)?.notifications as NotificationsApi;
	const secretService = useMemo(() => SecretsService.getInstance(), []);

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

	const secretKeys = useMemo(
		() =>
			secretModels.map(e => {
				return {
					form:
						Boolean(secretMap.get(e.key)?.exists) && e.immutable
							? getDummyPair(e, secretScope, secretMap.get(e.key), copySecretPath)
							: getFormPair(e),
					dummy: getDummyPair(e, secretScope, secretMap.get(e.key), copySecretPath),
				};
			}),
		[secretModels, secretMap]
	);

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
				secretModels.map(e => e.key),
				secretScope
			)
			.then(({ data }) => {
				console.log(data);
				if (data === false) {
					setExisting(false);
					setSecretMap(new Map());
				} else if (Array.isArray(data)) {
					setExisting(data.some(e => e.exists));
					setSecretMap(new Map(data.map((e: SecretResponseType) => [e.key, e])));
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
		const secrets: any = secretModels
			.filter(e => formData.get(e.key)?.toString().trim())
			.map(e => ({
				key: e.key,
				value: formData.get(e.key)?.toString().trim(),
			}));

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
							secretMap.set(s.key, { key: s.key, exists: true, lastModifiedDate: new Date() });
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
							{secretKeys.map(s => s.form)}
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
								{secretKeys.map(s => s.dummy)}

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
													secretModels
														.filter(e => !e.immutable && secretMap.get(e.key)?.exists)
														.map(e =>
															secretService.deleteSsmSecret(
																`/${secretScope.toLowerCase()}/${e.key}`
															)
														),
												])
													.then(() => {
														setUpdating(false);
														secretModels.forEach(
															e =>
																!e.immutable &&
																secretMap.set(e.key, {
																	key: e.key,
																	exists: false,
																})
														);
														setSecretMap(new Map([...secretMap.entries()]));
														nm.show({
															content: 'Cleared Credentials.',
															type: NotificationType.Success,
														});
														setExisting([...secretMap.entries()].some(e => e[1].exists));
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
								</div>
							</p>
						</>
					)}
				</div>
			)}
		</>
	);
};

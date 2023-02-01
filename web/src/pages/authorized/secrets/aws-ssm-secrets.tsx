import { NotificationsApi, NotificationType } from 'components/argo-core';
import { Context } from 'context/argo/ArgoUi';
import React, { FormEvent, useEffect } from 'react';
import { useState } from 'react';
import { SecretsService } from 'services/secrets/secrets.service';
import { ReactComponent as Copy } from 'assets/images/icons/copy.svg';

interface Props {
	secretScope: string;
	secretKey: any;
	saveCallback: (key: string) => void;
	closeCallback: () => void;
	scopeEditable?: boolean;
}

export const AWSSSMSecret: React.FC<Props> = ({
	secretScope,
	secretKey,
	saveCallback,
	closeCallback,
	scopeEditable,
}: Props) => {
	const [updating, setUpdating] = useState<boolean>(false);
	const [showForm, setShowForm] = useState<boolean>(false);
	const nm = React.useContext(Context)?.notifications as NotificationsApi;
	const [scope, setScope] = useState<string>(secretScope);
	const secretService = SecretsService.getInstance();

	useEffect(() => {
		if (!secretScope) {
			return;
		}
		setScope(secretScope);
	}, [secretScope]);

	const updateSecret = (e: FormEvent) => {
		e.preventDefault();
		if (scopeEditable && !scopeValid()) {
			nm.show({
				content: 'secret path not valid!',
				type: NotificationType.Error,
			});
			return;
		}
		const formData = new FormData(e.target as HTMLFormElement);
		const secretId = secretKey ? secretKey : formData.get('secret-id')?.toString().trim();
		const secretValue = formData.get('secret-value')?.toString().trim();
		if (secretId && secretValue) {
			setUpdating(true);
			secretService
				.putSsmSecret(`${scope}/${secretId}`, secretValue)
				.then(({ data }) => {
					if (data === false) {
						throw 'Not Updated';
					}
					if (data === true) {
						setUpdating(false);
						setShowForm(false);
						nm.show({
							type: NotificationType.Success,
							content: 'Secrets updated successfully!',
						});
						saveCallback(scope);
						closeCallback();
					}
				})
				.catch(err => {
					setUpdating(false);
					nm.show({
						type: NotificationType.Error,
						content: 'There was an error while updating secrets!',
					});
				});
		}
	};

	const deleteSecret = () => {
		if (secretKey) {
			secretService.deleteSsmSecret(`${scope}/${secretKey}`).then(({ data }) => {
				setUpdating(false);
				if (data === false) {
					nm.show({
						type: NotificationType.Error,
						content: 'Failed to delete secret!',
					});
					return;
				}
				nm.show({
					type: NotificationType.Success,
					content: 'Secret Deleted successfully!',
				});
				saveCallback(scope);
				closeCallback();
			});
		}
	};

	const addSecretForm = () => {
		return (
			<div>
				{showForm || !secretKey ? (
					<form className="secret-pair form" onSubmit={e => updateSecret(e)}>
						<h5 className="secrets-container__heading">{scopeEditable ? getEditableScopeNode() : scope}</h5>
						<label className="secret-pair__name">Secret Id</label>
						{secretKey ? (
							<label>{secretKey}</label>
						) : (
							<input required name="secret-id" className="secret-pair__input shadowy-input" type="text" />
						)}
						<label className="secret-pair__name">Secret Value</label>
						<textarea
							required
							name="secret-value"
							className="secret-pair__input shadowy-input"
							rows={3}
							style={{ outline: 'none', resize: 'vertical' }}></textarea>
						<div>
							<button
								type="submit"
								disabled={updating ? true : false}
								className="secret-pair__button secret-pair__save shadowy-input">
								{updating ? 'Updating...' : 'Save'}
							</button>

							{!scopeEditable && <button
								type="button"
								onClick={e => {
									const btnElement = e.nativeEvent.currentTarget as HTMLButtonElement;
									btnElement.form?.reset();
									setShowForm(false);
									closeCallback();
								}}
								disabled={updating ? true : false}
								className="secret-pair__button secret-pair__cancel shadowy-input">
								Close
							</button>}
						</div>
					</form>
				) : (
					<>
						<h5 className="secrets-container__heading">
							<span>{scope}</span>
							<button
								className="copy-to-clipboard"
								onClick={e => {
									navigator.clipboard.writeText(`/${scope}/${secretKey}`);
									nm.show({
										content: 'path copied successfully',
										type: NotificationType.Success,
										toastOptions: {
											autoClose: 2000,
										},
									});
								}}>
								<Copy />
							</button>
						</h5>
						<div className="secret-pair">
							<label className="secret-pair__name">Secret Id</label>
							<span className="secret-pair__dummy-value">{secretKey}</span>
							<label className="secret-pair__name">Secret Value</label>
							<span className="secret-pair__dummy-value">******</span>
							<div>
								<button
									disabled={updating ? true : false}
									onClick={() => {
										setShowForm(true);
									}}
									className="secret-pair__button secret-pair__update shadowy-input">
									Update
								</button>
								<button
									disabled={updating ? true : false}
									onClick={() => {
										const resp = window.confirm(`Are you sure, you want to delete "${secretKey}"?`);
										if (!resp) {
											return;
										}
										setUpdating(true);
										deleteSecret();
									}}
									className="secret-pair__button secret-pair__cancel shadowy-input">
									{updating ? 'Deleting...' : 'Delete'}
								</button>
								{!scopeEditable && <button
									type="button"
									disabled={updating ? true : false}
									onClick={e => {
										const btnElement = e.nativeEvent.target as HTMLButtonElement;
										btnElement.form?.reset();
										closeCallback();
									}}
									className="secret-pair__button secret-pair__cancel shadowy-input">
									Close
								</button>}
							</div>
						</div>
					</>
				)}
			</div>
		);
	};

	const scopeValid = () => {
		const tokens = scope.split('/');
		return tokens.length === 3 && tokens.every(t => t.length > 0);
	};

	const getEditableScopeNode = () => {
		if (!scope) {
			return <></>;
		}
		const tokens = scope.split('/');
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
					defaultValue={tokens[2]}
					onChange={e => {
						setScope(`${tokens[0]}/${tokens[1]}/${e.target.value}`);
					}}
				/>
			</span>
		);
	};

	return <div className="secret-container">{addSecretForm()}</div>;
};

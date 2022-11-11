import { BaseGuide } from './BaseGuide';
import React, { useContext, useState } from 'react';
import { IGuide } from './IGuide';
import ApiClient from 'utils/apiClient';
import { Context } from 'context/argo/ArgoUi';
import { NotificationsManager, NotificationType } from 'components/argo-core';
import AuthStore from 'auth/AuthStore';
import { ZLoaderCover } from 'components/atoms/loader/LoaderCover';
import { Organization } from 'models/user.models';

type Props = {
	baseClassName: string;
	ctx: any;
	nm?: NotificationsManager;
};

export class ConfiguringZlifecycle extends BaseGuide implements IGuide {
	private static instance: ConfiguringZlifecycle | null = null;
	private saveGithubCredentials: () => Promise<any> = () => Promise.resolve(null);
	private context: any = AuthStore.getOrganization()?.githubRepo ? AuthStore.getOrganization() : null;

	private ConfiguringZlifecycleUI: React.FC<Props> = ({ baseClassName, ctx, nm }) => {
		const cls = (className: string) => `${baseClassName}_section-guide${className}`;
		const formRef = React.useRef<HTMLFormElement>(null);
		const [settingUp, SettingUpProgress] = useState<boolean>(false);
		const [repoUrl, setRepoUrl] = useState<string>(ctx?.githubRepo || AuthStore.getOrganization()?.githubRepo || '');

		this.saveGithubCredentials = async (): Promise<any> => {
			if (!formRef.current) return false;
			if (!formRef.current.checkValidity()) {
				formRef.current.classList.add('invalid');
				return false;
			}

			formRef.current.classList.remove('invalid');
			const formData = new FormData(formRef.current);
			SettingUpProgress(true);
			const payload = {
				githubRepo: formData.get('githubRepo'),
			};

			try {
				const res = await ApiClient.patch<Organization>(`/orgs/${AuthStore.getOrganization()?.id}`, payload);
				this.context = res.data.githubRepo ? res.data : null;
				setRepoUrl(res.data.githubRepo);
			} catch (err) {
				nm?.show({
					content: 'Failed to update github repo url',
					type: NotificationType.Error,
				});
			}
			SettingUpProgress(false);
		};

		return (
			<>
				<ZLoaderCover loading={settingUp}>
					<section className={cls('')}>
						<div className={`${cls('_content')}`}>
							<form ref={formRef} className={`${cls('_form')}`}>
								<h6 className={`${cls('_heading')}`}>Step 1.</h6>
								<section className={`${cls('_form-group')}`}>
									<em>
										Create a github repo for eg.
										https://github.com/zl-zbank-tech/zlifecycle-config.git
									</em>
								</section>
								<h6 className={`${cls('_heading')}`}>Step 2.</h6>
								<section className={`${cls('_form-group')}`}>
									{repoUrl ? (
										<span>Github Repo is set to {repoUrl} you can update the URL here.</span>
									) : (
										<label className="required">Paste the repo web url here</label>
									)}
									<input
										type="text"
										pattern="git@.*.git|https://.*.git"
										required
										name="githubRepo"
										className="input"
										placeholder={
											repoUrl || 'https://github.com/zl-zbank-tech/zlifecycle-config.git'
										}
									/>
									<label className="error-msg">Please enter a valid repo URL</label>
									<div className="mt-5">
										<button
											type="button"
											className="shadowy-input btn btn__update"
											onClick={async () => await this.saveGithubCredentials()}>
											Update
										</button>
									</div>
								</section>
								<h6 className={`${cls('_heading')}`}>Step 3.</h6>
								<section className={`${cls('_form-group')}`}>
									<em>
										Install flow using{' '}
										<a href="https://github.com/apps/zlifecycle">
											https://github.com/apps/zlifecycle
										</a>
									</em>
									<label></label>
									<label>OR</label>
									<em>Contact us for other alternative.</em>
								</section>
							</form>
						</div>
					</section>
				</ZLoaderCover>
			</>
		);
	};

	static getInstance() {
		if (!ConfiguringZlifecycle.instance) {
			ConfiguringZlifecycle.instance = new ConfiguringZlifecycle('Configure CloudKnit');
		}
		return ConfiguringZlifecycle.instance;
	}

	async onNext(): Promise<any> {
		if (!this.context) {
			throw 'Please update the github repo url.';
		}
		return this.context;
	}

	render(baseClassName: string, ctx: any, nm?: NotificationsManager) {
		super.show(`.${baseClassName}_section-guide`);
		return <this.ConfiguringZlifecycleUI baseClassName={baseClassName} ctx={ctx} nm={nm} />;
	}
}

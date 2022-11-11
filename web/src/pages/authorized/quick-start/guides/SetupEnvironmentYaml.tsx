import { BaseGuide } from './BaseGuide';
import React, { useEffect, useState } from 'react';
import { IGuide } from './IGuide';
import AuthStore from 'auth/AuthStore';
import { ZEditor } from 'components/molecules/editor/Editor';
import { LocalStorage } from 'utils/localStorage/localStorage';
import { LocalStorageKey } from 'models/localStorage';
import { QuickStartContext } from '../QuickStart';
import { NotificationsManager } from 'components/argo-core';
import { ReactComponent as Copy } from 'assets/images/icons/copy.svg';

type Props = {
	baseClassName: string;
	ctx: any;
	nm?: NotificationsManager;
};

export class SetupEnvironmentYaml extends BaseGuide implements IGuide {
	private static instance: SetupEnvironmentYaml | null = null;
	private team_name = 'default';

	private SetupEnvironmentYamlUI: React.FC<Props> = ({ baseClassName, ctx, nm }) => {
		const cls = (className: string) => `${baseClassName}_section-guide${className}`;
		const [envName, setEnvName] = useState<string>(this.team_name);
		const user = AuthStore.getUser();
		const formRef = React.useRef<HTMLFormElement>(null);
		const envYaml = `apiVersion: stable.compuzest.com/v1
kind: Environment
metadata:
	name: ${user?.organizations[0].name}-hello-world
	namespace: ${user?.organizations[0].name}-config
spec:
	teamName: ${ctx?.teamName}
	envName: ${envName}
	components:
	- name: images
		type: terraform
		module:
			source: aws
			name: s3-bucket
		variables:
			- name: bucket
			  value: "${user?.organizations[0].name}-${envName}-images"
	- name: videos
		type: terraform
		dependsOn: [images]
		module:
			source: aws
			name: s3-bucket
		variables:
			- name: bucket
			  value: "${user?.organizations[0].name}-${envName}-videos"
			`;

		return (
			<section className={cls('')}>
				<h4 className={`${cls('_heading')}`}>Lets setup an Environment</h4>
				<div className={`${cls('_content')}`}>
					<form ref={formRef} className={`${cls('_form')}`}>
						<section className={`${cls('_form-group')}`}>
							<label className="required">Environment Name</label>
							<input
								type="text"
								required
								name="envName"
								className="input"
								placeholder="enter an env name"
								defaultValue={envName}
								onChange={e => {
									const val = e.target.value;
									if (super.rfcSubdomainValidation(val)) {
										setEnvName(val);
										formRef.current?.classList.remove('invalid');
									} else {
										formRef.current?.classList.add('invalid');
									}
								}}
							/>
						</section>
						<section className={`${cls('_form-group')}`}>
							<label>
								Copy below yaml and push it to zlifecycle-config repo under environments folder{' '}
								<button
									type="button"
									title="Copy YML"
									className="copy-to-clipboard"
									onClick={() => {
										this.copyToClipboard(envYaml, nm);
									}}>
									<Copy />
								</button>
							</label>
							<div>
								<ZEditor data={envYaml} language={'yaml'} readOnly={true} />
							</div>
						</section>
					</form>
				</div>
			</section>
		);
	};

	static getInstance() {
		if (!SetupEnvironmentYaml.instance) {
			SetupEnvironmentYaml.instance = new SetupEnvironmentYaml('Environment Setup');
		}
		return SetupEnvironmentYaml.instance;
	}

	onNext(): Promise<any> {
		return Promise.resolve({});
	}

	onFinish(): Promise<any> {
		LocalStorage.setItem<QuickStartContext>(LocalStorageKey.QUICK_START_STEP, { ctx: {}, step: -1 });
		window.location.href = `${process.env.REACT_APP_BASE_URL}`;
		return Promise.resolve({});
	}

	render(baseClassName: string, ctx: any, nm?: NotificationsManager) {
		super.show(`.${baseClassName}_section-guide`);
		return <this.SetupEnvironmentYamlUI baseClassName={baseClassName} ctx={ctx} nm={nm} />;
	}
}

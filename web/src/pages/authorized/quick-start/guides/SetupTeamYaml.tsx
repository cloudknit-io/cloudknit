import { BaseGuide } from './BaseGuide';
import React, { useEffect, useState } from 'react';
import { IGuide } from './IGuide';
import AuthStore from 'auth/AuthStore';
import { ZEditor } from 'components/molecules/editor/Editor';
import { ReactComponent as Copy } from 'assets/images/icons/copy.svg';
import { NotificationsManager } from 'components/argo-core';
import { LocalStorage } from 'utils/localStorage/localStorage';
import { QuickStartContext } from '../QuickStart';
import { LocalStorageKey } from 'models/localStorage';

type Props = {
	baseClassName: string;
	ctx: any;
	nm?: NotificationsManager;
};

export class SetupTeamYaml extends BaseGuide implements IGuide {
	private static instance: SetupTeamYaml | null = null;
	private team_name = 'default';
	private env_name = 'default';

	private SetupTeamYamlUI: React.FC<Props> = ({ baseClassName, ctx, nm }) => {
		const cls = (className: string) => `${baseClassName}_section-guide${className}`;
		const [teamName, setTeamName] = useState<string>(ctx?.teamName || this.team_name);
		const [envName, setEnvName] = useState<string>(ctx?.envName || this.env_name);
		const user = AuthStore.getUser();
		const formRef = React.useRef<HTMLFormElement>(null);
		const teamYaml = `apiVersion: stable.compuzest.com/v1
kind: Team
metadata:
  name: ${teamName}
  namespace: ${user?.selectedOrg.name}-config
spec:
  teamName: ${teamName}
  configRepo:
    source: ${ctx?.githubRepo || user?.selectedOrg.githubRepo}
    path: .

---

apiVersion: stable.compuzest.com/v1
kind: Environment
metadata:
  name: ${user?.organizations[0].name}-hello-world
  namespace: ${user?.organizations[0].name}-config
spec:
  teamName: ${teamName}
  envName: ${envName}
  components:
    - name: images
      type: terraform
      module:
        source: aws
        name: s3-bucket
      variables:
        - name: bucket
          value: "${user?.selectedOrg.name}-${envName}-images"
    - name: videos
      type: terraform
      dependsOn: [images]
      module:
        source: aws
        name: s3-bucket
      variables:
        - name: bucket
          value: "${user?.selectedOrg.name}-${envName}-videos"
`;

		return (
			<section className={cls('')}>
				<h4 className={`${cls('_heading')}`}>Provision First Environment</h4>
				<div className={`${cls('_content')}`}>
					<form ref={formRef} className={`${cls('_form')}`}>
						<section className={`${cls('_form-group')}`}>
							<label className="required">Team Name</label>
							<input
								type="text"
								required
								name="teamName"
								className="input"
								placeholder="enter a team name"
								defaultValue={teamName}
								onChange={e => {
									const val = e.target.value;
									if (this.rfcSubdomainValidation(val)) {
										setTeamName(val);
										this.team_name = val;
										formRef.current?.classList.remove('invalid');
									} else {
										formRef.current?.classList.add('invalid');
									}
								}}
							/>
							<em className="error-msg">
								Team name must consist of lower case alphanumeric characters, '-' or '.', and must start
								and end with an alphanumeric character
							</em>
						</section>
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
										this.env_name = val;
										formRef.current?.classList.remove('invalid');
									} else {
										formRef.current?.classList.add('invalid');
									}
								}}
							/>
							<em className="error-msg">
								Env name must consist of lower case alphanumeric characters, '-' or '.', and must start
								and end with an alphanumeric character
							</em>
						</section>
						<section className={`${cls('_form-group')}`}>
							<label>
								<span className="break">
									Add a hello-world.yaml file with following content in the root directory of{' '}
									<pre>{ctx?.githubRepo}</pre> repo
								</span>
							</label>
							<label className="mt-5">
								<button
									type="button"
									title="Copy YML"
									className="copy-to-clipboard"
									onClick={() => {
										this.copyToClipboard(teamYaml, nm);
									}}>
									<Copy />
								</button>
							</label>

							<div className="mt-10">
								<ZEditor height='30vh' data={teamYaml} readOnly={true} language={'yaml'} />
							</div>
						</section>
					</form>
				</div>
			</section>
		);
	};

	static getInstance() {
		if (!SetupTeamYaml.instance) {
			SetupTeamYaml.instance = new SetupTeamYaml('Environment Setup');
		}
		return SetupTeamYaml.instance;
	}

	onNext(): Promise<any> {
		return Promise.resolve({});
	}

	async onFinish(): Promise<any> {
		LocalStorage.setItem<QuickStartContext>(LocalStorageKey.QUICK_START_STEP, { ctx: {}, step: 0 });
		const org = await AuthStore.fetchOrganizationStatus();
		if (org.provisioned) {
			AuthStore.redirectToHome();
		} else {
			window.location.href = `${process.env.REACT_APP_BASE_URL}`;
			return Promise.resolve({});
		}
	}

	render(baseClassName: string, ctx: any, nm?: NotificationsManager) {
		super.show(`.${baseClassName}_section-guide`);
		return <this.SetupTeamYamlUI baseClassName={baseClassName} ctx={ctx} nm={nm} />;
	}
}

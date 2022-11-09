import { BaseGuide } from './BaseGuide';
import React, { useEffect, useState } from 'react';
import { IGuide } from './IGuide';
import AuthStore from 'auth/AuthStore';
import { ZEditor } from 'components/molecules/editor/Editor';

type Props = {
	baseClassName: string;
	ctx: any;
};


export class SetupTeamYaml extends BaseGuide implements IGuide {
	private static instance: SetupTeamYaml | null = null;
	private team_name = 'default';

	private SetupTeamYamlUI: React.FC<Props> = ({ baseClassName, ctx }) => {
		const cls = (className: string) => `${baseClassName}_section-guide${className}`;
		const [teamName, setTeamName] = useState<string>(ctx.teamName || this.team_name);
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
	path: environments
`;

		return (
			<section className={cls('')}>
				<h4 className={`${cls('_heading')}`}>Lets setup a team</h4>
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
							<em className='error-msg'>Team name must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character</em>
						</section>
						<section className={`${cls('_form-group')}`}>
							<label>Copy below yaml and push it to zlifecycle-config repo under teams</label>
							<div>
								<ZEditor
									data={teamYaml}
									readOnly={true}
									language={'yaml'}
								/>
							</div>
						</section>
					</form>
				</div>
			</section>
		);
	};

	static getInstance() {
		if (!SetupTeamYaml.instance) {
			SetupTeamYaml.instance = new SetupTeamYaml('Team Setup');
		}
		return SetupTeamYaml.instance;
	}

	onNext(): Promise<any> {
		return Promise.resolve({
			teamName: this.team_name,
		});
	}

	render(baseClassName: string, ctx: any) {
		super.show(`.${baseClassName}_section-guide`);
		return <this.SetupTeamYamlUI baseClassName={baseClassName} ctx={ctx} />;
	}
}

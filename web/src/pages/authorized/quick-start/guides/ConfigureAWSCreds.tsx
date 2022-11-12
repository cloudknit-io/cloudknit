import { BaseGuide } from './BaseGuide';
import React from 'react';
import { IGuide } from './IGuide';
import { AWSSecret } from 'pages/authorized/secrets/aws-secrets';
import AuthStore from 'auth/AuthStore';
import '../../profile/styles.scss';

type Props = {
	baseClassName: string;
};

export class ConfigureAWSCreds extends BaseGuide implements IGuide {
	private static instance: ConfigureAWSCreds | null = null;
	private awsSecretUpdated = false;

	private ConfigureAWSCredsUI: React.FC<Props> = ({ baseClassName }) => {
		const cls = (className: string) => `${baseClassName}_section-guide${className}`;
		return (
			<section className={cls('')}>
				<h4 className={`${cls('_heading')}`}>Provide AWS Credentials</h4>
				<div className={`${cls('_content')}`}>
					<form>
						<section className={`${cls('_form-group')}`}>
							<label>Note: Your environment will be provisioned on this AWS account.</label>
							<div className="secrets-container">
								<div className={`secret-info secret-info-ssm`}>
									<div className="secret-container">
										<AWSSecret
											existCallback={(exists) => {
												this.awsSecretUpdated = exists;
											}}
											closeCallback={() => {
												this.awsSecretUpdated = true;
											}}
											secretScope={AuthStore.getUser()?.selectedOrg.name || ''}
										/>
									</div>
								</div>
							</div>
						</section>
					</form>
				</div>
			</section>
		);
	};

	static getInstance() {
		if (!ConfigureAWSCreds.instance) {
			ConfigureAWSCreds.instance = new ConfigureAWSCreds('Configure AWS');
		}
		return ConfigureAWSCreds.instance;
	}

	async onNext(): Promise<any> {
		if (!this.awsSecretUpdated) {
			throw 'Please update aws secrets before proceeding';
		}
		return {};
	}

	render(baseClassName: string, ctx: any) {
		super.show(`.${baseClassName}_section-guide`);
		return <this.ConfigureAWSCredsUI baseClassName={baseClassName} />;
	}
}

import { BaseGuide } from './BaseGuide';
import React, { useEffect } from 'react';
import { IGuide } from './IGuide';
import { AWSSecret } from 'pages/authorized/secrets/aws-secrets';
import AuthStore from 'auth/AuthStore';
import { Subject } from 'rxjs';
import '../../profile/styles.scss';

type Props = {
	baseClassName: string;
};

export class ConfigureAWSCreds extends BaseGuide implements IGuide {
	private static instance: ConfigureAWSCreds | null = null;
	private awsSecretUpdated = false;
	private triggerSave = new Subject<void>();
	private promiseResolver: any = null;
	private saveSecret = async () => {
		const promise = new Promise<boolean>((r, e) => {
			this.promiseResolver = r;
		});
		this.triggerSave.next();
		return promise;
	}
	private saveCallback = (success: boolean) => {
		if (this.promiseResolver) {
			this.promiseResolver(success);
		}
	}

	private ConfigureAWSCredsUI: React.FC<Props> = ({ baseClassName }) => {
		const cls = (className: string) => `${baseClassName}_section-guide${className}`;

		return (
			<section className={cls('')}>
				<h4 className={`${cls('_heading')}`}>Provide AWS Credentials</h4>
				<div className={`${cls('_content')}`}>
					<form>
						<section className={`${cls('_form-group')} ${cls('_form-group-flex')}`}>
							<label>Note: Your environment will be provisioned on this AWS account.</label>
							<div>
							 Check <a href="https://docs.aws.amazon.com/powershell/latest/userguide/pstools-appendix-sign-up.html" target='_blank'>this</a> article to find out how to get following keys.
							</div>
							<div className="secrets-container">
								<div className={`secret-info secret-info-ssm`}>
									<div className="secret-container">
										<AWSSecret
											existCallback={(exists) => {
												this.awsSecretUpdated = exists;
											}}
											closeCallback={this.saveCallback}
											secretScope={AuthStore.getUser()?.selectedOrg.name || ''}
											externalSaveTrigger={this.triggerSave}
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
			this.awsSecretUpdated = await this.saveSecret();
		}
		if (!this.awsSecretUpdated) {
			return null;
		}
		return {};
	}

	render(baseClassName: string, ctx: any) {
		super.show(`.${baseClassName}_section-guide`);
		return <this.ConfigureAWSCredsUI baseClassName={baseClassName} />;
	}
}

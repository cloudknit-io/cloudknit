import { BaseGuide } from './BaseGuide';
import React, { useEffect, useState } from 'react';
import { IGuide } from './IGuide';
import { AWSSecret } from 'pages/authorized/secrets/aws-secrets';
import AuthStore from 'auth/AuthStore';
import { Subject } from 'rxjs';
import '../../profile/styles.scss';
import { SecretsService } from 'services/secrets/secrets.service';

type Props = {
	baseClassName: string;
};

export class ConfigureAWSCreds extends BaseGuide implements IGuide {
	private static instance: ConfigureAWSCreds | null = null;
	private awsSecretUpdated = false;
	private triggerSave = new Subject<void>();
	private promiseResolver: any = null;
	private credType: 0 | 1 = 0;
	private saveSecret = async () => {
		if (this.credType === 0) {
			const resp = await SecretsService.getInstance().setDefaultSsmSecret();
			return resp.data == true;
		}
		const promise = new Promise<boolean>((r, e) => {
			this.promiseResolver = r;
		});
		this.triggerSave.next();
		return promise;
	};
	private saveCallback = (success: boolean) => {
		if (this.promiseResolver) {
			this.promiseResolver(success);
		}
	};

	private ConfigureAWSCredsUI: React.FC<Props> = ({ baseClassName }) => {
		const cls = (className: string) => `${baseClassName}_section-guide${className}`;
		// 0: Cloudknit, 1: Customer
		const [cred, setCred] = useState<1 | 0>(this.credType);
		return (
			<section className={cls('')}>
				<h4 className={`${cls('_heading')}`}>Provide AWS Credentials</h4>
				<div className={`${cls('_content')}`}>
					<form>
						<section className={`${cls('_form-group')} ${cls('_form-group-flex')}`}>
							<label>Note: Your environment will be provisioned on this AWS account.</label>
							<div>
								Check{' '}
								<a
									href="https://docs.aws.amazon.com/powershell/latest/userguide/pstools-appendix-sign-up.html"
									target="_blank">
									this
								</a>{' '}
								article to find out how to get following keys.
							</div>
							<div className="options-container">
								<div className="options-container__option">
									<input
										checked={this.credType === 0}
										id="cred-type-0"
										name="cred-type"
										type="radio"
										onChange={() => {
											this.credType = 0;
											setCred(0);
										}}
									/>
									<label htmlFor="cred-type-0">Use Cloudknit Credentials</label>
								</div>
								<div className="options-container__option">
									<input
										id="cred-type-1"
										name="cred-type"
										type="radio"
										onChange={() => {
											this.credType = 1;
											setCred(1);
										}}
									/>
									<label htmlFor="cred-type-1">Use your own Credentials</label>
								</div>
							</div>
							{cred === 1 && (
								<div className="secrets-container">
									<div className={`secret-info secret-info-ssm`}>
										<div className="secret-container">
											<AWSSecret
												existCallback={exists => {
													this.awsSecretUpdated = exists;
												}}
												closeCallback={this.saveCallback}
												secretScope={AuthStore.getUser()?.selectedOrg.name || ''}
												externalSaveTrigger={this.triggerSave}
											/>
										</div>
									</div>
								</div>
							)}
							{cred === 0 && (
								<div className='cred-type-0-info'>
									* Use CloudKnit AWS credentials to try out the product. These credentials will be
									disabled after 5 days.
								</div>
							)}
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

import { BaseGuide } from './BaseGuide';
import React, { useState } from 'react';
import { IGuide } from './IGuide';
import ApiClient from 'utils/apiClient';
import { ENVIRONMENT_VARIABLES } from 'utils/environmentVariables';
import { Loader } from 'components/atoms/loader/Loader';

type Props = {
	baseClassName: string;
	ctx: any;
};

export class OAuthGuide extends BaseGuide implements IGuide {
	private static instance: OAuthGuide | null = null;
	private saveCredentials?: () => Promise<any>;

	private OAuthGuideUI: React.FC<Props> = ({ baseClassName, ctx }) => {
		const cls = (className: string) => `${baseClassName}_section-guide${className}`;
		const formRef = React.useRef<HTMLFormElement>(null);
		const [settingUp, SettingUpProgress] = useState<boolean>(false);

		this.saveCredentials = async (): Promise<any> => {
			if (!formRef.current) return false;
			if (!formRef.current.checkValidity()) {
				formRef.current.classList.add('invalid');
				return false;
			}
			formRef.current.classList.remove('invalid');
			const formData = new FormData(formRef.current);
			SettingUpProgress(true);
			const payload = {
				company: formData.get('company'),
				clientId: formData.get('clientId'),
				clientSecret: formData.get('clientSecret'),
				namespace: `${ENVIRONMENT_VARIABLES.REACT_APP_CUSTOMER_NAME}-system`,
			};
			try {
				const res = await ApiClient.post('/company/oath/credentials', payload);
				console.log(res);
				const res2 = await ApiClient.patch('/company/oath/credentials', payload);
				console.log(res2);
				SettingUpProgress(false);
				return res.data;
			} catch (err) {
				SettingUpProgress(false);
				return false;
			}
		};

		return (
			<section className={cls('')}>
				<h4 className={`${cls('_heading')}`}>
					This will allow you to login to your zLifecycle instance using github oAuth
				</h4>
				<div className={`${cls('_content')}`}>
					{settingUp ? (
						<div
							className="display-flex justify-center align-center"
							style={{ flexDirection: 'column', marginTop: '100px' }}>
							<Loader title="Setting up OAuth" height={32} width={32} />
							<h3>Setting up OAuth</h3>
						</div>
					) : (
						<form className={`${cls('_form')}`} onSubmit={e => e.preventDefault()} ref={formRef}>
							<section className={`${cls('_form-group')}`}>
								<label className="required">Github Org</label>
								<input
									name="company"
									required
									className="input"
									placeholder="Please enter your github org id"
								/>
								<label className="error-msg">github org id cannot be empty</label>
							</section>
							<section className={`${cls('_form-group')}`}>
								<label className="required">Github oAuth Client Id</label>
								<input
									name="clientId"
									required
									className="input"
									placeholder="Please enter your client id"
								/>
								<label className="error-msg">client id cannot be empty</label>
							</section>
							<section className={`${cls('_form-group')}`}>
								<label className="required">Github oAuth Client Secret</label>
								<input
									name="clientSecret"
									required
									className="input"
									placeholder="Please enter your client secret"
								/>
								<label className="error-msg">client secret cannot be empty</label>
							</section>
						</form>
					)}
				</div>
			</section>
		);
	};

	static getInstance() {
		if (!OAuthGuide.instance) {
			OAuthGuide.instance = new OAuthGuide('GitHub Auth');
		}
		return OAuthGuide.instance;
	}

	async onNext(): Promise<any> {
		if (!this.saveCredentials) return false;
		return this.saveCredentials();
	}

	render(baseClassName: string, ctx: any) {
		super.show(`.${baseClassName}_section-guide`);
		return <this.OAuthGuideUI baseClassName={baseClassName} ctx={ctx} />;
	}
}

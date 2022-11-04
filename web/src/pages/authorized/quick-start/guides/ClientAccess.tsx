import { BaseGuide } from './BaseGuide';
import React from 'react';
import { IGuide } from './IGuide';

type Props = {
	baseClassName: string;
};

export class ClientAccessGuide extends BaseGuide implements IGuide {
	private static instance: ClientAccessGuide | null = null;

	private ClientAccessGuideUI: React.FC<Props> = ({ baseClassName }) => {
		const cls = (className: string) => `${baseClassName}_section-guide${className}`;
		return (
			<section className={cls('')}>
				<h4 className={`${cls('_heading')}`}>Steps to allow zlifecycle to access your Github account</h4>
				<div className={`${cls('_content')}`}>
					<form>
						<section className={`${cls('_form-group')}`}>
							<em>
								-- Install flow using{' '}
								<a href="https://github.com/apps/zlifecycle">https://github.com/apps/zlifecycle</a>
							</em>
							<label></label>
							<label>OR</label>
							<em>Contact us for other alternative.</em>
						</section>
					</form>
				</div>
			</section>
		);
	};

	static getInstance() {
		if (!ClientAccessGuide.instance) {
			ClientAccessGuide.instance = new ClientAccessGuide('Client Access');
		}
		return ClientAccessGuide.instance;
	}

	render(baseClassName: string, ctx: any) {
		super.show(`.${baseClassName}_section-guide`);
		return <this.ClientAccessGuideUI baseClassName={baseClassName} />;
	}
}

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
import { CalMeet } from 'pages/authorized/dashboard/helpers';

type Props = {
	baseClassName: string;
	ctx: any;
	nm?: NotificationsManager;
};

export class SetupCalMeet extends BaseGuide implements IGuide {
	private static instance: SetupCalMeet | null = null;

	private SetupCalMeet: React.FC<Props> = ({ baseClassName, ctx, nm }) => {
		const cls = (className: string) => `${baseClassName}_section-guide${className}`;

		return (
			<section className={cls('')}>
				<h4 className={`${cls('_heading')}`}>Schedule Onboarding Meeting</h4>
				<div className={`${cls('_content')}`}>
					<CalMeet />
				</div>
			</section>
		);
	};

	static getInstance() {
		if (!SetupCalMeet.instance) {
			SetupCalMeet.instance = new SetupCalMeet('Schedule Meet');
		}
		return SetupCalMeet.instance;
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
		return <this.SetupCalMeet baseClassName={baseClassName} ctx={ctx} nm={nm} />;
	}
}

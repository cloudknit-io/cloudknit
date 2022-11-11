import { NotificationsManager, NotificationType } from 'components/argo-core';
import { uniqueId } from 'lodash';
import React from 'react';
import { IGuide } from './IGuide';

export class BaseGuide implements IGuide {
	stepId: string = uniqueId('step');
	stepName: string = uniqueId('step');

	private UI: React.FC<{}> = () => {
		return <></>;
	};

	constructor(stepName: string) {
		this.stepName = stepName;
	}

	rfcSubdomainValidation(testString: string) {
		const pattern = /^(?![0-9]+$)(?!.*-$)(?!-)[a-z0-9-]{1,63}$/g;
		return pattern.test(testString);
	}

	render(baseClassName = '', ctx: any) {
		return <this.UI />;
	}

	show(selector: string) {
		setImmediate(() => {
			document.querySelector(selector)?.classList.add('visible');
		});
	}

	copyToClipboard(content: string, nm?: NotificationsManager) {
		navigator.clipboard.writeText(content);
		nm?.show({
			content: 'Copied',
			type: NotificationType.Success,
			toastOptions: {
				autoClose: 1000,
			},
		});
	}
}

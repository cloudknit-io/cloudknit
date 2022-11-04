import { NotificationsManager } from 'components/argo-core';
import React, { useEffect, useState } from 'react';
import { guideIndex } from './guides';
import { IGuide } from './guides/IGuide';

export type Props = {
	stepId: string;
	ctx: any;
	nm: NotificationsManager
};

export const QuickStartContent: React.FC<Props> = ({ stepId, ctx, nm }) => {
	const guide = guideIndex.get(stepId) as IGuide;

	return <>{guide.render('quick-start-guide_container_right-container', ctx, nm)}</>;
};

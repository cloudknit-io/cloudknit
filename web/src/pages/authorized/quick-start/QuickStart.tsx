import { Notifications, NotificationsManager, NotificationType } from 'components/argo-core';
import { Context } from 'context/argo/ArgoUi';
import { LocalStorageKey } from 'models/localStorage';
import React, { useContext, useEffect, useState } from 'react';
import { LocalStorage } from 'utils/localStorage/localStorage';
import { guideIndex, guideKeys } from './guides';
import { QuickStartContent } from './QuickStartContent';
import { QuickStartIndex } from './QuickStartIndex';
import './styles.scss';

export interface QuickStartContext {
	ctx: any;
	step: number;
}

export const QuickStart: React.FC = () => {
	const [activeStepIndex, updateActiveStepIndex] = useState<number>(
		LocalStorage.getItem<QuickStartContext>(LocalStorageKey.QUICK_START_STEP)?.step || 0
	);
	const [next, nextInProgress] = useState<boolean>(false);
	const [ctx, updateCtx] = useState<any>(
		LocalStorage.getItem<QuickStartContext>(LocalStorageKey.QUICK_START_STEP)?.ctx
	);
	const nm: NotificationsManager = new NotificationsManager();

	useEffect(() => {
		LocalStorage.setItem<QuickStartContext>(LocalStorageKey.QUICK_START_STEP, { ctx, step: activeStepIndex });
	}, [activeStepIndex]);

	return (
		<div className="quick-start-guide">
			<h1 className="quick-start-guide_heading-main">
				<span className="quick-start-guide_heading-main_z">z</span>Lifecycle Setup Wi
				<span className="quick-start-guide_heading-main_z">z</span>ard
			</h1>
			<div className="quick-start-guide_container">
				<div className="quick-start-guide_container_left-container">
					<QuickStartIndex
						activeStepIndex={activeStepIndex}
						changeHandler={(index: number) => {
							// updateActiveStepIndex(index);
						}}
					/>
				</div>
				<div className="quick-start-guide_container_right-container">
					<QuickStartContent stepId={guideKeys[activeStepIndex]} ctx={ctx} nm={nm} />
					<div className="quick-start-guide_container_right-container_footer">
						<div className="quick-start-guide_container_right-container_progress-bar">
							<span>
								<span
									className="progress"
									style={{ width: `${(activeStepIndex / (guideKeys.length - 1)) * 100}%` }}></span>
							</span>
						</div>
						<button
							disabled={activeStepIndex === 0}
							onClick={() => {
								updateActiveStepIndex(activeStepIndex - 1);
							}}>
							Previous
						</button>
						<button
							onClick={async () => {
								const guide = guideIndex.get(guideKeys[activeStepIndex]);
								if (!guide) return;
								if (!guide.onNext) {
									if (activeStepIndex === guideKeys.length - 1) {
										await guide.onFinish?.call(guide);
									} else {
										updateActiveStepIndex(activeStepIndex + 1);
									}
									return;
								}
								try {
									nextInProgress(true);
									const res = await guide.onNext();
									if (res) {
										updateCtx({
											...ctx,
											...res,
										});
										if (activeStepIndex === guideKeys.length - 1) {
											await guide.onFinish?.call(guide);
										} else {
											updateActiveStepIndex(activeStepIndex + 1);
										}
									}
									nextInProgress(false);
								} catch (err) {
									nm?.show({
										content: err,
										type: NotificationType.Error,
									});
								}
							}}>
							{activeStepIndex === guideKeys.length - 1 ? 'Finish' : 'Next'}
						</button>
					</div>
				</div>
			</div>
			<Notifications notifications={nm.notifications} />
		</div>
	);
};

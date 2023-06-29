import { Notifications, NotificationsManager, NotificationType } from 'components/argo-core';
import { Loader } from 'components/atoms/loader/Loader';
import { Context } from 'context/argo/ArgoUi';
import { LocalStorageKey } from 'models/localStorage';
import React, { useContext, useEffect, useState } from 'react';
import { LocalStorage } from 'utils/localStorage/localStorage';
import { usePageHeader } from '../contexts/EnvironmentHeaderContext';
import { guideIndex, guideKeys } from './guides';
import { QuickStartContent } from './QuickStartContent';
import { QuickStartIndex } from './QuickStartIndex';
import './styles.scss';
import { calLink, CalMeet } from '../dashboard/helpers';
import { ZLoaderCover } from 'components/atoms/loader/LoaderCover';

export interface QuickStartContext {
	ctx: any;
	step: number;
}
export const QuickStart: React.FC = () => {
	const [activeStepIndex, updateActiveStepIndex] = useState<number>(0);
	const [next, nextInProgress] = useState<boolean>(false);
	const [scheduleMeet, setScheduleMeet] = useState<boolean>(true);
	const [calLoading, setCalLoading] = useState<boolean>(true);
	const [ctx, updateCtx] = useState<any>(
		LocalStorage.getItem<QuickStartContext>(LocalStorageKey.QUICK_START_STEP)?.ctx
	);
	const nm: NotificationsManager = new NotificationsManager();

	const { pageHeaderObservable, breadcrumbObservable } = usePageHeader();

	useEffect(() => {
		pageHeaderObservable.next({
			breadcrumbs: [],
			headerTabs: [],
			pageName: null,
			filterTitle: '',
			onSearch: () => {},
			buttonText: '',
			onViewChange: () => {},
		});

		window.Cal('on', {
			action: '__windowLoadComplete',
			callback: e => {
				setCalLoading(false);
			},
		});
	}, []);

	useEffect(() => {
		breadcrumbObservable.next(false);
	}, [breadcrumbObservable]);

	useEffect(() => {
		LocalStorage.setItem<QuickStartContext>(LocalStorageKey.QUICK_START_STEP, { ctx, step: activeStepIndex });
	}, [activeStepIndex]);

	return (
		<ZLoaderCover loading={calLoading}>
			<div className="quick-start-guide">
				<h1 className="quick-start-guide_heading-main">
					{scheduleMeet ? (
						'Setup Onboarding Meeting'
					) : (
						<>
							<span className="quick-start-guide_heading-main_z">Cloud</span>Knit Setup Wizard
						</>
					)}
				</h1>
				{scheduleMeet ? (
					<div className="schedule-meeting">
						<CalMeet />
						<button className="btn manual-setup" onClick={() => setScheduleMeet(false)}>
							Skip & Setup Yourself
						</button>
					</div>
				) : (
					<div className="quick-start-guide_container">
						<div className="quick-start-guide_container_left-container">
							<QuickStartIndex
								activeStepIndex={activeStepIndex}
								changeHandler={(index: number) => {
									updateActiveStepIndex(index);
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
											style={{
												width: `${(activeStepIndex / (guideKeys.length - 1)) * 100}%`,
											}}></span>
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
									disabled={next}
									onClick={async () => {
										if (next) {
											return;
										}
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
										} catch (err) {
											nm?.show({
												content: err,
												type: NotificationType.Error,
											});
										}
										nextInProgress(false);
									}}>
									{next ? (
										<Loader height={16} />
									) : activeStepIndex === guideKeys.length - 1 ? (
										'Finish'
									) : (
										'Next'
									)}
								</button>
							</div>
						</div>
					</div>
				)}
				<Notifications notifications={nm.notifications} />
			</div>
		</ZLoaderCover>
	);
};

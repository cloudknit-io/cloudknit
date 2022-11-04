import React, { useEffect, useState } from 'react';
import { guideIndex, guideKeys, guideValues } from './guides';
import { IGuide } from './guides/IGuide';

export type Props = {
	changeHandler: (index: number) => any;
    activeStepIndex: number;
};

export const QuickStartIndex: React.FC<Props> = ({ changeHandler, activeStepIndex }) => {
	const activeStep = guideValues[activeStepIndex].stepId;

	const cls = (className: string) => `quick-start-guide_container_left-container${className}`;

	const indexDOM = (step: IGuide) => {
		return (
			<div key={step.stepId} className={`${cls('_index_step-container')}`}>
				<button
					className={`${cls('_index_step-container_button')} ${
						activeStep === step.stepId ? cls('_index_step-container_button-active') : ''
					}`}
					onClick={e => {
						changeHandler(guideKeys.indexOf(step.stepId));
					}}>
					{step.stepName}
				</button>
			</div>
		);
	};

	return <div className={cls('_index')}>{guideValues.map(step => indexDOM(step))}</div>;
};

import './style.scss';

import classNames from 'classnames';
import React, { FC } from 'react';
export type LabelColors = 'orange' | 'light-green' | 'pink' | 'blue' | 'violet';

type Props = {
	color: LabelColors;
	text: string;
};

export const ZCardLabel: FC<Props> = ({ color, text }: Props) => {
	return (
		<div className={classNames('com-card-label', `com-card-label--${color}`)} title={text}>
			{text}
		</div>
	);
};

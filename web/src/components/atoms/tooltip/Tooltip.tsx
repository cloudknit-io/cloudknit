import './style.scss';

import Tippy from '@tippy.js/react';
import React, { FC, PropsWithChildren } from 'react';

export type Props = PropsWithChildren<any>;

export const ZTooltip: FC<Props> = props => {
	return (
		<Tippy animation="fade" arrow={true} className={'zlifecycle-tooltip'} {...props}>
			{props.children}
		</Tippy>
	);
};

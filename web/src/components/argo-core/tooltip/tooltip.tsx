import Tippy from '@tippy.js/react';
import React from 'react';

export const Tooltip = (props: any) => (
	<Tippy animation="fade" arrow="true" {...props}>
		{props.children}
	</Tippy>
);

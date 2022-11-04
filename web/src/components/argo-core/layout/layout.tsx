import './layout.scss';

import React from 'react';

export interface LayoutProps extends React.Props<any> {
	version?: () => React.ReactElement;
}

export const Layout = (props: LayoutProps) => {
	return <div className="layout">{props.children}</div>;
};

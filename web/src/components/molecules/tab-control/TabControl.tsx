import './style.scss';

import classNames from 'classnames';
import { ZTabs } from 'components/atoms/tabs/Tabs';
import { OptionItem } from 'models/general.models';
import React, { FC, PropsWithChildren, useEffect, useState } from 'react';

interface Props extends PropsWithChildren<any> {
	tabs?: OptionItem[];
	selected: string;
	onTabChange?: Function;
	className?: string;
}

export const ZTablControl: FC<Props> = (props: Props) => {
	const { tabs, children, selected, onTabChange, className = '' } = props;
	const [currentTab, setCurrentTab] = useState<string>(selected);

	useEffect(() => {
		setCurrentTab(selected);
	}, [selected]);

	const handleTabChange = (id: string): void => {
		setCurrentTab(id);
		onTabChange && onTabChange(id);
	};

	return (
		<div className={classNames('zlifecycle-tab-control', className)}>
			<nav className="zlifecycle-tab-control__tabs">
				{tabs && <ZTabs items={tabs} active={currentTab} onTabSelected={handleTabChange} />}
			</nav>
			<div className="zlifecycle-tab-control__tab-pane">
				{children.map((child: any) => {
					if (child?.props?.id !== currentTab) return undefined;
					return (
						<div key={child?.props?.id} id={child?.props?.id}>
							{child?.props?.children}
						</div>
					);
				})}
			</div>
		</div>
	);
};

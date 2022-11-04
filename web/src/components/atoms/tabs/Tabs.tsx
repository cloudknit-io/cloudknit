import './style.scss';

import classNames from 'classnames';
import { ZText } from 'components/atoms/text/Text';
import { OptionItem } from 'models/general.models';
import React, { FC } from 'react';

type Props = {
	items: OptionItem[];
	active: string;
	onTabSelected: Function;
};

export const ZTabs: FC<Props> = ({ items, active, onTabSelected }) => {
	if (items.length > 0 && !items.find(e => e.id === active)) {
		onTabSelected(items[0].id)
	}

	return (
		<>
			{items.map(tab => (
				<div
					key={`tab-${tab.id}`}
					className={classNames('nav-link', {
						'nav-link--active': tab.id === active,
					})}
					id={tab.id}
					onClick={(): void => onTabSelected(tab.id)}>
					<ZText.Body size="16" lineHeight="24" weight={active === tab.id ? 'bold' : 'regular'}>
						{tab.name}
					</ZText.Body>
					{tab.id === active && <div className="nav-link__active-tab" />}
				</div>
			))}
		</>
	);
};

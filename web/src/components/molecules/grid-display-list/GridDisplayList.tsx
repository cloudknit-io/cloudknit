import './style.scss';

import { ZText } from 'components/atoms/text/Text';
import { ZTooltip } from 'components/atoms/tooltip/Tooltip';
import { ListItem } from 'models/general.models';
import React, { FC } from 'react';

type Props = {
	items: ListItem[];
};

const renderTextValueWithLabel = (value: string | undefined, index: number) => {
	return (
		<div key={`display-grid-label-text-value${index}`} className="zlifecycle-grid-display-list__value">
			{!value || value === '' ? (
				<ZText.Body
					key={'list-item-value-' + value + index}
					className="zlifecycle-grid-display-list__value--muted"
					size="14"
					lineHeight="18">
					No data available
				</ZText.Body>
			) : (
				<ZTooltip content={value}>
					<span>
						<ZText.Body key={'list-item-value-' + value + index} size="14" lineHeight="18">
							{value}
						</ZText.Body>
					</span>
				</ZTooltip>
			)}
		</div>
	);
};

export const ZGridDisplayListWithLabel: FC<Props> = ({ items }: Props) => {
	return (
		<div className="zlifecycle-grid-display-list">
			{items.map((item, index) => (
				<React.Fragment key={'list-item-label-' + item.label + index}>
					<ZText.Body className="zlifecycle-grid-display-list__label" size="14" lineHeight="18">
						{`${item.label} ${item.label ? ':' : ''}`}
					</ZText.Body>
					<div key={'display-grid-labels-' + Math.random()}>
						{!item.value || typeof item.value === 'string' ? (
							renderTextValueWithLabel(item.value as any, index)
						) : (
							<div className="zlifecycle-grid-display-list__value ">{item.value}</div>
						)}
					</div>
				</React.Fragment>
			))}
		</div>
	);
};

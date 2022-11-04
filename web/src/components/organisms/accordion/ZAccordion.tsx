import './styles.scss';

import { ReactComponent as Chevron } from 'assets/images/icons/chevron-right.svg';
import React, { FC, useEffect, useState } from 'react';

export interface ZAccordionItem {
	accordionHeader: any;
	accordionContent: any;
	collapsed: boolean;
}

type ZAccordionProps = {
	items: ZAccordionItem[];
};

type ZAccordionPairProps = {
	item: ZAccordionItem;
};

export const ZAccordion: FC<ZAccordionProps> = ({ items }: ZAccordionProps) => {
	return (
		<>
			{items.map((item: ZAccordionItem, _i) => (
				<ZAccordionPair key={_i} item={item} />
			))}
		</>
	);
};

export const ZAccordionPair: FC<ZAccordionPairProps> = ({ item }: ZAccordionPairProps) => {
	const [isCollapsed, setCollapsed] = useState<boolean>(false);
	useEffect(() => {
		if ([true, false].includes(item.collapsed)) {
			setCollapsed(item.collapsed);
		}
	}, [item]);
	return (
		<section className={`zaccordion ${isCollapsed ? 'collapsed' : ''}`}>
			<div className="zaccordion-header" onClick={() => setCollapsed(!isCollapsed)}>
				{item.accordionHeader}
				<div className="chevron">
					<Chevron />
				</div>
			</div>
			<div className="zaccordion-content">{item.accordionContent}</div>
		</section>
	);
};

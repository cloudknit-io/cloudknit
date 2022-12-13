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
	const [localCollapse, setLocalCollapse] = useState<any>(null);
	useEffect(() => {
		if (localCollapse !== null) {
			return;
		}
		if ([true, false].includes(item.collapsed)) {
			setCollapsed(item.collapsed);
		}
	}, [item]);

	const collapsed = () => {
		if (localCollapse === null) {
			return isCollapsed;
		}
		return localCollapse;
	}

	return (
		<section className={`zaccordion ${collapsed() ? 'collapsed' : ''}`}>
			<div className="zaccordion-header" onClick={() => setLocalCollapse(localCollapse == null ? false : !localCollapse)}>
				{item.accordionHeader}
				<div className="chevron">
					<Chevron />
				</div>
			</div>
			<div className="zaccordion-content">{item.accordionContent}</div>
		</section>
	);
};

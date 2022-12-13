import './styles.scss';

import { ReactComponent as ArrowUp } from 'assets/images/icons/arrow_drop_up.svg';
import { ZText } from 'components/atoms/text/Text';
import { MenuItem, ZDropdownMenuJSX } from 'components/molecules/dropdown-menu/DropdownMenu';
import { BreadcrumbFilter } from 'models/breadcrumb.model';
import React, { FC } from 'react';
import { useHistory } from 'react-router';

export interface BreadcrumbProps extends React.Props<any> {
	title: string;
	filters?: BreadcrumbFilter[];
	path: string;
}

export const Breadcrumb: FC<BreadcrumbProps> = ({ title, path, filters }: BreadcrumbProps) => {
	const history = useHistory();

	const pushRoute = (path: string) => {
		history.push(path);
	};

	const filteredItems = (filters || []).filter(e => title.toLowerCase() !== e.title.toLowerCase());

	return (
		<button
			onClick={() => path !== '/all' && pushRoute(path)}
			className={`breadcrumb ${history.location.pathname === path ? 'breadcrumb-active' : ''}`}>
			<ZText.Body size="14" lineHeight="18">
				{title}
			</ZText.Body>
			{filteredItems.length > 0 ? <ArrowUp className="breadcrumb-dd-caret" /> : null}
			{
				<ZDropdownMenuJSX
					className="breadcrumb-ddmenu"
					isOpened={true}
					items={filteredItems.map<MenuItem>(e => ({
						action: (ev: any) => {
							ev.stopPropagation();
							pushRoute(e.routePath);
						},
						text: e.title,
					}))}
				/>
			}
		</button>
	);
};

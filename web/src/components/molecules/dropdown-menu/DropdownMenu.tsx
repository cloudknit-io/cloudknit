import './style.scss';

import classNames from 'classnames';
import React, { FC } from 'react';

export type MenuItem = {
	text: string;
	jsx?: JSX.Element;
	action: Function;
	selected?: boolean;
};

type Props = {
	className?: string;
	items: MenuItem[];
	isOpened: boolean;
	label?: string;
};

export const ZDropdownMenu: FC<Props> = ({ className = '', items, isOpened }) => {
	return (
		<div
			className={classNames(
				'zlifecycle-dropdown-menu',
				className,
				isOpened && 'zlifecycle-dropdown-menu--active'
			)}>
			<ul className="zlifecycle-dropdown-menu__list">
				{items.map((item, index) => (
					<li
						key={`${item.text}-${index}-list-item`}
						className="zlifecycle-dropdown-menu__list__item"
						onClick={() => item.action()}>
						{item.text}
					</li>
				))}
			</ul>
		</div>
	);
};

export const ZDropdownMenuJSX: FC<Props> = ({ className = '', items, isOpened, label }) => {
	return (
		<div
			className={classNames(
				'zlifecycle-dropdown-menu',
				className,
				isOpened && 'zlifecycle-dropdown-menu--active'
			)}>
			<label>{label}</label>
			<ul className="zlifecycle-dropdown-menu__list">
				{items.map((item, index) => (
					<li
						key={`${item.text}-${index}-list-item`}
						className={`zlifecycle-dropdown-menu__list__item ${item.selected && 'zlifecycle-dropdown-menu__list__item--selected'}`}
						onClick={e => item.action(e)}>
						{item.jsx} {item.text}
					</li>
				))}
			</ul>
		</div>
	);
};

import './style.scss';

import { ReactComponent as ClearIcon } from 'assets/images/icons/field/clear.svg';
import classNames from 'classnames';
import { ZText } from 'components/atoms/text/Text';
import React, { FC, PropsWithChildren } from 'react';

interface Props extends PropsWithChildren<any> {
	className?: string;
	isShown: boolean;
	onClose: () => void;
}

export const ZSidePanel: FC<Props> = ({ className = '', isShown, onClose, children }: Props) => {
	return (
		<div className={classNames('zlifecycle-side-panel', isShown && 'zlifecycle-side-panel--active', className)}>
			<div
				className={classNames(
					'zlifecycle-side-panel__shade',
					isShown && 'zlifecycle-side-panel__shade--active'
				)}
				onClick={() => onClose()}></div>
			<div className="zlifecycle-side-panel__content">{children}</div>
		</div>
	);
};

import './style.scss';

import classNames from 'classnames';
import { FC, PropsWithChildren } from 'react';

interface Props extends PropsWithChildren<any> {
	className?: string;
	isShown: boolean;
	header: JSX.Element | string;
	onClose: () => void;
}

export const ZModalPopup: FC<Props> = ({ className = '', isShown, header, onClose, children }: Props) => {
	return (
		<section className={classNames('zlifecycle-modal-popup-overlay', isShown && 'zlifecycle-modal-popup-overlay--active')}>
			<section className={classNames('zlifecycle-modal-popup', className)}>
				<header className="zlifecycle-modal-popup--header">{header}</header>
				<section className="zlifecycle-modal-popup--content">{children}</section>
			</section>
		</section>
	);
};

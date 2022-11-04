import 'components/atoms/button/Button.scss';

import classNames from 'classnames';
import { Loader } from 'components/atoms/loader/Loader';
import React from 'react';

type Props = {
	isLoading?: boolean;
	color?: string;
	block?: boolean;
	className?: string;
	type?: string;
	disabled?: boolean;
	onClick?: () => void;
};

export const Button: React.FC<Props> = ({ className, children, isLoading, color, block, onClick }) => {
	return (
		<button
			onClick={onClick}
			className={classNames(`${className} zlifecycle-btn argo-button--${color || 'primary'}`, {
				'argo-button--full-width': block,
			})}>
			{children}
			{isLoading && (
				<div className="spinner-border-sm">
					<Loader width={20} height={20} color="#ffffff" />
				</div>
			)}
		</button>
	);
};

import 'components/atoms/loader/styles.scss';

import classNames from 'classnames';
import { Loader } from 'components/atoms/loader/Loader';
import React, { FC, ReactNode } from 'react';

interface Props {
	children?: string | ReactNode | ReactNode[];
	loading?: boolean;
	fullScreen?: boolean;
}

export const ZLoaderCover: FC<Props> = ({ children, loading, fullScreen }: Props) => (
	<div className={classNames('zlifecycle-spinner-cover', { 'zlifecycle-spinner-cover__full-screen': fullScreen })}>
		{loading && (
			<div className="zlifecycle-spinner-cover__loader d-flex align-items-center justify-content-center bg-white opacity-50">
				<Loader />
			</div>
		)}
		{children}
	</div>
);

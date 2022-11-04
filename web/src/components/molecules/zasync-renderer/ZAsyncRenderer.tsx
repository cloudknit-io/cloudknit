import { Loader } from 'components/atoms/loader/Loader';
import React, { FC, useEffect, useState } from 'react';

type Props = {
	promise: Promise<any>;
	Component: React.FC<any>;
	hideLoader?: boolean;
	defaultValue?: any;
	componentProps?: any;
};

export const ZAsyncRenderer: FC<Props> = ({ promise, hideLoader, defaultValue, Component, componentProps }) => {
	const [data, setData] = useState<any>(null);

	useEffect(() => {
		setData(defaultValue);
	}, [defaultValue]);
	const nullOrUndefined = (obj: any) => obj === null || obj === undefined;
	useEffect(() => {
		promise
			.then(data => {
				if (!nullOrUndefined(data)) setData(data || '');
			})
			.catch((err: Error) => {
				setData(err.message);
			});
	}, [promise]);

	return (
		<>
			{hideLoader ? (
				<Component data={data} {...componentProps} />
			) : !nullOrUndefined(data) ? (
				<Component data={data} {...componentProps} />
			) : (
				<Loader height={16} width={16} />
			)}
		</>
	);
};

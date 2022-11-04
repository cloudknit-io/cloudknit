import { Loader } from 'components/atoms/loader/Loader';
import React, { FC, useEffect, useState } from 'react';
import { Subject } from 'rxjs';

type Props = {
	subject: Subject<any>;
	Component: React.FC<any>;
	hideLoader?: boolean;
	defaultValue?: any;
	componentProps?: any;
};

export const ZStreamRenderer: FC<Props> = ({ subject, hideLoader, defaultValue, Component, componentProps }) => {
	const [data, setData] = useState<any>(defaultValue);
	const nullOrUndefined = (obj: any) => obj === null || obj === undefined;
	useEffect(() => {
		const $ssub = subject.subscribe((data: any) => {
			if (!nullOrUndefined(data)) setData(data || '');
		});

		return () => $ssub.unsubscribe();
	}, [subject]);

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

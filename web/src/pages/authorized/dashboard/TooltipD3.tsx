import React, { FC } from 'react';

interface TooltipD3Data {
	data: any;
}

export const TooltipD3: FC<TooltipD3Data> = (props: TooltipD3Data) => {
	const { data } = props;
	return data == null ? (
		<></>
	) : (
		<div className={data.classNames} style={{ top: data.top, left: data.left }}>
			{data.card}
		</div>
	);
};

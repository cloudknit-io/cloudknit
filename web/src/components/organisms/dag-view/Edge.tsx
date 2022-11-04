import React, { FC } from 'react';

type Props = {
	linkData: any;
};

export const Edge: FC<Props> = ({ linkData }) => {
	const drawPath = (linkData: any) => {
		const { target, source } = linkData;
		const deltaY = target.y - source.y;
		return `M${source.x},${source.y + 60} V${source.y + deltaY / 2} H${target.x} V${target.y - 30}`;
	};

	return <path d={drawPath(linkData)} style={{ fill: 'none', stroke: '#BBB' }} />;
};

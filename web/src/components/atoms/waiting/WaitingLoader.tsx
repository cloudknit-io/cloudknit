import React, { FC } from 'react';
import './style.scss';

interface Props {
	color?: string;
	radius?: number;
	title?: string;
}

export const WaitingLoader: FC<Props> = ({ color, radius, title }: Props) => {
	radius = radius ?? 50;
	const height = radius * 2 + 40;
	const gap = (radius * 2) + 5;
	const width = (gap * 3) + 5;

	return (
		<svg width={width} height={height} className='waiting-loader' preserveAspectRatio="xMidYMid">
			{title ? <title>{title}</title> : null}
			<circle
				cx={radius + 5}
				cy={height / 2}
				fill={color || '#252625'}
				r={radius}>
			</circle>
			<circle
				cx={radius + gap + 5}
				cy={height / 2}
				fill={color || '#252625'}
				r={radius}>
			</circle>
			<circle
				cx={(gap * 2) + radius + 5}
				cy={height / 2}
				fill={color || '#252625'}
				r={radius}>
			</circle>
		</svg>
	);
};

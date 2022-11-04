import React, { FC } from 'react';

interface Props {
	color?: string;
	height?: number;
	width?: number;
	title?: string;
}

export const Loader: FC<Props> = ({ color, height, width, title }: Props) => (
	<svg width={width || 60} height={height || 60} viewBox="0 0 100 100" preserveAspectRatio="xMidYMid">
		{title ? <title>{title}</title> : null}
		<circle
			cx="50"
			cy="50"
			fill="none"
			stroke={color || '#252625'}
			strokeWidth="10"
			r="35"
			strokeDasharray="164.93361431346415 56.97787143782138"
			strokeLinecap="round">
			<animateTransform
				attributeName="transform"
				type="rotate"
				repeatCount="indefinite"
				dur="1s"
				values="0 50 50;360 50 50"
				keyTimes="0;1"
			/>
		</circle>
	</svg>
);

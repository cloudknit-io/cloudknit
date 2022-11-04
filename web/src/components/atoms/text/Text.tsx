import classNames from 'classnames';
import React, { Component, FC, ReactNode } from 'react';

type TextTypes = 'p' | 'h1' | 'h2' | 'h3' | 'h4' | 'h5' | 'h6';
type TextSizes = '12' | '14' | '16' | '18' | '20' | '24' | '28' | '36';
type TextColors = 'white' | 'default';
type TextWeight = 'regular' | 'bold' | 'extra-bold' | 'thin' | 'extra-thin';
type TextSpacing = '0' | '2' | '4';
type FontFamily = 'dm' | 'default';
type TextLineHeight = '16' | '18' | '24' | '26' | '32' | '36' | '46';

interface Props {
	type?: TextTypes;
	size?: TextSizes;
	color?: TextColors;
	family?: FontFamily;
	weight?: TextWeight;
	spacing?: TextSpacing;
	lineHeight?: TextLineHeight;
	upperCase?: boolean;
	className?: string;
	children: string | ReactNode;
}
const common = '';
const generatePreset = (
	props: Props,
	size?: TextSizes,
	color?: TextColors,
	family?: FontFamily,
	weight?: TextWeight,
	spacing?: TextSpacing,
	lineHeight?: TextLineHeight
): string => {
	return classNames(
		common,
		`size-${size || '2xl'}`,
		`color-${color || 'default'}`,
		`weight-${weight || 'normal'}`,
		`family-${family || 'default'}`,
		`spacing-${spacing || 'normal'}`,
		`leading-${lineHeight || 'normal'}`,
		{ 'text-uppercase': props.upperCase },
		props.className
	);
};

const Headline: FC<Props> = (props: Props) => {
	const h2Preset = generatePreset(
		props,
		props.size || '24',
		props.color || 'default',
		props.family || 'default',
		props.weight || 'extra-bold',
		props.spacing || '0',
		'36'
	);
	return <h2 className={classNames(h2Preset, 'headline')}>{props.children}</h2>;
};

const Body: FC<Props> = (props: Props) => {
	const h4Preset = generatePreset(
		props,
		props.size || '24',
		props.color || 'default',
		props.family || 'default',
		props.weight || 'regular',
		props.spacing,
		props.lineHeight || '36'
	);
	return <p className={h4Preset}>{props.children}</p>;
};

class ZText extends Component<Props, {}> {
	public static Headline = Headline;
	public static Body = Body;

	render(): ReactNode {
		const {
			size = 'base',
			weight = 'normal',
			spacing = 'normal',
			lineHeight = 'normal',
			upperCase = false,
			className,
			children,
		} = this.props;
		const preset = classNames(
			common,
			`size-${size}`,
			`weight-${weight}`,
			`spacing-${spacing}`,
			`leading-${lineHeight}`,
			{ 'text-uppercase': upperCase },
			className
		);
		return <span className={preset}>{children}</span>;
	}
}
export { Headline, Body, ZText };

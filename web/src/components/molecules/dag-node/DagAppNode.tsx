import React, { FC } from 'react';
import ReactDOMServer from 'react-dom/server';
import { getClassName, getTextWidth, getTime } from 'components/organisms/tree-view/node-figure-helper';
import { CostRenderer, getSyncStatusIcon } from '../cards/renderFunctions';
import { ZSyncStatus } from 'models/argo.models';

type Props = {
	id: string;
	componentStatus: string;
	name: string;
	syncFinishedAt: string;
	isSkipped: boolean;
	displayValue: string;
	HealthIcon: JSX.Element;
	SyncStatus: ZSyncStatus;
	Icon: JSX.Element;
	ExpandIcon: JSX.Element;
	onNodeClick: (...params: any) => any
};

export const ZDagAppNode: FC<Props> = ({
	id,
	componentStatus,
	name,
	syncFinishedAt,
	isSkipped,
	displayValue,
	HealthIcon,
	Icon,
	ExpandIcon,
	SyncStatus,
	onNodeClick
}: Props) => {
	let width = getTextWidth(name);
	const timeWidth = 40 + getTextWidth(getTime(syncFinishedAt));
	width = timeWidth > width ? timeWidth : width;

	return (
		<g onClick={(e) => {
			onNodeClick(id);
		}}>
			<rect
				height={70}
				width={width + 100 > 160 ? width + 100 : 160}
				fill="#fff"
				rx="10"
				className={`node node__pod ${id === 'root' ? 'root' : 'node__pod'}${getClassName(
					componentStatus || ''
				)} ${isSkipped ? ' striped' : ''}`}
			/>
			<text x="65" y="25" fill="#323232" fontFamily="DM Sans" fontWeight={'light'} fontSize="15px">
				{name}
			</text>
			{syncFinishedAt && (
				<text
					x={HealthIcon && ExpandIcon ? 125 : 105}
					y="51"
					fill="#323232"
					fontFamily="DM Sans"
					fontWeight={'light'}
					fontSize="14px">
					{' | ' + getTime(syncFinishedAt)}
				</text>
			)}
			{SyncStatus && (
				<g
					transform={`translate(${65},${38})`}>
						{getSyncStatusIcon(SyncStatus)}
					</g>
			)}
			{HealthIcon && (
				<g
					transform={`translate(${85},${38})`}
					dangerouslySetInnerHTML={{
						__html: ReactDOMServer.renderToString(HealthIcon),
					}}></g>
			)}
			{ExpandIcon && (
				<g
					transform={`translate(${HealthIcon ? 105 : 85},${38})`}
					dangerouslySetInnerHTML={{
						__html: ReactDOMServer.renderToString(ExpandIcon),
					}}	
					onClick={e => {
						e.stopPropagation();
						ExpandIcon.props.onClick(e);
					}}>
						
					</g>
			)}
			{Icon && (
				<g
					transform={`translate(${12},${12}) scale(${id === 'root' ? 0.25 : 0.35})`}
					dangerouslySetInnerHTML={{
						__html: ReactDOMServer.renderToString(Icon),
					}}></g>
			)}
		</g>
	);
};

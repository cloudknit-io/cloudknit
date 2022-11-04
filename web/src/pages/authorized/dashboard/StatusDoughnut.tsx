import { ZDropdownMenu, ZDropdownMenuJSX } from 'components/molecules/dropdown-menu/DropdownMenu';
import * as d3 from 'd3';
import { ZSyncStatus } from 'models/argo.models';
import React, { FC, useEffect, useState } from 'react';

import { TooltipD3 } from './TooltipD3';

export interface StatusDoughnutData {
	data: any;
}

export const StatusDoughnut: FC<StatusDoughnutData> = (props: StatusDoughnutData) => {
	const { data } = props;
	const [tooltipData, setTooltipData] = useState<any>(null);
	const containerDiv = React.useRef<HTMLDivElement>(null);
	const { height, width } = {
		height: 640,
		width: 640,
	};
	const dim = { height, width };

	const initializeDoughnut = (data: any) => {
		// const color: any = d3
		// 	.scaleOrdinal()
		// 	.domain(data.map((d: any) => d.name))
		// 	.range(d3.quantize(t => d3.interpolateSpectral(t * 0.8 + 0.1), data.length).reverse());
		const color: any = (d: any) => {
			// if (d.data.labels.type === 'project') {
			// 	return '#2ebc38';
			// } else if (d.data.labels.type === 'config') {
			// 	return '#fcb00f';
			// } else if (d.data.labels.type === 'environment') {
			// 	return '#463af6';
			// } else {
			// 	return 'gray';
			// }
			if (d.name === 'InSync' || d.name === 'Healthy') {
				return '#2ebc38';
			} else if (d.name === 'Provisioned') {
				return 'yellowgreen';
			} else if (d.name === 'Unknown') {
				return 'gray';
			} else {
				return 'gray';
			}
		};

		const radius = Math.min(dim.width, dim.height) / 2;

		const arc: any = d3
			.arc()
			.innerRadius(radius * 0.67)
			.outerRadius(radius - 1);

		const pie = d3
			.pie()
			.padAngle(0.005)
			.sort(null)
			.value((d: any) => d.value);
		const arcs = pie(data);

		const svg = d3
			.create('svg')
			.attr('viewBox', `${-dim.width / 2}, ${-dim.height / 2}, ${dim.width}, ${dim.height}`);

		svg.selectAll('path')
			.data(arcs)
			.join('path')
			.attr('fill', (d: any) => color(d.data))
			.attr('d', arc)
			.on('mousemove', (e, d: any) => {
				// console.log(d);
				setTooltipData({
					card: (
						<ZDropdownMenu
							isOpened={true}
							items={d.data.components.map((c: any) => ({
								text: c.componentName,
								action: () => {},
							}))}
						/>
					),
					top: e.layerY,
					left: e.layerX + 60,
					classNames: 'tooltip-d3',
				});
				(containerDiv.current?.parentElement?.style as any).zIndex = 99;
			})
			.append('title')
			.attr('style', (d: any) => `display : ${d.data.value === 0 ? 'none' : 'inherit'}`)
			.text((d: any) => `${d.data.name}: ${d.data.value.toLocaleString()}`);

		svg.append('g')
			.attr('font-family', 'sans-serif')
			.attr('font-size', 12)
			.attr('text-anchor', 'middle')
			.selectAll('text')
			.data(arcs)
			.join('text')
			.attr('transform', (d: any) => `translate(${arc.centroid(d)})`)
			.attr('style', (d: any) => `display : ${d.data.value === 0 ? 'none' : 'inherit'}`)
			.call(text =>
				text
					.append('tspan')
					.attr('y', '-0.4em')
					.attr('font-weight', 'bold')
					.text((d: any) => d.data.name)
			)
			.call(text =>
				text
					.filter(d => d.endAngle - d.startAngle > 0.25)
					.append('tspan')
					.attr('x', 0)
					.attr('y', '0.7em')
					.attr('style', (d: any) => `display : ${d.data.value === 0 ? 'none' : 'inherit'}`)
					.attr('fill-opacity', 0.7)
					.text((d: any) => d.data.value.toLocaleString())
			);

		return svg.node();
	};

	useEffect(() => {
		if (data) {
			const doughnut = initializeDoughnut(data);
			const container: HTMLDivElement = containerDiv.current as HTMLDivElement;
			if (container && doughnut) {
				const svg = container.querySelector('svg');
				svg && container.removeChild(svg);
				container.appendChild(doughnut);
			}
		}
	}, [data]);

	return (
		<div
			ref={containerDiv}
			onMouseLeave={() => {
				(containerDiv.current?.parentElement?.style as any).zIndex = 0;
				setTooltipData({
					...tooltipData,
					classNames: 'com-cards tooltip-d3 teams hide',
				});
			}}>
			<TooltipD3 data={tooltipData} />
		</div>
	);
};

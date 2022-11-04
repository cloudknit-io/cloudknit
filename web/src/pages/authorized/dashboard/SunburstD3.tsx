import * as d3 from 'd3';
import { ESyncStatus, ZSyncStatus } from 'models/argo.models';
import React, { FC, useEffect } from 'react';

import { getClassName } from './helpers';

interface SunburstData {
	data: any;
}

export const SunburstD3: FC<SunburstData> = (props: SunburstData) => {
	const { data } = props;
	const containerDiv = React.useRef<HTMLDivElement>(null);
	const { height, width } = { height: 640, width: 640 };

	const radius = width / 6;
	const format = d3.format(',d');

	const color = (d: any) => {
		if (d.data.labels.type === 'project') {
			return '#2ebc38';
		} else if (d.data.labels.type === 'config') {
			return d.data.componentStatus === 'in_sync' ? 'green' : 'red';
		} else if (d.data.labels.type === 'environment') {
			return '#463af6';
		} else {
			return 'gray';
		}
	};

	const arc = d3
		.arc()
		.startAngle((d: any) => d.x0)
		.endAngle((d: any) => d.x1)
		.padAngle((d: any) => Math.min((d.x1 - d.x0) / 2, 0.005))
		.padRadius(radius * 1.5)
		.innerRadius((d: any) => d.y0 * radius)
		.outerRadius((d: any) => Math.max(d.y0 * radius, d.y1 * radius - 1));

	const partition = (data: any) => {
		const root = d3
			.hierarchy(data)
			.sum(d => d.value)
			.sort((a: any, b: any) => b.value - a.value);
		return d3.partition().size([2 * Math.PI, root.height + 1])(root);
	};

	const initializeSunburst = (data: any) => {
		const root = partition(data);

		root.each((d: any) => (d.current = d));

		const svg = d3.create('svg').attr('viewBox', `0, 0, ${width}, ${height}`).style('font', '10px sans-serif');

		const g = svg.append('g').attr('transform', `translate(${width / 2},${height / 2})`);

		const path = g
			.append('g')
			.selectAll('path')
			.data(root.descendants().slice(1))
			.join('path')
			.attr('class', (d: any) => {
				// while (d.depth > 1) d = d.parent;
				return 'wedge' + getClassName(d.data.componentStatus || d.data.syncStatus);
			})
			.attr('fill-opacity', (d: any) => (arcVisible(d.current) ? (d.children ? 0.6 : 0.4) : 0))
			.attr('d', (d: any) => arc(d.current));

		path.filter((d: any) => d.children)
			.style('cursor', 'pointer')
			.on('click', clicked);

		path.append('title').text((d: any) => d.data.componentStatus || d.data.syncStatus);

		const label = g
			.append('g')
			.attr('pointer-events', 'none')
			.attr('text-anchor', 'middle')
			.style('user-select', 'none')
			.selectAll('text')
			.data(root.descendants().slice(1))
			.join('text')
			.attr('dy', '0.35em')
			.attr('fill-opacity', (d: any) => +labelVisible(d.current))
			.attr('transform', (d: any) => labelTransform(d.current))
			.text((d: any) => d.data.componentName || d.data.displayValue);

		const parent = g
			.append('circle')
			.datum(root)
			.attr('r', radius)
			.attr('fill', 'none')
			.attr('pointer-events', 'all')
			.on('click', clicked);

		function clicked(event: any, p: any) {
			parent.datum(p.parent || root);

			root.each(
				(d: any) =>
					(d.target = {
						x0: Math.max(0, Math.min(1, (d.x0 - p.x0) / (p.x1 - p.x0))) * 2 * Math.PI,
						x1: Math.max(0, Math.min(1, (d.x1 - p.x0) / (p.x1 - p.x0))) * 2 * Math.PI,
						y0: Math.max(0, d.y0 - p.depth),
						y1: Math.max(0, d.y1 - p.depth),
					})
			);

			const t: any = g.transition().duration(750);

			path.transition(t)
				.tween('data', (d: any) => {
					const i = d3.interpolate(d.current, d.target);
					return (t: any) => (d.current = i(t));
				})
				.filter(function (d: any): any {
					const t: any = this;
					return +t.getAttribute('fill-opacity') || arcVisible(d.target);
				})
				.attr('fill-opacity', (d: any) => (arcVisible(d.target) ? (d.children ? 0.6 : 0.4) : 0))
				.attrTween('d', (d: any) => (): any => arc(d.current));

			label
				.filter(function (d: any): any {
					const t: any = this;
					return +t.getAttribute('fill-opacity') || labelVisible(d.target);
				})
				.transition(t)
				.attr('fill-opacity', (d: any) => +labelVisible(d.target))
				.attrTween('transform', (d: any) => () => labelTransform(d.current));
		}

		function arcVisible(d: any) {
			return d.y1 <= 3 && d.y0 >= 1 && d.x1 > d.x0;
		}

		function labelVisible(d: any) {
			return d.y1 <= 3 && d.y0 >= 1 && (d.y1 - d.y0) * (d.x1 - d.x0) > 0.03;
		}

		function labelTransform(d: any) {
			const x = (((d.x0 + d.x1) / 2) * 180) / Math.PI;
			const y = ((d.y0 + d.y1) / 2) * radius;
			return `rotate(${x - 90}) translate(${y},0) rotate(${x < 180 ? 0 : 180})`;
		}

		return svg.node();
	};

	useEffect(() => {
		if (data) {
			const sunburst = initializeSunburst({
				name: 'root',
				children: data,
			});
			const container: HTMLDivElement = containerDiv.current as HTMLDivElement;
			if (container && sunburst) {
				container.firstChild && container.removeChild(container.firstChild);
				container.appendChild(sunburst);
			}
		}
	}, [data]);

	return <div ref={containerDiv}></div>;
};

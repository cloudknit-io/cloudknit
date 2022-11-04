import * as d3 from 'd3';
import React, { FC, useEffect } from 'react';

interface TagBarchartData {
	data: any;
}

export const TagBarchartD3: FC<TagBarchartData> = (props: TagBarchartData) => {
	const { data } = props;
	const containerDiv = React.useRef<HTMLDivElement>(null);
	const { height, width } = {
		height: 640,
		width: 640,
	};
	const margin = { top: 30, right: 0, bottom: 30, left: 30 };

	const initializeBarchart = (data: any) => {
		const yAxis = (g: any) =>
			g
				.attr('transform', `translate(${margin.left},0)`)
				.call(d3.axisLeft(y).ticks(null, data.format))
				.call((g: any) => g.select('.domain').remove())
				.call((g: any) =>
					g
						.append('text')
						.attr('x', -margin.left)
						.attr('y', 10)
						.attr('fill', 'currentColor')
						.attr('text-anchor', 'start')
						.text(data.y)
				);
		const xAxis = (g: any) =>
			g.attr('transform', `translate(0,${height - margin.bottom})`).call(
				d3
					.axisBottom(x)
					.tickFormat((i: any) => data[i].name)
					.tickSizeOuter(0)
			);
		const x: any = d3
			.scaleBand()
			.domain(d3.range(data.length).map(e => e.toString()))
			.range([margin.left, width - margin.right])
			.padding(0.1);

		const y: any = d3
			.scaleLinear()
			.domain([0, Number(d3.max(data, (d: any) => d.value))])
			.nice()
			.range([height - margin.bottom, margin.top]);

		// const format = x.tickFormat(20, data.format);

		const svg = d3.create('svg').attr('viewBox', `0, 0, ${width}, ${height}`);

		svg.append('g')
			.selectAll('rect')
			.data(data)
			.join('rect')
			.attr('x', (d, i) => x(i))
			.attr('y', (d: any) => y(d.value))
			.attr('height', (d: any) => y(0) - y(d.value))
			.attr('width', x.bandwidth())
			.attr('fill', (d: any) => {
				return d.color;
			})
			.on('mouseover', (event: any, data: any) => {
				console.log(data.components.map((e: any) => e.componentName));
			});

		svg.append('g').call(xAxis);

		svg.append('g').call(yAxis);

		return svg.node();
	};

	useEffect(() => {
		if (data) {
			const colors = ['#2ebc38', '#fcb00f', '#463af6', 'gray', 'teal'];
			const getTag = () => ['ML', 'AI', 'DB', 'WS', 'API'][Math.floor(Math.random() * 5)];
			const map: any = {};
			const len = data.length;
			data.forEach((e: any) => {
				const name = getTag();
				if (map[name]) {
					map[name].components.push(e);
				} else {
					map[name] = {};
					map[name].components = [e];
				}
			});

			const bd = Object.keys(map).map((k: any) => ({
				name: k,
				components: map[k].components,
				value: map[k].components.length / len,
				color: colors.pop(),
			}));

			const barchart = initializeBarchart(bd);
			const container: HTMLDivElement = containerDiv.current as HTMLDivElement;
			if (container && barchart) {
				container.firstChild && container.removeChild(container.firstChild);
				container.appendChild(barchart);
			}
		}
	}, [data]);

	return <div ref={containerDiv}></div>;
};

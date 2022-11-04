import * as d3 from 'd3';
import React, { FC, useEffect } from 'react';

interface CloudBarchartData {
	data: any;
}

export const CloudBarchartD3: FC<CloudBarchartData> = (props: CloudBarchartData) => {
	const { data } = props;
	const containerDiv = React.useRef<HTMLDivElement>(null);
	const { height, width } = {
		height: 640,
		width: 640,
	};
	const margin = { top: 30, right: 0, bottom: 30, left: 50 };

	const initializeBarchart = (data: any) => {
		const yAxis = (g: any) =>
			g.attr('transform', `translate(${margin.left},0)`).call(
				d3
					.axisLeft(y)
					.tickFormat((i: any) => data[i].name)
					.tickSizeOuter(0)
			);
		const xAxis = (g: any) =>
			g
				.attr('transform', `translate(0,${margin.top})`)
				.call(d3.axisTop(x).ticks(width / 80, data.format))
				.call((g: any) => g.select('.domain').remove());

		const x: any = d3
			.scaleLinear()
			.domain([0, Number(d3.max(data, (d: any) => d.value))])
			.range([margin.left, width - margin.right]);

		const y: any = d3
			.scaleBand()
			.domain(d3.range(data.length) as Iterable<string>)
			.rangeRound([margin.top, height - margin.bottom])
			.padding(0.1);

		const format = x.tickFormat(20, data.format);

		const svg = d3.create('svg').attr('viewBox', `0, 0, ${width}, ${height}`);

		svg.append('g')
			.selectAll('rect')
			.data(data)
			.join('rect')
			.attr('x', x(0))
			.attr('y', (d: any, i: any) => y(i))
			.attr('width', (d: any) => x(d.value) - x(0))
			.attr('height', y.bandwidth())
			.attr('fill', (d: any) => {
				return d.color;
			})
			.on('mouseover', (event: any, data: any) => {
				console.log(data.components.map((e: any) => e.componentName));
			});

		svg.append('g')
			.attr('fill', 'white')
			.attr('text-anchor', 'end')
			.attr('font-family', 'sans-serif')
			.attr('font-size', 12)
			.selectAll('text')
			.data(data)
			.join('text')
			.attr('x', (d: any) => x(d.value))
			.attr('y', (d: any, i: any) => y(i) + y.bandwidth() / 2)
			.attr('dy', '0.35em')
			.attr('dx', -4)
			.text((d: any) => `${Number(d.value * 100).toFixed(2)}%`)
			.call(text =>
				text
					.filter((d: any) => x(d.value) - x(0) < 20) // short bars
					.attr('dx', +4)
					.attr('fill', 'black')
					.attr('text-anchor', 'start')
			);

		svg.append('g').call(xAxis);

		svg.append('g').call(yAxis);

		return svg.node();
	};

	useEffect(() => {
		if (data) {
			const colors = ['#2ebc38', '#fcb00f', '#463af6', 'gray', 'teal'];
			// if (d.data.labels.type === 'project') {
			// 	return '#2ebc38';
			// } else if (d.data.labels.type === 'config') {
			// 	return '#fcb00f';
			// } else if (d.data.labels.type === 'environment') {
			// 	return '#463af6';
			const getTag = () => ['AWS', 'GCP', 'AZURE', 'On-Prem'][Math.floor(Math.random() * 4)];
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

import * as d3 from 'd3';
import React, { FC, useEffect } from 'react';

interface HistoryCalenderData {
	data: any;
}

export const HistoryCalender: FC<HistoryCalenderData> = (props: HistoryCalenderData) => {
	const containerDiv = React.useRef<HTMLDivElement>(null);
	const weekday = ['weekday', 'sunday'][0];
	const cellSize = 17;
	const formatValue = (v: any) => {
		return `Deployed ${v === 1 ? v + ' time' : v + ' times'} on this day`;
	};
	const formatClose = (v: any) => {
		return `Costing : $${Number(v).toFixed(3)}`;
	};
	const formatDate = d3.utcFormat('%x');
	const formatDay = (i: any) => 'SMTWTFS'[i];
	const formatMonth = d3.utcFormat('%b');
	const timeWeek = weekday === 'sunday' ? d3.utcSunday : d3.utcMonday;
	const countDay = weekday === 'sunday' ? (i: any) => i : (i: any) => (i + 6) % 7;
	const pathMonth = (t: any) => {
		const n = weekday === 'weekday' ? 5 : 7;
		const d = Math.max(0, Math.min(n, countDay(t.getUTCDay())));
		const w = timeWeek.count(d3.utcYear(t), t);
		return `${
			d === 0
				? `M${w * cellSize},0`
				: d === n
				? `M${(w + 1) * cellSize},0`
				: `M${(w + 1) * cellSize},0V${d * cellSize}H${w * cellSize}`
		}V${n * cellSize}`;
	};

	const { height, width } = { height: cellSize * (weekday === 'weekday' ? 7 : 9), width: 640 };
	const color = (data: any): any => {
		if (data < 3) {
			return 'rgba(0,0,255, 0.2)';
		} else if (data <= 6) {
			return 'rgba(0,0,255, 0.4)';
		} else if (data <= 9) {
			return 'rgba(0,0,255, 0.6)';
		} else if (data > 9) {
			return 'rgba(0,0,255, 0.8)';
		}
		return 'rgba(0,0,0,0.1)';
	};
	const initializeHistoryCalender = (data: any) => {
		const years = d3.groups(data, (d: any) => d.date.getUTCFullYear()).reverse();
		const svg = d3
			.create('svg')
			.attr('viewBox', `0, 0, ${width}, ${height * years.length}`)
			.attr('font-family', 'sans-serif')
			.attr('font-size', 10);

		const year = svg
			.selectAll('g')
			.data(years)
			.join('g')
			.attr('transform', (d, i) => `translate(40.5,${height * i + cellSize * 1.5})`);

		year.append('text')
			.attr('x', -5)
			.attr('y', -5)
			.attr('font-weight', 'bold')
			.attr('text-anchor', 'end')
			.text(([key]) => key);

		year.append('g')
			.attr('text-anchor', 'end')
			.selectAll('text')
			.data(weekday === 'weekday' ? d3.range(1, 6) : d3.range(7))
			.join('text')
			.attr('x', -5)
			.attr('y', i => (countDay(i) + 0.5) * cellSize)
			.attr('dy', '0.31em')
			.text(formatDay);

		year.append('g')
			.selectAll('rect')
			.data(
				weekday === 'weekday'
					? ([, values]) => values.filter(d => ![0, 6].includes(d.date.getUTCDay()))
					: ([, values]) => values
			)
			.join('rect')
			.attr('width', cellSize - 1)
			.attr('height', cellSize - 1)
			.attr('x', d => timeWeek.count(d3.utcYear(d.date), d.date) * cellSize + 0.5)
			.attr('y', d => countDay(d.date.getUTCDay()) * cellSize + 0.5)
			.attr('fill', (d: any) => color(d.value))
			.append('title')
			.text(
				d => `${formatDate(d.date)}
          ${formatValue(d.value)}${
					d.close === undefined
						? ''
						: `
          ${formatClose(d.close)}`
				}`
			);

		const month = year
			.append('g')
			.selectAll('g')
			.data(([, values]) => d3.utcMonths(d3.utcMonth(values[0].date), values[values.length - 1].date))
			.join('g');

		month
			.filter((d, i: any) => i)
			.append('path')
			.attr('fill', 'none')
			.attr('stroke', '#fff')
			.attr('stroke-width', 3)
			.attr('d', pathMonth);

		month
			.append('text')
			.attr('x', d => timeWeek.count(d3.utcYear(d), timeWeek.ceil(d)) * cellSize + 2)
			.attr('y', -5)
			.text(formatMonth);

		return svg.node();
	};

	useEffect(() => {
		const data = props.data.flat().filter((e: any) => e);
		if (data?.length > 0) {
			const map: any = {};
			for (let i = 0; i < data.length; i++) {
				const e = new Date(data[i].deployStartedAt).toISOString().split('T')[0];
				if (!map[e]) {
					map[e] = {};
				}
				map[e] = {
					date: new Date(data[i].deployStartedAt),
					value: map[e].value ? map[e].value + 1 : 1,
					close: Math.random() * 1000,
				};
			}

			const calender = initializeHistoryCalender(
				Object.keys(map).map(k => ({
					date: map[k].date,
					value: map[k].value,
					close: map[k].close,
				}))
			);
			const container: HTMLDivElement = containerDiv.current as HTMLDivElement;
			if (container && calender) {
				container.firstChild && container.removeChild(container.firstChild);
				container.appendChild(calender);
			}
		}
	}, [props.data]);
	return <div ref={containerDiv}></div>;
};

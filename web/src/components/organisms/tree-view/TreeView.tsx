import './styles.scss';

import { ReactComponent as ArrowUp } from 'assets/images/icons/chevron-right.svg';
import * as d3 from 'd3';
import dagreD3, { GraphLabel } from 'dagre-d3';
import { renderSyncStatus } from 'pages/authorized/environment-components/helpers';
import React, { FC, useEffect, useRef, useState } from 'react';
import ReactDOMServer from 'react-dom/server';

import { getCurveType } from './curve-helper';
import { getShapeLabel, updateShapeLabel } from './shape-helper';
// @ts-ignore
import css_vars from './styles.scss';
import { cleanDagNodeCache } from './node-figure-helper';

interface Props {
	arrowType: any;
	data: any;
	nodeSep: any;
	rankSep: any;
	ranker: any;
	onNodeClick: any;
	zoomEventHandlers: any;
	environmentId: string;
	rankDir: 'TB' | 'LR' | '';
	deploymentType: 'blue' | 'green' | 'grey' | '';
}

const colorLegend = [
	{
		key: 'Succeeded',
		value: css_vars.successful,
		order: 5,
	},
	{
		key: 'Failed',
		value: css_vars.failed,
		order: 6,
	},
	{
		key: 'Waiting for approval',
		value: css_vars.pending,
		order: 2,
	},
	{
		key: 'In Progress',
		value: css_vars.initializing,
		order: 0,
	},
	{
		key: 'Destroyed/Not Provisioned',
		value: css_vars.destroyed,
		order: 5,
	},
];

const Tree: FC<Props> = ({
	environmentId,
	arrowType,
	data,
	nodeSep,
	rankSep,
	ranker,
	rankDir,
	onNodeClick,
	zoomEventHandlers,
	deploymentType,
}: Props) => {
	const containerRef = useRef<HTMLDivElement>(null);
	const [dagGraph, setDagGraph] = useState<any>();
	const [isPositionSet, setPosition] = useState<boolean>(false);
	const [legendVisibility, toggleLegend] = useState<boolean>(true);

	const initializeDagGraph = () => {
		const g: any = new dagreD3.graphlib.Graph({ directed: true });
		// Set an object for the graph label
		g.setGraph({});
		g.graph().rankdir = rankDir || 'TB';
		g.graph().ranker = ranker;
		g.graph().ranksep = rankSep;
		g.graph().nodesep = nodeSep;
		// Default to assigning a new object as a label for each new edge.
		g.setDefaultEdgeLabel(function () {
			return {};
		});

		const svg = d3.select(`#deployment-type-${deploymentType} .tree-node`);
		svg.attr('height', '80vh');
		setDagGraph(g);
	};

	const removeOldNodes = (g: any, data: any) => {
		if (g.nodes().length === 0) {
			cleanDagNodeCache('temp');
		}
		const activeNodes = data.map((e: any) => e.id);
		const oldNodes = g.nodes().filter((e: string) => !activeNodes.includes(e));
		oldNodes.forEach((e: string) => g.removeNode(e));
	};

	const setNewNode = (item: any, g: any, curve: any) => {
		const label = getShapeLabel({ ...item });

		g.setNode(item.id, label);

		if (item.dependsOn) {
			item.dependsOn.forEach((dependency: any) => {
				g.setEdge(dependency, item.id, {
					curve,
					style: 'stroke: blue; fill:none; stroke-width: 1px; stroke-dasharray: 5, 5;',
					arrowheadStyle: 'fill: blue',
				});
			});
		}
	};

	useEffect(() => {
		if (!dagGraph || environmentId) {
			initializeDagGraph();
		}
	}, [environmentId]);

	useEffect(() => {
		if (!dagGraph || data.length === 0) {
			return;
		}
		const g = dagGraph;
		const curve = getCurveType(arrowType);
		const svg = d3.select(`#deployment-type-${deploymentType} .tree-node`);
		const render = new dagreD3.render();
		const inner: any = svg.select('g');
		removeOldNodes(g, data);
		data.forEach((item: any) => {
			const exists = updateShapeLabel({ ...item });
			if (!exists) {
				setNewNode(item, g, curve);
			}
		});

		inner.on('mouseover', (event: any) => {
			const containerDiv = (event.target as Element)?.closest('.roundedCorners');
			if (containerDiv) {
				if (containerDiv.id === 'root') {
					return;
				}
				const { top, right } = { right: event.layerX, top: event.layerY };
				const tooltip = containerRef.current?.querySelector('.node_tooltip') as HTMLElement;
				if (tooltip) {
					tooltip.innerHTML = '';
					const id = containerDiv.id;
					const renderData = data.find((e: any) => e.name === id);
					if (renderData) {
						tooltip.innerHTML = ReactDOMServer.renderToString(renderSyncStatus(renderData));
					} else {
						return;
					}
					tooltip.style.top = top + 20 + 'px';
					tooltip.style.left = right + 20 + 'px';
					tooltip.style.transform = 'scale(1)';
					tooltip.style.opacity = '1';
				}
			}
		});

		inner.on('mouseout', (event: MouseEvent) => {
			const tooltip = containerRef.current?.querySelector('.node_tooltip') as HTMLElement;
			if (tooltip) {
				tooltip.style.transform = 'scale(0)';
			}
		});

		render(inner, g);

		if (!isPositionSet) {
			const zoom: any = d3.zoom().on('zoom', function (e) {
				inner.attr('transform', e.transform);
			});
			setImmediate(() => {
				const svgElem = svg.node() as Element;
				const gElem = svgElem.querySelector('#t-g');
				const svgWidth = svgElem.getBoundingClientRect().width || 2;
				const graphWidth = gElem?.getBoundingClientRect().width || 2;
				gElem?.classList.add('animate-children');
				svg.call(
					zoom.transform,
					d3.zoomIdentity.translate((Number(svgWidth || 0) - graphWidth) / 2, 20).scale(zoom.transform.k || 1)
				);
				gElem?.setAttribute('transform', `translate(${(svgWidth - graphWidth) / 2.5}, 0)`);
				setImmediate(() => {
					containerRef.current?.classList.add('show-graph');
					setPosition(true);
				});
			});
		}
	}, [arrowType, ranker, data, dagGraph]);

	useEffect(() => {
		if (dagGraph) {
			const svg = d3.select(`#deployment-type-${deploymentType} .tree-node`);
			const inner: any = svg.select('g');

			const zoom: any = d3.zoom().on('zoom', function (e) {
				if (e.transform.k <= 0.2) {
					e.transform.k = 0.2;
					return;
				}

				if (e.transform.k >= 2) {
					e.transform.k = 2;
					return;
				}
				inner.attr('transform', e.transform);
			});

			svg.call(zoom).node();
			svg.call(zoom).on('wheel.zoom', null).node();
			svg.call(zoom).on('dblclick.zoom', null).node();

			zoomEventHandlers(
				() => {
					zoom.scaleBy(svg.transition().duration(300), 1.2);
				},
				() => {
					zoom.scaleBy(svg.transition().duration(300), 0.8);
				},
				() => {
					// svg.call(zoom, d3.zoomIdentity);
				}
			);
		}
	}, [dagGraph, zoomEventHandlers]);

	useEffect(() => {
		if (dagGraph) {
			const svg = d3.select(`#deployment-type-${deploymentType} .tree-node`);
			const inner: any = svg.select('g');
			dagGraph.graph().ranker = ranker;
			dagGraph.graph().rankSep = rankSep;
			dagGraph.graph().nodesep = nodeSep;
			const render = new dagreD3.render();
			render(inner, dagGraph);
		}
	}, [dagGraph, nodeSep, rankSep, ranker]);

	useEffect(() => {
		if (dagGraph) {
			const curve = getCurveType(arrowType);
			dagGraph.edges().forEach(({v, w}: any) => {
				dagGraph.setEdge(v, w, {
					curve,
					style: 'stroke: blue; fill:none; stroke-width: 1px; stroke-dasharray: 5, 5;',
					arrowheadStyle: 'fill: blue',
				});
			});
			const svg = d3.select(`#deployment-type-${deploymentType} .tree-node`);
			const render = new dagreD3.render();
			const inner: any = svg.select('g');
			render(inner, dagGraph);
		}
	}, [dagGraph, arrowType]);

	return (
		<>
			<div
				id={`deployment-type-${deploymentType}`}
				className="tree-view-container flex items-center justify-center dag-graph"
				style={{ width: '100%' }}
				ref={containerRef}>
				{!deploymentType && (
					<div className={`modifier color-legend color-legend-${legendVisibility ? 'show' : 'hide'}`}>
						{/* <div className="color-legend_toggler">
							<button
								title="Color Legend"
								onClick={() => {
									toggleLegend(!legendVisibility);
								}}>
								<ArrowUp />
							</button>
						</div> */}
						<div className="color-legend_status">
							<div>
								<label>Status:</label>
								{colorLegend
									.sort((a, b) => a.order - b.order)
									.map(color => (
										<span className="color-legend_value" key={color.key}>
											<label style={{ background: color.value }}></label>
											<label>{color.key}</label>
										</span>
									))}
							</div>
						</div>
					</div>
				)}
				{!deploymentType && <div className="node_tooltip"></div>}
				<svg className="tree-node" height="100%" width="100%">
					<defs>
						<filter
							id="drop_shadow"
							filterUnits="objectBoundingBox"
							x="-50%"
							y="-50%"
							width="200%"
							height="200%">
							<feDropShadow dx="0" dy="0" stdDeviation="20" floodColor="#ccc" floodOpacity="1" />
						</filter>
						<pattern
							id="pattern-stripe"
							width="6"
							height="4"
							patternUnits="userSpaceOnUse"
							patternTransform="rotate(45)">
							<rect width="3" height="4" transform="translate(0,0)" fill="white"></rect>
						</pattern>
						<mask id="mask-stripe">
							<rect x="0" y="0" width="100%" height="100%" fill="url(#pattern-stripe)" />
						</mask>
					</defs>
					<g id="t-g" />
				</svg>
			</div>
		</>
	);
};
export default Tree;

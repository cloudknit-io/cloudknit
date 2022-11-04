import { getSVGNode } from 'components/organisms/tree-view/shape-helper';
import React from 'react';
import { ReactElement } from 'react';
import ReactDOMServer from 'react-dom/server';
import {
	PropertyBinding,
	YComponent,
	YEnvironmentBuilder,
	YMetaData,
	YOutput,
	YTag,
	YVariable,
	YVariableFile,
} from 'models/environment-builder';
import { ReactComponent as DeleteSvg } from '../assets/images/icons/card-status/sync/delete.svg';
import { ReactComponent as EditSvg } from '../assets/images/icons/edit.svg';
import { ReactComponent as Chevron } from '../assets/images/icons/chevron-right.svg';
import { uniqueId } from 'lodash';
import { getCurveType } from 'components/organisms/tree-view/curve-helper';
import classNames from 'classnames';

const getTextWidth = (name: string): number => {
	const canvas = document.createElement('canvas');
	const ctx: CanvasRenderingContext2D = canvas.getContext('2d') as CanvasRenderingContext2D;
	ctx.font = 'bold 15px "DM Sans"';
	let { width } = ctx.measureText(name);
	canvas.remove();
	return width;
};

const calcPoints = (e: any, graph: any, d3: any) => {
	const { v, w } = e;
	let edge = graph.edge(v, w),
		tail = graph.node(v),
		head = graph.node(w);
	let points = edge.points.slice(1, edge.points.length - 1);
	points = [{ x: tail.x, y: tail.y }];
	points.push(intersectRect(head, points[points.length - 1]));
	// console.log(edge);
	// const deleteButton = edge.elem.querySelector('.delete-edge');
	// if (deleteButton) {
	// 	const small = tail.y > points[1].y;
	// 	const {x, y} = edge.elem.querySelector('path:first-child').getBoundingClientRect();
	// 	deleteButton.style.transform = `translate(${x - 200}px, ${points[1].y - (small ? -100 : 100)}px)`;
	// }

	return d3
		.line()
		.x(function (d: any) {
			return d.x;
		})
		.y(function (d: any) {
			return d.y;
		})(points);
};

const intersectRect = (node: any, point: any) => {
	let x = node.x;
	let y = node.y;
	let dx = point.x - x;
	let dy = point.y - y;
	let w = parseInt(document.getElementById(node.customId)?.getAttribute('width') || '2') / 2;
	let h = parseInt(document.getElementById(node.customId)?.getAttribute('height') || '2') / 2;
	let sx = 0,
		sy = 0;
	if (Math.abs(dy) * w > Math.abs(dx) * h) {
		if (dy < 0) {
			h = -h;
		}
		sx = dy === 0 ? 0 : (h * dx) / dy;
		sy = h;
	} else {
		if (dx < 0) {
			w = -w;
		}
		sx = w;
		sy = dx === 0 ? 0 : (w * dy) / dx;
	}
	return {
		x: x + sx,
		y: y + sy,
	};
};

const startLineDraw = (d: any, graph: any, d3: any, connectorDot: any) => {
	const node = graph.node(d.subject);
	const points = [{ x: node.x, y: node.y }];
	points.push({
		x: d.x,
		y: d.y,
	});

	document.querySelector('.connector-line')?.setAttribute(
		'd',
		d3
			.line()
			.x(function (e: any) {
				return e.x;
			})
			.y(function (e: any) {
				return e.y;
			})(points)
	);
};

const canSetEdge = (fromNode: any, toNode: any): boolean => {
	if (toNode.dependsOn?.get(fromNode.id) || toNode.kind === 'Environment') {
		return false;
	}
	return true;
};

export const moveNode = function (this: any, graph: any, subject: any, nx: any, ny: any, d3: any) {
	const node = this,
		selectedNode = graph.node(subject);
	const prevX = selectedNode.x,
		prevY = selectedNode.y;

	selectedNode.x += nx;
	selectedNode.y += ny;
	node.setAttribute('transform', 'translate(' + selectedNode.x + ',' + selectedNode.y + ')');

	// let dx = selectedNode.x - prevX,
	//   dy = selectedNode.y - prevY;
	let dx = selectedNode.x - prevX,
		dy = selectedNode.y - prevY;
	// setImmediate(() => {
	graph.edges().forEach(function (e: any) {
		if (e.v == subject || e.w == subject) {
			const edge = graph.edge(e.v, e.w);
			translateEdge(graph.edge(e.v, e.w), dx, dy);
			document.getElementById(edge.customId)?.setAttribute('d', calcPoints(e, graph, d3));
		}
	});
};

const translateEdge = (e: any, dx: any, dy: any) => {
	e.points.forEach(function (p: any) {
		p.x = p.x + dx;
		p.y = p.y + dy;
	});
};

export const nodeDragEnd = (
	d: any,
	graph: any,
	fromElem: any,
	rerender: any,
	setGraph: any,
	connector: any,
	builder: any,
	d3: any
) => {
	const toElem = d.sourceEvent.target.closest('.node')?.querySelector('[data-node-id]');
	if (toElem && fromElem && toElem !== fromElem) {
		const fromNode = fromElem.getAttribute('data-node-id');
		const toNode = toElem.getAttribute('data-node-id');
		if (canSetEdge(graph.node(fromNode).nodeData, graph.node(toNode).nodeData)) {
			setEdge(graph, fromNode, toNode, builder);
			setGraph(graph);
			rerender();
			calcPoints(
				{
					v: fromNode,
					w: toNode,
				},
				graph,
				d3
			);
		}
	}
	connector.remove();
	fromElem = null;
};

export const setEdge = (graph: any, fromNode: any, toNode: any, builder: YEnvironmentBuilder) => {
	if (!graph.node(fromNode).nodeData.kind) {
		graph.node(toNode).nodeData.dependsOn.set(graph.node(fromNode).nodeData.id, graph.node(fromNode).nodeData);
	}
	builder.spec.components.set(graph.node(toNode).nodeData.id, graph.node(toNode).nodeData);
	graph.setEdge(fromNode, toNode, {
		curve: getCurveType('0'),
		arrowheadStyle: 'fill: #333',
	});
};

export const nodeDragMove = function (this: any, d: any, graph: any, d3: any, layout: any, connectorDot: any) {
	if (connectorDot) {
		startLineDraw(d, graph, d3, connectorDot);
		return;
	}
	const selectedNode = graph.node(d.subject);
	moveNode.call(this, graph, d.subject, d.dx, d.dy, d3);
	layout.set(d.subject, {
		x: selectedNode.x,
		y: selectedNode.y,
	});
};

export const nodeDragStart = (d: any, connector: any) => {
	d.sourceEvent.stopPropagation();
	if (!d.sourceEvent.target.classList.contains('connector-dot')) {
		return;
	}
	const connectorLine = document.createElementNS('http://www.w3.org/2000/svg', 'path');
	connectorLine.classList.add('connector-line');
	connector.appendChild(connectorLine);
	return d.sourceEvent.currentTarget.closest('.node').querySelector('[data-node-id]');
};

export const getNode = (id: string, name: string, icon: ReactElement, showPopUp: any, deleteNode: any) => {
	let width = getTextWidth(name) + 80;
	width = width > 160 ? width : 160;

	const groupNode = getSVGNode(
		{
			class: 'roundedCorners',
			id: id,
			filter: `url(#drop_shadow)`,
		},
		'g',
		{
			onmouseenter: (ev: any) => {
				ev.currentTarget.classList.add('visible');
			},
			onmouseleave: (ev: any) => {
				if (!ev.currentTarget.classList.contains('connecting')) {
					ev.currentTarget.classList.remove('visible');
				}
			},
		}
	);
	const rectContainer = getSVGNode(
		{
			height: '70px',
			width: `${width}px`,
			fill: id === 'root' ? 'lightblue' : 'yellowgreen',
			rx: '10',
			class: ``,
		},
		'rect'
	);

	const nameTextNode = getSVGNode(
		{
			x: '65',
			y: '37',
			fill: 'black',
			'font-family': 'DM Sans',
			'font-weight': 'bold',
			'font-size': '14px',
		},
		'text',
		{
			innerHTML: name,
		}
	);

	const img = getSVGNode(
		{
			transform: `translate(${12},${12}) scale(${id === 'root' ? 0.25 : 0.35})`,
		},
		'g',
		{
			innerHTML: ReactDOMServer.renderToString(icon),
		}
	);

	const groupNodeOverlay = getSVGNode(
		{
			class: 'node-overlay',
		},
		'g'
	);

	const rectContainerOverlay = getSVGNode(
		{
			height: '70px',
			width: `${width}px`,
			rx: '10',
			class: `rect-container`,
		},
		'rect'
	);

	const circleConnectCenter = getSVGNode(
		{
			cx: `${width / 2}px`,
			cy: `${35}px`,
			r: '8',
			class: 'connector-dot',
		},
		'circle'
	);

	const deleteNodeSvg = getSVGNode(
		{
			transform: `translate(${width - 10},${-5})`,
			class: 'delete-node',
		},
		'g',
		{
			onclick: () => {
				if (window.confirm('Are you sure you want to delete this node?')) {
					deleteNode();
				}
			},
		}
	);
	const deleteNodeContainer = getSVGNode(
		{
			cx: `${6}px`,
			cy: `${5}px`,
			r: '8',
			fill: '#f00',
			class: 'delete-node__container',
		},
		'circle'
	);

	const closeTextContainer = getSVGNode(
		{
			x: '2.5',
			y: '9',
			class: 'delete-node__text',
		},
		'text',
		{
			innerHTML: 'x',
		}
	);

	deleteNodeSvg.appendChild(deleteNodeContainer);
	deleteNodeSvg.appendChild(closeTextContainer);

	const editNodeSvg = getSVGNode(
		{
			transform: `translate(${width / 3 - 10},${25})`,
			class: 'edit-node',
		},
		'g',
		{
			innerHTML: ReactDOMServer.renderToString(<EditSvg />),
			onclick: () => showPopUp(),
		}
	);

	groupNode.appendChild(rectContainer);
	groupNode.appendChild(nameTextNode);
	groupNode.appendChild(img);
	groupNode.appendChild(groupNodeOverlay);
	groupNodeOverlay.appendChild(rectContainerOverlay);
	groupNodeOverlay.appendChild(circleConnectCenter);
	id !== 'root' && groupNodeOverlay.appendChild(deleteNodeSvg);
	groupNodeOverlay.appendChild(editNodeSvg);
	return groupNode;
};

export const getDeleteEdgeNode = (deleteEdge: any) => {
	return getSVGNode(
		{
			class: 'delete-edge',
		},
		'g',
		{
			innerHTML: ReactDOMServer.renderToString(<DeleteSvg />),
			onclick: () => {
				if (window.confirm('Are you sure you want to delete this edge?')) {
					deleteEdge();
				}
			},
		}
	);
};

export const threeWayToggle = (property: PropertyBinding, setter: any) => {
	return (
		<>
			<label htmlFor={property.key}>{property.label}</label>
			<div
				className="toggler"
				onClick={ev => {
					const value = (ev.nativeEvent.target as HTMLButtonElement).value;
					if (value === 'none') {
						property.set(null);
					} else {
						property.set(value === 'true');
					}
					setter();
				}}>
				<button className={`toggler__false ${property.value === false && 'selected'}`} value="false">
					F
				</button>
				<button className={`toggler__none ${property.value === null && 'selected'}`} value="none">
					-
				</button>
				<button className={`toggler__true ${property.value === true && 'selected'}`} value="true">
					T
				</button>
			</div>
		</>
	);
};

export const getlabelInputPair = (property: PropertyBinding, setter: any) => {
	return (
		<>
			<label htmlFor={property.key}>{property.label}</label>
			<input
				id={property.key}
				type="text"
				name={property.key}
				placeholder={property.label}
				className="shadowy-input environment-builder-form__input"
				defaultValue={property.value}
				onChange={e => {
					if (!e.target.value.trim()) return;
					property.set(e.target.value.trim());
					setter();
				}}
			/>
		</>
	);
};

export const getDropDownList = (
	dropdownData: Set<string>,
	onClickHandler: any,
	classHandler: any,
	readOnly: boolean = false,
	defaultValue: string = '',
	label: string = '',
	additionalInputClasses: string[] = [],
) => {
	const uid = uniqueId('drop-down-');
	const showDropDown = () => {
		const ul = document.querySelector(`#component-builder-drop-down-ul-${uid}`);
		ul?.classList.add('show');
	};
	return (
		<div
			className="component-builder-drop-down"
			onMouseLeave={() => {
				const ul = document.querySelector(`#component-builder-drop-down-ul-${uid}`);
				ul?.classList.remove('show');
			}}>
			<label>{label}</label>
			<Chevron className="drop-down-icon" onClick={showDropDown} />
			<input
				key={uid}
				readOnly={readOnly}
				className={classNames(["shadowy-input", "environment-builder-form__secret-tuple-input", ...additionalInputClasses])}
				type="text"
				placeholder={readOnly ? '' : 'Search...'}
				onClick={showDropDown}
				defaultValue={defaultValue}
				onChange={ev => {
					const v = ev.currentTarget.value.trim();
					const ul = document.querySelector(`#component-builder-drop-down-ul-${uid}`);
					ul?.querySelectorAll('li').forEach(e => {
						if (!e.innerText.includes(v)) {
							e.style.display = 'none';
						} else {
							e.style.display = 'block';
						}
					});
				}}
			/>
			<ul id={`component-builder-drop-down-ul-${uid}`}>
				{[...dropdownData.values()].map(e => (
					<li
						className={classHandler(e)}
						onClick={ev => {
							onClickHandler(e, ev);
						}}
						key={e}>
						{e}
					</li>
				))}
			</ul>
		</div>
	);
};

export const environmentBlueprints = [
	{
		id: 'EKS with Postgres',
		getBlueprint: (addNode: any, graph: any, eb: YEnvironmentBuilder) => {
			const networking = getNetworkingNode();
			const staticAssets = getStaticAssetsNode();
			const platformEks = getPlatformEKSNode();
			const platformEc2 = getPlatformEC2Node();
			const postgres = getPostgres();
			const eksAddons = getEksAddons();
			addNode(networking, false);
			addNode(staticAssets, false);
			addNode(platformEks, false);
			addNode(platformEc2, false);
			addNode(postgres, false);
			addNode(eksAddons, false);
			setEdge(graph, 'root', networking.id, eb);
			setEdge(graph, 'root', staticAssets.id, eb);
			setEdge(graph, networking.id, platformEc2.id, eb);
			setEdge(graph, networking.id, platformEks.id, eb);
			setEdge(graph, platformEks.id, eksAddons.id, eb);
			setEdge(graph, platformEc2.id, postgres.id, eb);
		},
	},
	{
		id: 'Machine Learning',
		getBlueprint: (addNode: any, graph: any, eb: YEnvironmentBuilder) => {
			const networking = getNetworkingNode();
			const staticAssets = getStaticAssetsNode();
			const platformEks = getPlatformEKSNode();
			const platformEc2 = getPlatformEC2Node();
			const postgres = getPostgres();
			const eksAddons = getEksAddons();
			addNode(networking, false);
			addNode(staticAssets, false);
			addNode(platformEks, false);
			addNode(platformEc2, false);
			addNode(postgres, false);
			addNode(eksAddons, false);
			setEdge(graph, 'root', networking.id, eb);
			setEdge(graph, 'root', staticAssets.id, eb);
			setEdge(graph, networking.id, platformEc2.id, eb);
			setEdge(graph, networking.id, platformEks.id, eb);
			setEdge(graph, platformEks.id, eksAddons.id, eb);
			setEdge(graph, platformEc2.id, postgres.id, eb);
		},
	},
	{
		id: 'Data Analytics',
		getBlueprint: (addNode: any, graph: any, eb: YEnvironmentBuilder) => {
			const networking = getNetworkingNode();
			const staticAssets = getStaticAssetsNode();
			const platformEks = getPlatformEKSNode();
			const platformEc2 = getPlatformEC2Node();
			const postgres = getPostgres();
			const eksAddons = getEksAddons();
			addNode(networking, false);
			addNode(staticAssets, false);
			addNode(platformEks, false);
			addNode(platformEc2, false);
			addNode(postgres, false);
			addNode(eksAddons, false);
			setEdge(graph, 'root', networking.id, eb);
			setEdge(graph, 'root', staticAssets.id, eb);
			setEdge(graph, networking.id, platformEc2.id, eb);
			setEdge(graph, networking.id, platformEks.id, eb);
			setEdge(graph, platformEks.id, eksAddons.id, eb);
			setEdge(graph, platformEc2.id, postgres.id, eb);
		},
	},
];

const getNetworkingNode = () => {
	const node = new YComponent('networking');
	node.module.Name.set('vpc');
	node.module.source = 'aws';
	node.variablesFile = new YVariableFile(
		'https://github.com/zl-zlab-tech/checkout-team-config.git',
		'staging/tfvars/networking.tfvars'
	);
	node.tags.add(new YTag('componentType', 'app'));
	node.outputList.add('private_subnets');
	node.outputs.set('private_subnets', new YOutput('private_subnets', true));
	// set template values
	return node;
};

const getStaticAssetsNode = () => {
	const node = new YComponent('static-assets');
	node.module.source = 'aws';
	node.module.Name.set('s3-bucket');
	const variable = new YVariable('bucket', 'dev-zlab-checkout-staging-static-assets');
	node.inputList.set('bucket', variable);
	node.variables.set('bucket', variable);
	node.tags.add(new YTag('componentType', 'data'));
	node.tags.add(new YTag('cloudProvider', 'aws'));
	// set template values
	return node;
};

const getPlatformEKSNode = () => {
	const node = new YComponent('platform-eks');
	node.module.source = 'aws';
	node.module.Name.set('s3-bucket');

	node.variablesFile = new YVariableFile(
		'https://github.com/zl-zlab-tech/checkout-team-config.git',
		'staging/tfvars/platform-eks.tfvars'
	);
	node.tags.add(new YTag('cloudProvider', 'aws'));
	return node;
};

const getPlatformEC2Node = () => {
	const node = new YComponent('platform-ec2');
	node.module.source = 'aws';
	node.module.Name.set('ec2-instance');
	node.variablesFile = new YVariableFile(
		'https://github.com/zl-zlab-tech/checkout-team-config.git',
		'staging/tfvars/ec2.tfvars'
	);
	const variable = new YVariable('subnet_id', '', 'networking.private_subnets[0]');
	variable.choosenDisposition = 'ValueFrom';
	node.inputList.set('subnet_id', variable);
	node.variables.set('subnet_id', variable);
	return node;
};

const getPostgres = () => {
	const node = new YComponent('postgres');
	node.module.source = 'aws';
	node.module.Name.set('s3-bucket');
	const variable = new YVariable('bucket', 'dev-zlab-checkout-staging-postgres', '');
	variable.choosenDisposition = 'Value';
	node.inputList.set('bucket', variable);
	node.variables.set('bucket', variable);
	node.tags.add(new YTag('componentType', 'data'));
	return node;
};

const getEksAddons = () => {
	const node = new YComponent('eks-addons');
	node.module.source = 'aws';
	node.module.Name.set('s3-bucket');
	node.outputList.add('s3_bucket_arn');
	node.outputs.set('name', new YOutput('s3_bucket_arn'));
	return node;
};

import './styles.scss';

import { ZEditor } from 'components/molecules/editor/Editor';
import * as d3 from 'd3';
import dagreD3 from 'dagre-d3';
import { YComponent, YEnvironmentBuilder } from 'models/environment-builder';
import { LocalStorageKey } from 'models/localStorage';
import React, { useEffect, useState } from 'react';
import { useRef } from 'react';
import './styles.scss';
import { ReactComponent as ComputeIcon } from 'assets/images/icons/DAG-View/config.svg';
import { ReactComponent as LayersIcon } from 'assets/images/icons/DAG-View/environment-icon.svg';
import { ReactComponent as Copy } from 'assets/images/icons/copy.svg';
import { ReactComponent as PullRequest } from 'assets/images/icons/git-pull-request.svg';
import {
	getDeleteEdgeNode,
	getNode,
	moveNode,
	nodeDragEnd,
	nodeDragMove,
	nodeDragStart,
	setEdge,
	environmentBlueprints,
} from 'helpers/environment-builder.helper';
import { breadcrumbObservable, pageHeaderObservable } from '../contexts/EnvironmentHeaderContext';
import YAML from 'yaml';
import { TerraformModules } from './TerraformModules';
import ApiClient from 'utils/apiClient';
import { ComponentNodeForm } from './ComponentNodeForm';
import { EnvironmentNodeForm } from './EnvironmentNodeForm';
import { Context } from 'context/argo/ArgoUi';
import { NotificationsApi } from 'components/argo-core/notifications/notification-manager';
import { NotificationType } from 'components/argo-core';
import { EnvironmentTemplates } from './EnvironmentBlueprints';

export const EnvironmentBuilder: React.FC = () => {
	const [graph, setGraph] = useState<any>(null);
	const notificationsManager = React.useContext(Context)?.notifications as NotificationsApi;
	const [builder, setBuilder] = useState<YEnvironmentBuilder | null>(null);
	const [view, setView] = useState<DagViewType>(DagViewType.YAML);
	const [layout, setLayout] = useState<Map<string, any>>(new Map<string, any>());
	const [selectedNodeData, setSelectedNodeData] = useState<any>();
	const [popupVisible, setPopup] = useState<boolean>(false);
	const [yaml, setYaml] = useState<string>('');
	const [environmentBlueprint, setEnvironmentBlueprint] = useState<string>('');
	let fromElem: any = null;
	const connector = useRef<any>(null);
	const containerRef = useRef<HTMLDivElement>(null);
	const contextMenuContainer = useRef<HTMLDivElement>(null);

	let connectorDot: any = null,
		selectedEdge: any = null;

	const nodeDrag = d3
		.drag()
		.on('start', function (d: any) {
			if (d.sourceEvent.target.classList.contains('connector-dot')) {
				connectorDot = d.sourceEvent.target;
			}
			fromElem = nodeDragStart(d, connector.current);
			if (connectorDot) {
				const node = fromElem.parentElement;
				node.querySelector('.roundedCorners').classList.add('visible');
				node.querySelector('.roundedCorners').classList.add('connecting');
			}
		})
		.on('drag', function (d: any) {
			nodeDragMove.call(this, d, graph, d3, layout, connectorDot);
			setGraph(graph);
		})
		.on('end', function (d: any) {
			if (!Boolean(connectorDot)) {
				return;
			}
			const node = fromElem.parentElement;
			node.querySelector('.roundedCorners').classList.remove('visible');
			node.querySelector('.roundedCorners').classList.remove('connecting');
			connectorDot = null;
			nodeDragEnd(
				d,
				graph,
				fromElem,
				rerender,
				setGraph,
				connector.current.querySelector('.connector-line'),
				builder,
				d3
			);
		});

	useEffect(() => {
		pageHeaderObservable.next({
			breadcrumbs: [],
			headerTabs: [],
			pageName: 'Environment Builder',
			filterTitle: '',
			onSearch: () => {},
			onViewChange: () => {},
			buttonText: '',
			checkBoxFilters: null,
		});

		breadcrumbObservable.next({ [LocalStorageKey.TEAMS]: [] });
		initializeDagGraph();
		window.onclick = (e: MouseEvent) => {
			if ((e.target as HTMLElement)?.closest('.context-menu, g.node, button.diff-checker')) {
				return;
			}
			closeNodeUpdatePopup();
		};

		window.onkeydown = (e: KeyboardEvent) => {
			if (e.key === 'Escape') {
				closeNodeUpdatePopup();
			}
		};
	}, []);

	const closeNodeUpdatePopup = () => {
		const contextMenu = contextMenuContainer.current;
		if (contextMenu) {
			contextMenu.querySelector('form')?.reset();
		}
		setPopup(false);
	};

	useEffect(() => {
		if (builder || !graph) {
			return;
		}
		setImmediate(() => {
			const eb = new YEnvironmentBuilder();
			addEnvironmentNode(eb);
		});
	}, [graph]);

	useEffect(() => {
		if (!environmentBlueprint || !graph) {
			return;
		}
		const blueprint = environmentBlueprints.find(e => e.id === environmentBlueprint);
		if (!blueprint) {
			return;
		}
		const eb = new YEnvironmentBuilder();
		addEnvironmentNode(eb);
		blueprint.getBlueprint(addNode, graph, eb);
		setEnvironmentBlueprint('');
		rerender();
	}, [environmentBlueprint, graph]);

	useEffect(() => {
		updateYaml();
	}, [builder]);

	const initializeDagGraph = () => {
		const g = new dagreD3.graphlib.Graph({ directed: true });
		g.setGraph({
			nodesep: 30,
			ranksep: 150,
			rankdir: 'TB',
			ranker: 'network-simplex',
			marginx: 20,
			marginy: 20,
		});

		g.setDefaultEdgeLabel(function () {
			return {};
		});

		const svg = d3.select('.dag-editor');
		svg.attr('height', '100%');
		svg.attr('width', '100%');
		setGraph(g);
		return g;
	};

	const addEnvironmentNode = (eb: YEnvironmentBuilder) => {
		const existing = graph.node('root');
		graph.setNode('root', {
			label: getNode('root', eb.metadata.name, <LayersIcon />, showPopUp.bind(null, eb), () => {}),
			labelType: 'svg',
			style: 'fill: transparent;',
			paddingX: 0,
			paddingY: 0,
			nodeData: eb,
		});
		setBuilder(eb);

		rerender();

		const node = graph.node('root');
		if (existing) {
			return;
		}
		const rect = node.elem.getBoundingClientRect();
		if (containerRef.current) node.x = (containerRef.current.getBoundingClientRect().width - rect.width) / 2;
		node.y = rect.height;

		node.elem.setAttribute('transform', 'translate(' + node.x + ',' + node.y + ')');
		node.elem.setAttribute('rx', '10px');
		node.elem.setAttribute('ry', '10px');
		layout.set('root', {
			x: node.x,
			y: node.y,
		});
	};

	const addNode = (nodeData: YComponent, render: boolean = true, x?: number, y?: number) => {
		graph.setNode(nodeData.id, {
			label: getNode(
				nodeData.Name.value,
				nodeData.Name.value,
				<ComputeIcon />,
				showPopUp.bind(null, nodeData),
				deleteNode.bind(null, nodeData)
			),
			labelType: 'svg',
			style: 'fill: transparent;',
			paddingX: 0,
			paddingY: 0,
			nodeData,
		});

		if (builder?.spec.components.has(nodeData.id)) {
			builder.spec.components.set(nodeData.id, nodeData);
			builder.spec.components.forEach(comp => {
				if (comp.dependsOn.has(nodeData.id)) {
					comp.dependsOn.set(nodeData.id, nodeData);
				}
			});
		}
		render && rerender();

		if (!x || !y) {
			return;
		}
		const node = graph.node(nodeData.id);
		node.x = x;
		node.y = y;
		node.elem.setAttribute('transform', 'translate(' + node.x + ',' + node.y + ')');
		node.elem.setAttribute('rx', '10');
		node.elem.setAttribute('ry', '10');
		layout.set(nodeData.id, {
			x,
			y,
		});
	};

	const updateIds = (svg: any) => {
		svg.selectAll('g.node > rect').attr('id', function (d: any) {
			return 'custom_' + d;
		});
		svg.selectAll('g.node > rect').attr('data-node-id', function (d: any) {
			return '' + d;
		});
		svg.selectAll('g.edgePath > path').attr('id', function (e: any) {
			return e.v + '-' + e.w;
		});

		console.log(graph.edges());
		graph.nodes().forEach(function (v: any) {
			const node = graph.node(v);
			node.customId = 'custom_' + v;
		});
		graph.edges().forEach(function (e: any) {
			const edge = graph.edge(e.v, e.w);
			edge.customId = e.v + '-' + e.w;
		});
		graph.nodes().forEach(function (v: any) {
			const node = graph.node(v);
			const defaultTransform = node.elem.getAttribute('transform');
			const parts = /translate\(\s*([^\s,)]+)[ ,]?([^\s,)]+)?/.exec(defaultTransform)?.map(e => parseInt(e)) || [
				0,
				0,
				0,
			];
			if (layout.has(v)) {
				const xy = layout.get(v);
				moveNode.call(node.elem, graph, v, xy.x - parts[1], xy.y - parts[2], d3);
			} else {
				node.x = parts[1];
				node.y = parts[2];
				node.elem.setAttribute('transform', 'translate(' + node.x + ',' + node.y + ')');
				node.elem.setAttribute('rx', '10');
				node.elem.setAttribute('ry', '10');
				layout.set(v, {
					x: parts[1],
					y: parts[2],
				});
			}
		});
	};

	const attachEvents = (svg: any) => {
		nodeDrag(svg.selectAll('g.node'));
		svg.selectAll('g.edgePath').on('click', function (this: any, d: any) {
			d.preventDefault();
			deleteEdge(this);
			// selectedEdge?.classList?.remove('selected');
			// this.classList.add('selected');
			// selectedEdge = this;
		});

		svg.selectAll('g.node').on('click', function (this: any, d: any) {
			d.preventDefault();
		});
	};

	const showPopUp = (nodeData: any) => {
		setSelectedNodeData(nodeData);
		setPopup(true);
		const div = contextMenuContainer.current;
		if (div) {
			div.style.top = `0px`;
			div.style.left = `0px`;
		}
	};

	const nodeUpdateCallback = () => {
		if (selectedNodeData.kind === 'Environment') {
			addEnvironmentNode(selectedNodeData);
		} else {
			addNode(selectedNodeData);
		}
	};

	const deleteEdge = (ref: any) => {
		const edges: any = d3.select(ref).data()[0];
		const vcomponent: YComponent = graph.node(edges.v).nodeData;
		const wcomponent: YComponent = graph.node(edges.w).nodeData;
		graph.removeEdge(edges.v, edges.w);
		wcomponent.dependsOn.delete(vcomponent.id);
		setGraph(graph);
		rerender();
	};

	const deleteNode = (selectedNodeData: any) => {
		closeNodeUpdatePopup();
		const node = selectedNodeData;
		let parentNode = 'root';
		let children: string[] = [];
		graph.edges().forEach((e: any) => {
			if (node.id === e.w) {
				parentNode = e.v;
			}
			if (node.id === e.v) {
				children.push(e.w);
			}
			if (node.id === e.w || node.id === e.v) deleteEdge(graph.edge(e.v, e.w).elem);
		});
		builder?.spec.components.delete(node.id);
		graph.removeNode(node.id);
		if (builder && parentNode && children.length > 0) {
			children.forEach(c => {
				setEdge(graph, parentNode, c, builder);
			});
		}
		setGraph(graph);
		rerender();
	};

	const addRemoveEdgeSvg = (selection: d3.Selection<d3.BaseType, unknown, d3.BaseType, unknown>) => {
		selection.each(function (this: any, d: any) {
			// this.classList.remove('selected');
			if (this.querySelector('.delete-edge')) {
				return;
			}
			// this.appendChild(getDeleteEdgeNode(deleteEdge.bind(null, this)));
		});
	};

	const rerender = () => {
		const render = new dagreD3.render();
		const svg = d3.select('.dag-editor');
		const inner: any = svg.select('g');
		render(inner, graph);
		updateIds(svg);
		addRemoveEdgeSvg(svg.selectAll('g.edgePath'));
		attachEvents(svg);
		updateYaml();
	};

	const updateYaml = () => {
		if (builder) setYaml(YAML.stringify(builder.yamlObject()));
	};

	const renderPopup = () => {
		if (!popupVisible) return;
		return (
			selectedNodeData && (
				<div className="context-menu" style={{ display: 'flex' }} ref={contextMenuContainer}>
					{selectedNodeData.kind === 'Environment' ? (
						<EnvironmentNodeForm
							environmentSetter={{ selectedNodeData, setSelectedNodeData }}
							updateCallback={nodeUpdateCallback}
						/>
					) : (
						// getEnvironmentNodeForm({ selectedNodeData, setSelectedNodeData }, nodeUpdateCallback)
						<ComponentNodeForm
							componentSetter={{ selectedNodeData, setSelectedNodeData }}
							updateCallback={nodeUpdateCallback}
							deleteCallback={() => {}}
						/>
					)}
				</div>
			)
		);
	};

	const renderPullRequest = () => {
		return (
			<button className="controls-btn view-toggle" title="Open Pull Request">
				<PullRequest viewBox="0 0 16 16" height={24} width={24} />{' '}
				<strong style={{ marginLeft: 5 }}>Pull Request</strong>
			</button>
		);
	};

	return (
		<>
			<div className="controls">
				<div>
					<button
						className="controls-btn view-toggle"
						onClick={() => {
							setView(view === DagViewType.DAG ? DagViewType.YAML : DagViewType.DAG);
						}}>
						{view === DagViewType.DAG ? DagViewType.YAML : DagViewType.DAG}
					</button>
				</div>
				<div>
					{renderPullRequest()}
					<button
						title="Copy YML"
						className="controls-btn view-toggle"
						onClick={e => {
							navigator.clipboard.writeText(yaml);
							notificationsManager.show(
								{
									content: 'YML Copied!',
									type: NotificationType.Success,
								},
								2000
							);
						}}>
						<Copy />
					</button>
				</div>
			</div>
			<div className="dag-editor-container">
				<div>
					<TerraformModules />
					<hr/>
					<EnvironmentTemplates />
				</div>
				<div
					onDrop={ev => {
						try {
							const data = JSON.parse(ev.dataTransfer.getData('text'));
							if (data.type === 'template') {
								initializeDagGraph();
								setEnvironmentBlueprint(data.id);
							} else {
								const comp = new YComponent(data.name, '', data.name, '');
								ApiClient.get(`/terraform-external/modules/aws/${data.name}`).then(({ data }) => {
									const { root } = data as any;
									const { outputs, inputs } = root;
									comp.outputList = new Set<string>(outputs.map((e: any) => e.name));
									comp.inputList = new Map<string, any>(inputs.map((e: any) => [e.name, e]));
									setSelectedNodeData(comp);
								});
								addNode(comp, true, (ev.nativeEvent as any).layerX, (ev.nativeEvent as any).layerY);
								console.log(ev.dataTransfer.getData('text'));
							}
							containerRef.current?.classList.remove('drop-container');
						} catch (err) {
							containerRef.current?.classList.remove('drop-container');
						}
					}}
					onDragOver={ev => {
						ev.preventDefault();
						containerRef.current?.classList.add('drop-container');
					}}
					onDragLeave={ev => {
						containerRef.current?.classList.remove('drop-container');
					}}
					ref={containerRef}
					className="svg-container"
					onDoubleClick={(e: any) => {
						const component = new YComponent();
						// component.variablesFile.path = component.name + '.tfvars';
						addNode(component, true, e.nativeEvent.layerX, e.nativeEvent.layerY);
					}}>
					<svg className="dag-editor">
						<defs>
							<filter
								id="drop_shadow"
								filterUnits="objectBoundingBox"
								x="-50%"
								y="-50%"
								width="200%"
								height="200%">
								<feDropShadow dx="0" dy="0" stdDeviation="5" floodColor="#aaa" floodOpacity="1" />
							</filter>
						</defs>
						<g id="graph-container">
							<g id="connector-line" ref={connector}></g>
						</g>
					</svg>
				</div>
				<div className={`generated-yaml ${view === DagViewType.YAML ? '' : 'collapsed'}`}>
					<ZEditor
						data={yaml}
						options={{
							fontSize: '16px',
						}}
					/>
				</div>
				{renderPopup()}
			</div>
		</>
	);
};

enum DagViewType {
	DAG = 'DAG View',
	YAML = 'Split View',
}

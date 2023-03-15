import { ReactComponent as ComputeIcon } from 'assets/images/icons/DAG-View/config.svg';
import { ReactComponent as LayersIcon } from 'assets/images/icons/DAG-View/environment-icon.svg';
import dagre from 'dagre';
import { AuditStatus, ESyncStatus, ZSyncStatus } from 'models/argo.models';
import { Component, Environment } from 'models/entity.type';
import { useCallback } from 'react';
import { Edge, getSimpleBezierPath, Handle, MarkerType, Node, Position, useStore } from 'reactflow';
import { DagNode } from './DagNode';

export const generateRootNode = (environment: Environment) => {
	const data = {
		icon: <LayersIcon />,
		...environment,
	};
	return (
		<>
			<Handle
				id="a"
				className="targetHandle"
				style={{ zIndex: 2, top: 30 }}
				position={Position.Bottom}
				type="source"
				isConnectable={true}
			/>
			<Handle
				id="b"
				className="targetHandle"
				style={{
					top: 0,
				}}
				position={Position.Top}
				type="target"
				isConnectable={true}
			/>
			<DagNode
				data={{
					cost: data.estimatedCost,
					name: data.name,
					icon: data.icon,
					status: data.status as ZSyncStatus,
					timestamp: data.lastReconcileDatetime,
					operation: 'Provision',
					isSkipped: false,
				}}
			/>
		</>
	);
};

export const generateComponentNode = (component: Component) => {
	const data = {
		...component,
		icon: <ComputeIcon />,
		isSkipped: [AuditStatus.SkippedProvision, AuditStatus.SkippedDestroy].includes(component.lastAuditStatus),
	};

	return (
		<>
			<Handle
				id="a"
				className="targetHandle"
				style={{ zIndex: 2, top: 30 }}
				position={Position.Bottom}
				type="source"
				isConnectable={true}
			/>
			<Handle
				id="b"
				className="targetHandle"
				style={{ top: 0 }}
				position={Position.Top}
				type="target"
				isConnectable={true}
			/>
			<DagNode
				data={{
					cost: data.estimatedCost,
					name: data.name,
					icon: data.icon,
					status: data.status as ZSyncStatus,
					timestamp: data.lastReconcileDatetime,
					operation: data.isDestroyed ? 'Destroy' : 'Provision',
					isSkipped: data.isSkipped,
				}}
			/>
		</>
	);
};

export const initializeLayout = () => {
	const dagreGraph = new dagre.graphlib.Graph();
	dagreGraph.setDefaultEdgeLabel(() => ({}));

	const nodeWidth = 250;
	const nodeHeight = 60;

	const getLayoutedElements = (nodes: any, edges: any, direction = 'TB') => {
		const isHorizontal = direction === 'LR';
		dagreGraph.setGraph({ rankdir: direction });

		nodes.forEach((node: any) => {
			dagreGraph.setNode(node.id, { width: nodeWidth, height: nodeHeight });
		});
		edges.forEach((edge: any) => {
			dagreGraph.setEdge(edge.source, edge.target);
		});

		dagre.layout(dagreGraph);
		nodes.forEach((node: any) => {
			const nodeWithPosition = dagreGraph.node(node.id);
			node.targetPosition = isHorizontal ? 'left' : 'top';
			node.sourcePosition = isHorizontal ? 'right' : 'bottom';

			// We are shifting the dagre node position (anchor=center center) to the top left
			// so it matches the React Flow node anchor point (top left).
			node.position = {
				x: nodeWithPosition.x - nodeWidth / 2,
				y: (nodeWithPosition.y - nodeHeight / 2) * 1.2,
			};

			return node;
		});

		return { nodes, edges };
	};

	return {
		getLayoutedElements,
	};
};

const getNode = (id: string, label: JSX.Element): Node => {
	return {
		id,
		data: { label },
		position: { x: 0, y: 0 },
		style: { padding: 0, border: 0 },
		draggable: false,
		selectable: false,
	};
};

const getEdge = (id: string, source: string, target: string): Edge => {
	return {
		id,
		source,
		target,
		type: 'smart',
		markerEnd: { type: MarkerType.ArrowClosed, width: 15, height: 15, color: '#333' },
		sourceHandle: 'a',
		targetHandle: 'b',
	};
};

export const generateNodesAndEdges = (environment: Environment) => {
	environment.dag = [
		{
			name: 'app-addons',
			tags: [
				{
					name: 'componentName',
					value: 'app-addons',
				},
			],
			type: 'terraform',
			module: {
				path: 'common-app-addons',
				source: 'https://github.com/zl-compuzest-tech/terraform-modules.git',
			},
			secrets: [
				{
					key: 'auth0_web_client_id',
					name: 'AUTH0_WEB_CLIENT_ID',
					scope: 'environment',
				},
				{
					key: 'auth0_web_secret',
					name: 'AUTH0_WEB_SECRET',
					scope: 'environment',
				},
				{
					key: 'auth0_api_client_id',
					name: 'AUTH0_API_CLIENT_ID',
					scope: 'environment',
				},
				{
					key: 'auth0_api_secret',
					name: 'AUTH0_API_SECRET',
					scope: 'environment',
				},
				{
					key: 'db_password',
					name: 'zlifecycle_api_db_password',
					scope: 'team',
				},
				{
					key: 'shared_aws_access_key_id',
					name: 'shared_aws_access_key_id',
					scope: 'team',
				},
				{
					key: 'shared_aws_secret_access_key',
					name: 'shared_aws_secret_access_key',
					scope: 'team',
				},
			],
			dependsOn: ['k8s-addons'],
			overlayFiles: [
				{
					paths: ['common/tf/data-common.tf', 'common/tf/app-addons.tf'],
					source: 'git@github.com:zl-compuzest-tech/dev-config.git',
				},
			],
			variablesFile: {
				path: 'common/tfvars/app-addons.tfvars',
				source: 'git@github.com:zl-compuzest-tech/dev-config.git',
			},
		},
		{
			name: 'app-db',
			tags: [
				{
					name: 'componentName',
					value: 'app-db',
				},
			],
			type: 'terraform',
			module: {
				path: 'common-db',
				source: 'https://github.com/zl-compuzest-tech/terraform-modules.git',
			},
			secrets: [
				{
					key: 'db_password',
					name: 'rds_master_password',
					scope: 'team',
				},
			],
			dependsOn: ['networking', 'vpn', 'k8s-addons'],
			variablesFile: {
				path: 'common/tfvars/app-db.tfvars',
				source: 'git@github.com:zl-compuzest-tech/dev-config.git',
			},
			destroyProtection: true,
		},
		{
			name: 'argo-workflow',
			tags: [
				{
					name: 'componentName',
					value: 'argo-workflow',
				},
			],
			type: 'terraform',
			module: {
				path: 'argo-workflow',
				source: 'https://github.com/zl-compuzest-tech/terraform-modules.git',
			},
			secrets: [
				{
					key: 'git_ssh_private_key',
					name: 'git_ssh_private_key',
					scope: 'team',
				},
			],
			dependsOn: ['networking', 'eks', 'k8s-addons'],
			overlayFiles: [
				{
					paths: ['common/tf/app-addons.tf'],
					source: 'git@github.com:zl-compuzest-tech/dev-config.git',
				},
			],
			variablesFile: {
				path: 'common/tfvars/argo-workflow.tfvars',
				source: 'git@github.com:zl-compuzest-tech/dev-config.git',
			},
			destroyProtection: true,
		},
		{
			name: 'eks',
			tags: [
				{
					name: 'componentName',
					value: 'eks',
				},
			],
			type: 'terraform',
			module: {
				path: 'common-k8s',
				source: 'https://github.com/zl-compuzest-tech/terraform-modules.git',
			},
			outputs: [
				{
					name: 'oidc_provider_arn',
				},
				{
					name: 'cluster_oidc_issuer_url',
				},
				{
					name: 'cluster_id',
				},
				{
					name: 'cluster_endpoint',
				},
				{
					name: 'cluster_certificate_authority_data',
				},
			],
			dependsOn: ['networking'],
			overlayFiles: [
				{
					paths: ['common/tf/eks.tf'],
					source: 'git@github.com:zl-compuzest-tech/dev-config.git',
				},
			],
			variablesFile: {
				path: 'common/tfvars/eks.tfvars',
				source: 'git@github.com:zl-compuzest-tech/dev-config.git',
			},
		},
		{
			name: 'k8s-addons',
			tags: [
				{
					name: 'componentName',
					value: 'k8s-addons',
				},
			],
			type: 'terraform',
			module: {
				path: 'common-k8s-addons',
				source: 'https://github.com/zl-compuzest-tech/terraform-modules.git',
			},
			outputs: [
				{
					name: 'zlifecycle_system_namespace',
				},
				{
					name: 'zlifecycle_executor_namespace',
				},
			],
			dependsOn: ['networking', 'eks'],
			overlayFiles: [
				{
					paths: ['common/tf/k8s-addons.tf'],
					source: 'git@github.com:zl-compuzest-tech/dev-config.git',
				},
			],
			variablesFile: {
				path: 'common/tfvars/k8s-addons.tfvars',
				source: 'git@github.com:zl-compuzest-tech/dev-config.git',
			},
		},
		{
			name: 'networking',
			type: 'terraform',
			module: {
				name: 'vpc',
				source: 'aws',
			},
			outputs: [
				{
					name: 'vpc_id',
				},
				{
					name: 'public_subnets',
				},
				{
					name: 'private_subnets',
				},
			],
			variablesFile: {
				path: 'common/tfvars/networking.tfvars',
				source: 'git@github.com:zl-compuzest-tech/dev-config.git',
			},
			destroyProtection: true,
		},
		{
			name: 'system',
			type: 'terraform',
			module: {
				path: 'system',
				source: 'https://github.com/zl-compuzest-tech/terraform-modules.git',
			},
			variablesFile: {
				path: 'common/tfvars/system.tfvars',
				source: 'git@github.com:zl-compuzest-tech/dev-config.git',
			},
			destroyProtection: true,
		},
		{
			name: 'vpn',
			tags: [
				{
					name: 'componentName',
					value: 'vpn',
				},
			],
			type: 'terraform',
			module: {
				path: 'vpn',
				source: 'https://github.com/zl-compuzest-tech/terraform-modules.git',
			},
			outputs: [
				{
					name: 'private_ip',
				},
				{
					name: 'public_ip',
				},
				{
					name: 'security_group_id',
				},
			],
			dependsOn: ['networking'],
			variablesFile: {
				path: 'common/tfvars/vpn.tfvars',
				source: 'git@github.com:zl-compuzest-tech/dev-config.git',
			},
			destroyProtection: true,
		},
	] as any;
	const nodes: Node[] = [];
	const edges: Edge[] = [];

	const envNode = getNode(environment.argoId.toString(), generateRootNode(environment));
	nodes.push(envNode);

	const components = [
		{
			"id": 11,
			"name": "app-addons",
			"type": "terraform",
			"status": "provisioned",
			"estimatedCost": 0,
			"lastReconcileDatetime": "2023-03-12T08:05:19.000Z",
			"duration": -1,
			"lastWorkflowRunId": "dev-common-app-addons-nxt6v",
			"isDestroyed": false,
			"costResources": [],
			"envId": 2,
			"orgId": 1,
			"lastAuditStatus": "Provisioned",
			"argoId": "key0"
		},
		{
			"id": 10,
			"name": "app-db",
			"type": "terraform",
			"status": "provisioned",
			"estimatedCost": 13.56,
			"lastReconcileDatetime": "2023-03-06T01:33:18.000Z",
			"duration": -1,
			"lastWorkflowRunId": "dev-common-app-db-b75kg",
			"isDestroyed": false,
			"costResources": [
				{
					"name": "module.app-db.aws_db_instance.this",
					"tags": {
						"Name": "dev-zlifecycle-db",
						"Environment": "dev"
					},
					"metadata": {},
					"hourlyCost": "0.0185753424657534245",
					"monthlyCost": "13.56",
					"costComponents": [
						{
							"name": "Database instance (on-demand, Single-AZ, db.t3.micro)",
							"unit": "hours",
							"price": "0.017",
							"hourlyCost": "0.017",
							"monthlyCost": "12.41",
							"hourlyQuantity": "1",
							"monthlyQuantity": "730"
						},
						{
							"name": "Storage (general purpose SSD, gp2)",
							"unit": "GB",
							"price": "0.115",
							"hourlyCost": "0.0015753424657534245",
							"monthlyCost": "1.15",
							"hourlyQuantity": "0.0136986301369863",
							"monthlyQuantity": "10"
						}
					]
				},
				{
					"name": "module.app-db.aws_route53_record.db",
					"metadata": {},
					"hourlyCost": null,
					"monthlyCost": null,
					"costComponents": [
						{
							"name": "Standard queries (first 1B)",
							"unit": "1M queries",
							"price": "0.4",
							"hourlyCost": null,
							"monthlyCost": null,
							"hourlyQuantity": null,
							"monthlyQuantity": null
						},
						{
							"name": "Latency based routing queries (first 1B)",
							"unit": "1M queries",
							"price": "0.6",
							"hourlyCost": null,
							"monthlyCost": null,
							"hourlyQuantity": null,
							"monthlyQuantity": null
						},
						{
							"name": "Geo DNS queries (first 1B)",
							"unit": "1M queries",
							"price": "0.7",
							"hourlyCost": null,
							"monthlyCost": null,
							"hourlyQuantity": null,
							"monthlyQuantity": null
						}
					]
				}
			],
			"envId": 2,
			"orgId": 1,
			"lastAuditStatus": "skipped_provision",
			"argoId": "key1"
		},
		{
			"id": 9,
			"name": "argo-workflow",
			"type": "terraform",
			"status": "provisioned",
			"estimatedCost": 0,
			"lastReconcileDatetime": "2023-03-07T02:50:43.000Z",
			"duration": -1,
			"lastWorkflowRunId": "dev-common-argo-workflow-c7xd6",
			"isDestroyed": false,
			"costResources": [],
			"envId": 2,
			"orgId": 1,
			"lastAuditStatus": "skipped_provision",
			"argoId": "key2"
		},
		{
			"id": 5,
			"name": "eks",
			"type": "terraform",
			"status": "provisioned",
			"estimatedCost": 587.288,
			"lastReconcileDatetime": "2023-03-06T01:30:25.000Z",
			"duration": -1,
			"lastWorkflowRunId": "dev-common-eks-mqwxm",
			"isDestroyed": false,
			"costResources": [
				{
					"name": "module.eks.module.eks.aws_cloudwatch_log_group.this[0]",
					"tags": {
						"Cluster": "dev-eks",
						"Terraform": "true"
					},
					"metadata": {},
					"hourlyCost": null,
					"monthlyCost": null,
					"costComponents": [
						{
							"name": "Data ingested",
							"unit": "GB",
							"price": "0.5",
							"hourlyCost": null,
							"monthlyCost": null,
							"hourlyQuantity": null,
							"monthlyQuantity": null
						},
						{
							"name": "Archival Storage",
							"unit": "GB",
							"price": "0.03",
							"hourlyCost": null,
							"monthlyCost": null,
							"hourlyQuantity": null,
							"monthlyQuantity": null
						},
						{
							"name": "Insights queries data scanned",
							"unit": "GB",
							"price": "0.005",
							"hourlyCost": null,
							"monthlyCost": null,
							"hourlyQuantity": null,
							"monthlyQuantity": null
						}
					]
				},
				{
					"name": "module.eks.module.eks.aws_eks_cluster.this[0]",
					"tags": {
						"Cluster": "dev-eks",
						"Terraform": "true"
					},
					"metadata": {},
					"hourlyCost": "0.1",
					"monthlyCost": "73",
					"costComponents": [
						{
							"name": "EKS cluster",
							"unit": "hours",
							"price": "0.1",
							"hourlyCost": "0.1",
							"monthlyCost": "73",
							"hourlyQuantity": "1",
							"monthlyQuantity": "730"
						}
					]
				},
				{
					"name": "module.eks.module.system_managed_node_group[0].aws_eks_node_group.this[0]",
					"tags": {
						"Name": "dev-eks-system",
						"CalicoCNI": "4303740921264266020"
					},
					"metadata": {},
					"hourlyCost": "0.70450410958904108",
					"monthlyCost": "514.288",
					"subresources": [
						{
							"name": "module.eks.aws_launch_template.launch_template",
							"metadata": {},
							"hourlyCost": "0.70450410958904108",
							"monthlyCost": "514.288",
							"subresources": [
								{
									"name": "block_device_mapping[0]",
									"metadata": {},
									"hourlyCost": "0.0273972602739726",
									"monthlyCost": "20",
									"costComponents": [
										{
											"name": "Storage (general purpose SSD, gp2)",
											"unit": "GB",
											"price": "0.1",
											"hourlyCost": "0.0273972602739726",
											"monthlyCost": "20",
											"hourlyQuantity": "0.273972602739726",
											"monthlyQuantity": "200"
										}
									]
								}
							],
							"costComponents": [
								{
									"name": "Instance usage (Linux/UNIX, on-demand, t3.xlarge)",
									"unit": "hours",
									"price": "0.1664",
									"hourlyCost": "0.6656",
									"monthlyCost": "485.888",
									"hourlyQuantity": "4",
									"monthlyQuantity": "2920"
								},
								{
									"name": "EC2 detailed monitoring",
									"unit": "metrics",
									"price": "0.3",
									"hourlyCost": "0.01150684931506848",
									"monthlyCost": "8.4",
									"hourlyQuantity": "0.0383561643835616",
									"monthlyQuantity": "28"
								},
								{
									"name": "CPU credits",
									"unit": "vCPU-hours",
									"price": "0.05",
									"hourlyCost": "0",
									"monthlyCost": "0",
									"hourlyQuantity": "0",
									"monthlyQuantity": "0"
								}
							]
						}
					]
				}
			],
			"envId": 2,
			"orgId": 1,
			"lastAuditStatus": "skipped_provision",
			"argoId": "key3"
		},
		{
			"id": 6,
			"name": "k8s-addons",
			"type": "terraform",
			"status": "provisioned",
			"estimatedCost": 0.5,
			"lastReconcileDatetime": "2023-03-06T01:31:52.000Z",
			"duration": -1,
			"lastWorkflowRunId": "dev-common-k8s-addons-f44mh",
			"isDestroyed": false,
			"costResources": [
				{
					"name": "module.k8s-addons.aws_route53_zone.internal",
					"metadata": {},
					"hourlyCost": "0.0006849315068493",
					"monthlyCost": "0.5",
					"costComponents": [
						{
							"name": "Hosted zone",
							"unit": "months",
							"price": "0.5",
							"hourlyCost": "0.0006849315068493",
							"monthlyCost": "0.5",
							"hourlyQuantity": "0.0013698630136986",
							"monthlyQuantity": "1"
						}
					]
				}
			],
			"envId": 2,
			"orgId": 1,
			"lastAuditStatus": "skipped_provision",
			"argoId": "key4"
		},
		{
			"id": 4,
			"name": "networking",
			"type": "terraform",
			"status": "provisioned",
			"estimatedCost": 32.85,
			"lastReconcileDatetime": "2023-03-06T01:28:54.000Z",
			"duration": -1,
			"lastWorkflowRunId": "dev-common-networking-56np8",
			"isDestroyed": false,
			"costResources": [
				{
					"name": "module.networking.aws_nat_gateway.this[0]",
					"tags": {
						"Name": "dev-common-us-east-1a",
						"Terraform": "true",
						"Environment": "dev-common",
						"kubernetes.io/cluster/dev-eks": "shared"
					},
					"metadata": {},
					"hourlyCost": "0.045",
					"monthlyCost": "32.85",
					"costComponents": [
						{
							"name": "NAT gateway",
							"unit": "hours",
							"price": "0.045",
							"hourlyCost": "0.045",
							"monthlyCost": "32.85",
							"hourlyQuantity": "1",
							"monthlyQuantity": "730"
						},
						{
							"name": "Data processed",
							"unit": "GB",
							"price": "0.045",
							"hourlyCost": null,
							"monthlyCost": null,
							"hourlyQuantity": null,
							"monthlyQuantity": null
						}
					]
				}
			],
			"envId": 2,
			"orgId": 1,
			"lastAuditStatus": "skipped_provision",
			"argoId": "key5"
		},
		{
			"id": 8,
			"name": "system",
			"type": "terraform",
			"status": "provisioned",
			"estimatedCost": 0,
			"lastReconcileDatetime": "2023-03-06T01:28:50.000Z",
			"duration": -1,
			"lastWorkflowRunId": "dev-common-system-7jfhx",
			"isDestroyed": false,
			"costResources": [
				{
					"name": "module.system.aws_s3_bucket.system",
					"metadata": {},
					"hourlyCost": null,
					"monthlyCost": null,
					"subresources": [
						{
							"name": "Standard",
							"metadata": {},
							"hourlyCost": null,
							"monthlyCost": null,
							"costComponents": [
								{
									"name": "Storage",
									"unit": "GB",
									"price": "0.023",
									"hourlyCost": null,
									"monthlyCost": null,
									"hourlyQuantity": null,
									"monthlyQuantity": null
								},
								{
									"name": "PUT, COPY, POST, LIST requests",
									"unit": "1k requests",
									"price": "0.005",
									"hourlyCost": null,
									"monthlyCost": null,
									"hourlyQuantity": null,
									"monthlyQuantity": null
								},
								{
									"name": "GET, SELECT, and all other requests",
									"unit": "1k requests",
									"price": "0.0004",
									"hourlyCost": null,
									"monthlyCost": null,
									"hourlyQuantity": null,
									"monthlyQuantity": null
								},
								{
									"name": "Select data scanned",
									"unit": "GB",
									"price": "0.002",
									"hourlyCost": null,
									"monthlyCost": null,
									"hourlyQuantity": null,
									"monthlyQuantity": null
								},
								{
									"name": "Select data returned",
									"unit": "GB",
									"price": "0.0007",
									"hourlyCost": null,
									"monthlyCost": null,
									"hourlyQuantity": null,
									"monthlyQuantity": null
								}
							]
						}
					]
				}
			],
			"envId": 2,
			"orgId": 1,
			"lastAuditStatus": "skipped_provision",
			"argoId": "key6"
		},
		{
			"id": 7,
			"name": "vpn",
			"type": "terraform",
			"status": "provisioned",
			"estimatedCost": 10.792,
			"lastReconcileDatetime": "2023-03-06T01:30:16.000Z",
			"duration": -1,
			"lastWorkflowRunId": "dev-common-vpn-t4fv5",
			"isDestroyed": false,
			"costResources": [
				{
					"name": "module.vpn.module.vpn.aws_cloudwatch_metric_alarm.default[0]",
					"metadata": {},
					"hourlyCost": "0.00013698630136986",
					"monthlyCost": "0.1",
					"costComponents": [
						{
							"name": "Standard resolution",
							"unit": "alarm metrics",
							"price": "0.1",
							"hourlyCost": "0.00013698630136986",
							"monthlyCost": "0.1",
							"hourlyQuantity": "0.0013698630136986",
							"monthlyQuantity": "1"
						}
					]
				},
				{
					"name": "module.vpn.module.vpn.aws_instance.default[0]",
					"tags": {
						"Name": "common-dev-ec2-vpn",
						"Stage": "dev",
						"Namespace": "common"
					},
					"metadata": {},
					"hourlyCost": "0.01464657534246575",
					"monthlyCost": "10.692",
					"subresources": [
						{
							"name": "root_block_device",
							"metadata": {},
							"hourlyCost": "0.00136986301369863",
							"monthlyCost": "1",
							"costComponents": [
								{
									"name": "Storage (general purpose SSD, gp2)",
									"unit": "GB",
									"price": "0.1",
									"hourlyCost": "0.00136986301369863",
									"monthlyCost": "1",
									"hourlyQuantity": "0.0136986301369863",
									"monthlyQuantity": "10"
								}
							]
						}
					],
					"costComponents": [
						{
							"name": "Instance usage (Linux/UNIX, on-demand, t3.micro)",
							"unit": "hours",
							"price": "0.0104",
							"hourlyCost": "0.0104",
							"monthlyCost": "7.592",
							"hourlyQuantity": "1",
							"monthlyQuantity": "730"
						},
						{
							"name": "EC2 detailed monitoring",
							"unit": "metrics",
							"price": "0.3",
							"hourlyCost": "0.00287671232876712",
							"monthlyCost": "2.1",
							"hourlyQuantity": "0.0095890410958904",
							"monthlyQuantity": "7"
						},
						{
							"name": "CPU credits",
							"unit": "vCPU-hours",
							"price": "0.05",
							"hourlyCost": "0",
							"monthlyCost": "0",
							"hourlyQuantity": "0",
							"monthlyQuantity": "0"
						}
					]
				}
			],
			"envId": 2,
			"orgId": 1,
			"lastAuditStatus": "skipped_provision",
			"argoId": "key7"
		}
	]; //EntityStore.getInstance().getComponentsByEnvId(environment.id);

	for (let e of environment.dag) {
		const item = components.find(c => c.name === e.name) as Component;
		if (!item) continue;

		nodes.push(getNode(item.argoId.toString(), generateComponentNode(item)));

		if (!e.dependsOn?.length) {
			edges.push(
				getEdge(`e${environment.argoId}-${item.argoId}`, environment.argoId.toString(), item.argoId.toString())
			);
		} else {
			e.dependsOn.forEach(d => {
				const dc = components.find(ic => ic.name === d);
				if (dc) {
					edges.push(getEdge(`e${dc.argoId}-${item.argoId}`, dc.argoId.toString(), item.argoId.toString()));
				}
			});
		}
	}

	return {
		nodes,
		edges,
	};
};

export const colorLegend = [
	{
		key: 'Succeeded',
		value: '#D3E4CD',
		order: 5,
	},
	{
		key: 'Failed',
		value: '#FDE9F2',
		order: 6,
	},
	{
		key: 'Waiting for approval',
		value: '#FFC85C',
		order: 2,
	},
	{
		key: 'In Progress',
		value: 'linear-gradient(to right,rgba(226, 194, 185, 0.5),rgba(226, 194, 185, 1),rgb(190, 148, 137))',
		order: 0,
	},
	{
		key: 'Destroyed/Not Provisioned',
		value: '#ddd',
		order: 5,
	},
];

export const getClassName = (status: string): string => {
	switch (status) {
		case ZSyncStatus.Initializing:
		case ZSyncStatus.InitializingApply:
			return '--initializing';
		case ZSyncStatus.RunningPlan:
		case ZSyncStatus.RunningDestroyPlan:
			return '--running';
		case ZSyncStatus.CalculatingCost:
			return '--running';
		case ZSyncStatus.WaitingForApproval:
			return '--pending';
		case ZSyncStatus.Provisioning:
		case ZSyncStatus.Destroying:
			return '--waiting';
		case ZSyncStatus.Provisioned:
			return '--successful';
		case ZSyncStatus.Destroyed:
		case ZSyncStatus.NotProvisioned:
			return '--destroyed';
		case ZSyncStatus.Skipped:
			return '--skipped';
		case ZSyncStatus.SkippedReconcile:
			return '--skipped-reconcile';
		case ZSyncStatus.PlanFailed:
		case ZSyncStatus.ValidationFailed:
		case ZSyncStatus.ApplyFailed:
		case ZSyncStatus.ProvisionFailed:
		case ZSyncStatus.DestroyFailed:
		case ZSyncStatus.OutOfSync:
		case ESyncStatus.OutOfSync:
			return '--failed';
		case ZSyncStatus.InSync:
		case ESyncStatus.Synced:
			return '--successful';
		default:
			return '--unknown';
	}
};

function getNodeIntersection(intersectionNode: any, targetNode: any) {
	// https://math.stackexchange.com/questions/1724792/an-algorithm-for-finding-the-intersection-point-between-a-center-of-vision-and-a
	const {
		width: intersectionNodeWidth,
		height: intersectionNodeHeight,
		positionAbsolute: intersectionNodePosition,
	} = intersectionNode;
	const targetPosition = targetNode.positionAbsolute;

	const w = intersectionNodeWidth / 2;
	const h = intersectionNodeHeight / 2;

	const x2 = intersectionNodePosition.x + w;
	const y2 = intersectionNodePosition.y + h;
	const x1 = targetPosition.x + w;
	const y1 = targetPosition.y + h;

	const xx1 = (x1 - x2) / (2 * w) - (y1 - y2) / (2 * h);
	const yy1 = (x1 - x2) / (2 * w) + (y1 - y2) / (2 * h);
	const a = 1 / (Math.abs(xx1) + Math.abs(yy1));
	const xx3 = a * xx1;
	const yy3 = a * yy1;
	const x = w * (xx3 + yy3) + x2;
	const y = h * (-xx3 + yy3) + y2;

	return { x, y };
}

// returns the position (top,right,bottom or right) passed node compared to the intersection point
function getEdgePosition(node: any, intersectionPoint: any) {
	const n = { ...node.positionAbsolute, ...node };
	const nx = Math.round(n.x);
	const ny = Math.round(n.y);
	const px = Math.round(intersectionPoint.x);
	const py = Math.round(intersectionPoint.y);

	if (px <= nx + 1) {
		return Position.Left;
	}
	if (px >= nx + n.width - 1) {
		return Position.Right;
	}
	if (py <= ny + 1) {
		return Position.Top;
	}
	if (py >= n.y + n.height - 1) {
		return Position.Bottom;
	}

	return Position.Top;
}

// returns the parameters (sx, sy, tx, ty, sourcePos, targetPos) you need to create an edge
export function getEdgeParams(source: any, target: any) {
	const sourceIntersectionPoint = getNodeIntersection(source, target);
	const targetIntersectionPoint = getNodeIntersection(target, source);

	const sourcePos = getEdgePosition(source, sourceIntersectionPoint);
	const targetPos = getEdgePosition(target, targetIntersectionPoint);

	return {
		sx: sourceIntersectionPoint.x,
		sy: sourceIntersectionPoint.y,
		tx: targetIntersectionPoint.x,
		ty: targetIntersectionPoint.y,
		sourcePos,
		targetPos,
	};
}

export function createNodesAndEdges() {
	const nodes = [];
	const edges = [];
	const center = { x: window.innerWidth / 2, y: window.innerHeight / 2 };

	nodes.push({ id: 'target', data: { label: 'Target' }, position: center });

	for (let i = 0; i < 8; i++) {
		const degrees = i * (360 / 8);
		const radians = degrees * (Math.PI / 180);
		const x = 250 * Math.cos(radians) + center.x;
		const y = 250 * Math.sin(radians) + center.y;

		nodes.push({ id: `${i}`, data: { label: 'Source' }, position: { x, y } });

		edges.push({
			id: `edge-${i}`,
			target: 'target',
			source: `${i}`,
			type: 'floating',
			markerEnd: {
				type: MarkerType.Arrow,
			},
		});
	}

	return { nodes, edges };
}

export function FloatingEdge({ id, source, target, markerEnd, style }: any) {
	const sourceNode = useStore(useCallback(store => store.nodeInternals.get(source), [source]));
	const targetNode = useStore(useCallback(store => store.nodeInternals.get(target), [target]));

	if (!sourceNode || !targetNode) {
		return null;
	}

	const { sx, sy, tx, ty, sourcePos, targetPos } = getEdgeParams(sourceNode, targetNode);

	const [edgePath] = getSimpleBezierPath({
		sourceX: sx,
		sourceY: sy,
		sourcePosition: sourcePos,
		targetPosition: targetPos,
		targetX: tx,
		targetY: ty,
	});

	return <path id={id} className="react-flow__edge-path" d={edgePath} markerEnd={markerEnd} style={style} />;
}
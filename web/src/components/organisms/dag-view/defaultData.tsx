import { ReactComponent as CloudIcon } from 'assets/images/icons/DAG-View/Cloud.svg';
import { ReactComponent as ComputeIcon } from 'assets/images/icons/DAG-View/compute.svg';
import { ReactComponent as EC2Icon } from 'assets/images/icons/DAG-View/ec2.svg';
import { ReactComponent as EKSIcon } from 'assets/images/icons/DAG-View/eks.svg';
import { ReactComponent as LayersIcon } from 'assets/images/icons/DAG-View/Layers.svg';
import { ReactComponent as LoadBalancerIcon } from 'assets/images/icons/DAG-View/load balancer.svg';
import { ReactComponent as NetworkingIcon } from 'assets/images/icons/DAG-View/Networking.svg';
import { ReactComponent as RDSIcon } from 'assets/images/icons/DAG-View/RDS.svg';
import { ReactComponent as S3Icon } from 'assets/images/icons/DAG-View/S3.svg';
import { ZSyncStatus } from 'models/argo.models';
import { DagNode } from 'models/dag.models';
import React, { SVGProps } from 'react';

export const svgProps: SVGProps<SVGSVGElement> = {
	x: '-12',
	y: '-12',
};

export const defaultData: DagNode = {
	name: 'Development',
	status: ZSyncStatus.OutOfSync,
	icon: <LayersIcon {...svgProps} />,
	children: [
		{
			name: 'Networking',
			status: ZSyncStatus.OutOfSync,
			icon: <NetworkingIcon {...svgProps} />,
			children: [
				{
					name: 'RDS',
					status: ZSyncStatus.InSync,
					icon: <RDSIcon {...svgProps} />,
				},
				{
					name: 'EKS',
					status: ZSyncStatus.OutOfSync,
					icon: <EKSIcon {...svgProps} />,
				},
				{
					name: 'EC2',
					status: ZSyncStatus.OutOfSync,
					icon: <EC2Icon {...svgProps} />,
					children: [
						{
							name: 'Lorem',
							status: ZSyncStatus.InSync,
							icon: <RDSIcon {...svgProps} />,
						},
						{
							name: 'Ipsum',
							status: ZSyncStatus.OutOfSync,
							icon: <EKSIcon {...svgProps} />,
						},
						{
							name: 'Dolor',
							status: ZSyncStatus.OutOfSync,
							icon: <ComputeIcon {...svgProps} />,
						},
					],
				},
			],
		},
		{
			name: 'S3',
			status: ZSyncStatus.InSync,
			icon: <S3Icon {...svgProps} />,
			children: [
				{
					name: 'Noam',
					status: ZSyncStatus.InSync,
					icon: <ComputeIcon {...svgProps} />,
				},
				{
					name: 'Noam',
					status: ZSyncStatus.InSync,
					icon: <ComputeIcon {...svgProps} />,
				},
				{
					name: 'Noam',
					status: ZSyncStatus.InSync,
					icon: <ComputeIcon {...svgProps} />,
					children: [
						{
							name: 'Noam',
							status: ZSyncStatus.InSync,
							icon: <ComputeIcon {...svgProps} />,
						},
						{
							name: 'Noam',
							status: ZSyncStatus.InSync,
							icon: <ComputeIcon {...svgProps} />,
						},
						{
							name: 'Noam',
							status: ZSyncStatus.InSync,
							icon: <ComputeIcon {...svgProps} />,
						},
					],
				},
			],
		},
		{
			name: 'Awan',
			status: ZSyncStatus.RunningPlan,
			icon: <CloudIcon {...svgProps} />,
			children: [
				{
					name: 'Enoch',
					status: ZSyncStatus.RunningPlan,
					icon: <CloudIcon {...svgProps} />,
				},
				{
					name: 'Enoch',
					status: ZSyncStatus.InSync,
					icon: <ComputeIcon {...svgProps} />,
					children: [
						{
							name: 'Enoch',
							status: ZSyncStatus.RunningPlan,
							icon: <CloudIcon {...svgProps} />,
						},
						{
							name: 'Enoch',
							status: ZSyncStatus.InSync,
							icon: <ComputeIcon {...svgProps} />,
						},
						{
							name: 'Enoch',
							status: ZSyncStatus.RunningPlan,
							icon: <LoadBalancerIcon {...svgProps} />,
						},
					],
				},
				{
					name: 'Enoch',
					status: ZSyncStatus.RunningPlan,
					icon: <LoadBalancerIcon {...svgProps} />,
				},
			],
		},
	],
};

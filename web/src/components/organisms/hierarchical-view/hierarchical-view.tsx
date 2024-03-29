import './styles.scss';

import { CostRenderer, currency } from 'components/molecules/cards/renderFunctions';
import React from 'react';
import { FC } from 'react';
import { Component } from 'models/entity.type';

type ResourceProps = {
	resource: Resource;
	depth: number;
};

type ResourcesProps = {
	resources: Resource[];
	depth: number;
};

type Resource = {
	name: string;
	hourlyCost: string;
	monthlyCost: string;
	subresources: Resource[];
	monthlyQuantity: any;
	unit: any;
	costComponents: [];
};

export type Props = {
	data: Component
}

export const HierarchicalView: FC<Props> = ({ data }) => {
	return (
		<>
			{data?.costResources?.length > 0 ? (
				<div className="hierarchy-container">
					<ul className="hierarchy-container_columns">
						<li>
							<span>Name</span>
						</li>
						<li>
							<span>Quantity</span>
						</li>
						<li>
							<span>Unit</span>
						</li>
						<li>
							<span>Monthly Cost</span>
						</li>
					</ul>
					<Hierarchy resources={data.costResources} depth={0} />
					<ul className="hierarchy-container_footer">
						<li>
							<span></span>
						</li>
						<li>
							<span></span>
						</li>
						<li>
							<span>Total Cost</span>
						</li>
						<li>
							<span>
								<CostRenderer data={data.estimatedCost} />
							</span>
						</li>
					</ul>
				</div>
			) : (
				'No resources were found!'
			)}
		</>
	);
};

const Hierarchy: FC<ResourcesProps> = ({ resources, depth = 0 }: ResourcesProps) => {
	return (
		<div className="hierarchy">
			{(resources || []).map((resource: any) => (
				<Node resource={resource} depth={depth} />
			))}
		</div>
	);
};

const Node: FC<ResourceProps> = ({ resource, depth }: ResourceProps) => {
	return (
		<>
			<div style={{ display: 'flex', paddingLeft: `${depth * 2}em` }}>
				{depth > 0 && <span className="node-hierarchy-l"></span>}
				<ul>
					<li>
						<span>{resource.name}</span>
					</li>
					<li>
						<span>{resource.monthlyQuantity}</span>
					</li>
					<li>
						<span>{resource.unit}</span>
					</li>
					<li>
						{!(resource?.subresources?.length > 0 || resource?.costComponents?.length > 0) && (
							<span>{`$${currency(Number(resource.monthlyCost))}`}</span>
						)}
					</li>
				</ul>
			</div>
			<Hierarchy resources={resource.subresources} depth={depth + 1} />
			<Hierarchy resources={resource.costComponents} depth={depth + 1} />
		</>
	);
};

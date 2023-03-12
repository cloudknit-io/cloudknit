import { Environment } from 'models/entity.type';
import { TreeReconcile } from 'pages/authorized/environments/helpers';
import { Reconciler } from 'pages/authorized/environments/Reconciler';
import React from 'react';
import { SmallText } from '../workflow-diagram/WorkflowDiagram';
import { colorLegend } from './tree-view.helper';

export type TreeViewControlProps = {
	environment?: Environment;
};

export const TreeViewControls: React.FC<TreeViewControlProps> = ({ environment }) => {
	return (
		<div className="node__title">
			<div></div>
			<div className="dag-controls">
				{environment && <Reconciler environment={environment} template={TreeReconcile} />}
			</div>
			<div className="color-legend-control">
				<div className="color-legend-control_status">
					<div>
						<label>Status:</label>
						{colorLegend
							.sort((a, b) => a.order - b.order)
							.map(color => (
								<span className="color-legend-control_value" key={color.key}>
									<label style={{ background: color.value }}></label>
									<label>{color.key}</label>
								</span>
							))}
					</div>
				</div>
				<SmallText data={'* Costs are monthly estimates calculated at the time of last reconciliation'}/>
			</div>
		</div>
	);
};

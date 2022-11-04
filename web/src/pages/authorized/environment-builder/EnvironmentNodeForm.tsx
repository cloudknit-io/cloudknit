import { YEnvironmentBuilder, YMetaData } from '../../../models/environment-builder';
import React from 'react';
import { getlabelInputPair, threeWayToggle } from 'helpers/environment-builder.helper';

type envSetter = {
	selectedNodeData: YEnvironmentBuilder;
	setSelectedNodeData: (val: any) => any;
};
interface Props {
	environmentSetter: envSetter;
	updateCallback: () => void;
}

export const EnvironmentNodeForm: React.FC<Props> = (props: Props) => {
	const { environmentSetter, updateCallback } = props;
	const { selectedNodeData, setSelectedNodeData } = environmentSetter;

	const setMetaData = () => {
		selectedNodeData.metadata = new YMetaData(
			selectedNodeData.spec.TeamName.value,
			selectedNodeData.spec.EnvName.value
		);
		setSelectedNodeData(Object.assign(new YEnvironmentBuilder(), selectedNodeData));
		updateCallback();
	};

	return (
		<>
			<div className="component-node-form-header">
				<h4 className="component-node-form-header__heading">Update Environment</h4>
			</div>
			<form
				key={`${selectedNodeData.id}`}
				noValidate
				onSubmit={e => {
					e.preventDefault();
				}}
				className="environment-builder-form">
				<div className="environment-builder-form__group flex-direction-column">
					<h5>Environment Info</h5>
					<div className="overflow-container">
						<div className="environment-builder-form__group">
							<div className="flex-basis-3 flex-basis-3__col">
								{getlabelInputPair(selectedNodeData.spec.TeamName, setMetaData)}
							</div>
							<div className="flex-basis-3 flex-basis-3__col">
								{getlabelInputPair(selectedNodeData.spec.EnvName, setMetaData)}
							</div>
							<div className="flex-basis-3 flex-basis-3__col">
								{threeWayToggle(selectedNodeData.spec.Teardown, setMetaData)}
							</div>
							<div className="flex-basis-3 flex-basis-3__col">
								{threeWayToggle(selectedNodeData.spec.AutoApprove, setMetaData)}
							</div>
						</div>
					</div>
				</div>
			</form>
		</>
	);
};

import { YComponent, YOutput, YSecret, YTag, YVariable, YVariableFile } from '../../../models/environment-builder';
import React, { useEffect } from 'react';
import { ReactComponent as Delete } from '../../../assets/images/icons/card-status/sync/delete.svg';
import { getDropDownList, getlabelInputPair, threeWayToggle } from 'helpers/environment-builder.helper';

type ComponentSetter = {
	selectedNodeData: YComponent;
	setSelectedNodeData: (val: YComponent) => any;
};

interface Props {
	componentSetter: ComponentSetter;
	updateCallback: () => void;
	deleteCallback: () => void;
}

export const ComponentNodeForm: React.FC<Props> = (props: Props) => {
	const { componentSetter, updateCallback, deleteCallback } = props;
	const { selectedNodeData, setSelectedNodeData } = componentSetter;

	const setNode = () => {
		setSelectedNodeData(Object.assign(new YComponent(), selectedNodeData));
		updateCallback();
	};

	const scrollerUpdate = (ev: any) => {
		const scroller = ev.currentTarget.closest('.update-scroll')?.querySelector('.scroll');
		setImmediate(() =>
			scroller?.scrollBy({
				behavior: 'smooth',
				top: scroller?.scrollHeight,
			})
		);
	};

	const getInlineVariableList = () => {
		return (
			selectedNodeData.inputList.size > 0 && (
				<div className="environment-builder-form__group flex-direction-column update-scroll">
					<h5>Inline Variables </h5>
					<div className="overflow-container flex-direction-column">
						<div className="environment-builder-form__group a-flex-end">
							<div className="flex-basis-3">
								{getDropDownList(
									new Set([...selectedNodeData.inputList.keys()]),
									(val: string, ev: React.MouseEvent<HTMLLIElement, MouseEvent>) => {
										if (!selectedNodeData.variables.has(val)) {
											scrollerUpdate(ev);
										}
										toggleInput(val);
									},
									(val: string) => `${selectedNodeData.variables.has(val) ? 'selected' : ''}`
								)}
							</div>
						</div>

						{selectedNodeData.variables.size > 0 && (
							<div className="overflow-container flex-direction-column scroll">
								<div className="environment-builder-form__group flex-direction-column">
									{[...selectedNodeData.variables.values()].map(e => (
										<div
											key={`${e.Name.value}`}
											className="environment-builder-form__group a-flex-end">
											<b>{e.Name.value}</b>
											<div className="flex-basis-3" key={`${e.Name.value}-${e.Value.key}`}>
												{getDropDownList(
													new Set(['Value', 'ValueFrom']),
													(val: 'Value' | 'ValueFrom') => {
														e.choosenDisposition = val;
														setNode();
													},
													(val: 'Value' | 'ValueFrom') =>
														e.choosenDisposition === val ? 'selected' : '',
													true,
													e.choosenDisposition,
													'Disposition'
												)}
											</div>
											<div
												className="flex-basis-3"
												key={`${e.Name.value}-${e.choosenDisposition}`}>
												{getlabelInputPair(e[e.choosenDisposition], setNode)}
											</div>
											<div>
												<button
													className="environment-builder-form__secret-tuple-remove"
													onClick={ev => {
														ev.stopPropagation();
														selectedNodeData.variables.delete(e.Name.value);
														setNode();
													}}>
													<Delete />
												</button>
											</div>
										</div>
									))}
								</div>
							</div>
						)}
					</div>
				</div>
			)
		);
	};

	const getOutputList = () => {
		return (
			selectedNodeData.outputList.size > 0 && (
				<div className="environment-builder-form__group flex-direction-column update-scroll">
					<h5>Outputs </h5>
					<div className="overflow-container flex-direction-column">
						<div className="environment-builder-form__group a-flex-end">
							<div className="flex-basis-3">
								{getDropDownList(
									selectedNodeData.outputList,
									(val: string, ev: any) => {
										if (!selectedNodeData.variables.has(val)) {
											scrollerUpdate(ev);
										}
										toggleOutput(val);
									},
									(val: string) => `${selectedNodeData.outputs.has(val) ? 'selected' : ''}`
								)}
							</div>
						</div>
					</div>
					<div className="environment-builder-form__secret--outputs scroll">
						{[...selectedNodeData.outputs.values()].map(e => (
							<div id={e.name} className="environment-builder-form__secret-tuple background">
								<label>{e.name}</label>
								<button
									className="environment-builder-form__secret-tuple-remove"
									onClick={ev => {
										ev.stopPropagation();
										toggleOutput(e.name);
									}}>
									<Delete />
								</button>
							</div>
						))}
					</div>
				</div>
			)
		);
	};

	const toggleOutput = (name: string) => {
		if (selectedNodeData.outputs.has(name)) {
			selectedNodeData.outputs.delete(name);
		} else {
			selectedNodeData.outputs.set(name, new YOutput(name));
		}
		setNode();
	};

	const toggleInput = (name: string) => {
		if (selectedNodeData.variables.has(name)) {
			selectedNodeData.variables.delete(name);
		} else {
			selectedNodeData.variables.set(name, new YVariable(name, '', ''));
		}

		setNode();
	};

	const getSecretList = () => {
		return (
			<div className="environment-builder-form__group flex-direction-column">
				<h5>
					Secrets{' '}
					<button
						className="environment-builder-form__secret-tuple-remove"
						onClick={() => {
							selectedNodeData.secrets.add(new YSecret());
							setNode();
						}}>
						+
					</button>
				</h5>
				<div className="overflow-container flex-direction-column">
					{[...selectedNodeData.secrets.values()].map((e, i) => (
						<div className="environment-builder-form__group a-flex-end" key={`${e.id}`}>
							<div className="flex-basis-3">{getlabelInputPair(e.Key, setNode)}</div>
							<div className="flex-basis-3">{getlabelInputPair(e.Name, setNode)}</div>
							<div className="flex-basis-3">
								{getDropDownList(
									new Set(['environment', 'team', 'org']),
									(val: string) => {
										e.Scope.set(val);
										setNode();
									},
									(val: string) => `${val === e.Scope.value ? 'selected' : ''}`,
									true,
									e.Scope.value,
									e.Scope.label
								)}
							</div>
							<div>
								<button
									className="environment-builder-form__secret-tuple-remove"
									onClick={ev => {
										ev.stopPropagation();
										selectedNodeData.secrets.delete(e);
										setNode();
									}}>
									<Delete />
								</button>
							</div>
						</div>
					))}
				</div>
			</div>
		);
	};

	const getVariablesFile = () => {
		if (!selectedNodeData.variablesFile) selectedNodeData.variablesFile = new YVariableFile('', '');
		return (
			<div className="environment-builder-form__group flex-direction-column">
				<h5>Variable Files </h5>
				<div className="overflow-container flex-direction-column">
					<div className="environment-builder-form__group a-flex-end" key={`${selectedNodeData.variablesFile.id}`}>
						<div className="flex-basis-3">{getlabelInputPair(selectedNodeData.variablesFile.Source, setNode)}</div>
						<div className="flex-basis-3">{getlabelInputPair(selectedNodeData.variablesFile.Path, setNode)}</div>
					</div>
				</div>
			</div>
		);
	};

	const getTagList = () => {
		return (
			<div className="environment-builder-form__group flex-direction-column">
				<h5>
					Tags{' '}
					<button
						className="environment-builder-form__secret-tuple-remove"
						onClick={() => {
							selectedNodeData.tags.add(new YTag());
							setNode();
						}}>
						+
					</button>
				</h5>
				<div className="overflow-container flex-direction-column">
					{[...selectedNodeData.tags.values()].map((e, i) => (
						<div key={`${e.id}`} className="environment-builder-form__group a-flex-end">
							<div className="flex-basis-3">{getlabelInputPair(e.Name, setNode)}</div>
							<div className="flex-basis-3">{getlabelInputPair(e.Value, setNode)}</div>
							<div>
								<button
									className="environment-builder-form__secret-tuple-remove"
									onClick={ev => {
										ev.stopPropagation();
										selectedNodeData.tags.delete(e);
										setNode();
									}}>
									<Delete />
								</button>
							</div>
						</div>
					))}
				</div>
			</div>
		);
	};

	return (
		<>
			<div className="component-node-form-header">
				<h4 className="component-node-form-header__heading">Update Component</h4>
			</div>
			<form
				key={`${selectedNodeData.id}`}
				noValidate
				onSubmit={e => {
					e.preventDefault();
				}}
				className="environment-builder-form">
				<div className="environment-builder-form__group flex-direction-column">
					<h5>Component Info</h5>
					<div className="overflow-container">
						<div className="environment-builder-form__group a-flex-end">
							<div className="flex-basis-3">{getlabelInputPair(selectedNodeData.Name, setNode)}</div>
							<div className="flex-basis-3">
								{getlabelInputPair(selectedNodeData.module.Name, setNode)}
							</div>
							<div className="flex-basis-3">
								{getlabelInputPair(selectedNodeData.module.Path, setNode)}
							</div>
							<div className="flex-basis-3">{threeWayToggle(selectedNodeData.AutoApprove, setNode)}</div>
							<div className="flex-basis-3">
								{threeWayToggle(selectedNodeData.DestroyProtection, setNode)}
							</div>
						</div>
					</div>
				</div>
				{getInlineVariableList()}
				{getVariablesFile()}
				{getOutputList()}
				{getSecretList()}
				{getTagList()}
			</form>
		</>
	);
};

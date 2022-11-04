import React, { useState } from 'react';
import './styles.scss';
import { ReactComponent as Avatar } from 'assets/images/icons/notification.svg';
import { ZEditor } from 'components/molecules/editor/Editor';
import { Context } from 'context/argo/ArgoUi';
import { StateFileService } from 'services/state-file/StateFile.service';

export const ZYamlValidator = () => {
	const [showValidator, toggleValidator] = useState<boolean>(false);
	const [validating, setValidating] = useState<boolean | null>(null);
	const [isValid, setValid] = useState<boolean | null>(null);
	const [yaml, setYaml] = useState<string>('');
	const notificationsApi = React.useContext(Context)?.notifications;

	const validateYaml = async () => {
		const val = await StateFileService.getInstance().validateYaml(yaml);
		setValid(val);
	};

	return (
		<>
			<div className="top-bar__validator">
				<button
					className="yml-icon"
					onClick={() => {
						toggleValidator(true);
					}}>
					YML
				</button>
			</div>
			{showValidator && (
				<div className="yamlvalidator-container">
					<div className="modal">
						<h4>Validate YAML</h4>
						<div
							className={`editor-container ${isValid === false ? 'invalid' : ''} ${
								isValid === true ? 'valid' : ''
							}`}>
							<ZEditor
								data={''}
								language={'yaml'}
								onChange={e => {
									setYaml(e || '');
								}}
                                height={'100%'}
							/>
						</div>
						<div className="btn-container">
							<button
								className={`btn green ${validating ? 'disabled' : ''}`}
								onClick={async () => {
									setValidating(true);
									notificationsApi?.show({
										content: 'Validating yaml',
										type: 0,
									});
									await validateYaml();
									setValidating(false);
								}}>
								{validating ? 'Validating...' : 'Validate'}
							</button>
							<button
								className={`btn red ${validating ? 'disabled' : ''}`}
								onClick={() => {
									toggleValidator(false);
								}}>
								Close
							</button>
						</div>
					</div>
				</div>
			)}
		</>
	);
};

import { ZEditor } from 'components/molecules/editor/Editor';
import { ZAsyncRenderer } from 'components/molecules/zasync-renderer/ZAsyncRenderer';
import React, { useEffect, useState } from 'react';
import { AuditService } from 'services/audit/audit.service';
import { StateFileService } from 'services/state-file/StateFile.service';
import { ReactComponent as Back } from 'assets/images/icons/chevron-right.svg';
import { ReactComponent as DownloadFile } from 'assets/images/icons/download.svg';
import { ReactComponent as EditFile } from 'assets/images/icons/edit.svg';
import './styles.scss';
import { Loader } from 'components/atoms/loader/Loader';
import { uniqueId } from 'lodash';
import { NotificationsApi, NotificationType } from 'components/argo-core';
import { Context } from 'context/argo/ArgoUi';

type StateFileViewProps = {
	teamName: string;
	environmentName: string;
	componentName: string;
	workflowName: string;
	ilRepo: string;
};

export const StateFileView: React.FC<StateFileViewProps> = ({
	teamName,
	environmentName,
	componentName,
	workflowName,
	ilRepo,
}: StateFileViewProps) => {
	const service = StateFileService.getInstance();
	const [editView, setEditView] = useState<boolean>(false);
	const [componentList, setComponentList] = useState<string[]>([]);
	const [selectedComponents, setSelectedComponents] = useState<Set<string>>(new Set<string>());
	const [currentData, setCurrentData] = useState<any>(null);
	const [fetch, setFetch] = useState<boolean>(false);
	const [rendererKey, setRendererKey] = useState<any>(uniqueId(workflowName));
	const [fetchMessage, setFetchMessage] = useState<string>('Fetching resources...');
	const contextApi = React.useContext(Context);

	useEffect(() => {
		AuditService.getInstance()
		.fetchStateFile(teamName, environmentName, componentName, rendererKey)
		.then(({ data }) => {
			setCurrentData(data);
			return data;
		});

		return () => {
			setCurrentData(null);
		}
	}, [teamName, environmentName, componentName]);

	useEffect(() => {
		setRendererKey(uniqueId(workflowName));
	}, [workflowName]);

	useEffect(() => {
		if (!editView || !workflowName) {
			return;
		}
		setFetch(true);
		const repo = ilRepo.includes("::") ? ilRepo.split("::")[1] : ilRepo.split(':').slice(-1)[0];
		service
			.fetchStateListComponents(teamName, environmentName, componentName, repo)
			.then((response: any) => {
				setFetch(false);
				if (!response?.resources) {
					(contextApi?.notifications as NotificationsApi).show({
						content: 'No resources were found.',
						type: NotificationType.Error,
					});
					return;
				}
				setComponentList(response.resources);
			})
			.catch(err => {
				setFetch(false);
				(contextApi?.notifications as NotificationsApi).show({
					content: err.message || 'There was an error fetching resources',
					type: NotificationType.Error,
				});
			});
		return () => {
			setComponentList([]);
		}
	}, [workflowName, editView]);

	const downloadData = () => {
		if (!currentData) return;
		const blob = new Blob([currentData], { type: 'text' });
		const downloadLink = document.createElement('a');
		downloadLink.href = window.URL.createObjectURL(blob);
		downloadLink.download = 'tfState.json';
		document.body.appendChild(downloadLink);
		downloadLink.click();
		document.body.removeChild(downloadLink);
	};

	const toggleSelectedComponents = (component: string, value: boolean) => {
		if (value) {
			selectedComponents.add(component);
		} else {
			selectedComponents.delete(component);
		}
		setSelectedComponents(new Set(selectedComponents));
	};

	const getLoader = (message: string) => {
		return (
			<div style={{ display: 'flex', justifyContent: 'center' }}>
				<Loader title={message} height={20} width={20} />
				<span>{message}</span>
			</div>
		);
	};

	const getStateComponentEditor = () => {
		if (fetch) {
			return getLoader(fetchMessage);
		}
		return (
			<>
				<header className="state-list-header">
					<span>
						<input
							type="checkbox"
							checked={selectedComponents.size === componentList.length}
							onChange={ev => {
								if (ev.target.checked) {
									setSelectedComponents(new Set(componentList));
								} else {
									setSelectedComponents(new Set());
								}
							}}
						/>
					</span>{' '}
					<span className="header-text">Components</span>
					<span>
						<button
							disabled={selectedComponents.size === 0}
							onClick={ev => {
								const proceed = window.confirm('Are you sure you want to delete the resources from terraform state?');
								if (!proceed) {
									return;
								}
								setFetchMessage(`Deleting selected resources...`);
								(contextApi?.notifications as NotificationsApi).show({
									content: 'Deleting resources',
									type: NotificationType.Warning,
								});
								setFetch(true);
								setFetchMessage('Deleting resources...');
								const repo = ilRepo.includes("::") ? ilRepo.split("::")[1] : ilRepo.split(':').slice(-1)[0];
								service
									.deleteComponents(teamName, environmentName, componentName, repo, [
										...selectedComponents.values(),
									])
									.then(({ resources }: any) => {
										setFetch(false);
										if (!resources) {
											(contextApi?.notifications as NotificationsApi).show({
												content: 'No resources were found.',
												type: NotificationType.Error,
											});
											return;
										}
										(contextApi?.notifications as NotificationsApi).show({
											content: 'Deletion Successful.',
											type: NotificationType.Success,
										});
										setComponentList(resources);
										setSelectedComponents(new Set());
										setRendererKey(uniqueId(workflowName));
									})
									.catch(err => {
										setFetch(false);
										(contextApi?.notifications as NotificationsApi).show({
											content: err.message || 'There was an error deleting resources',
											type: NotificationType.Error,
										});
									});
							}}>
							Delete
						</button>
					</span>
				</header>
				<ul className="state-list-components">
					{componentList.map(c => (
						<li
							key={c}
							className={`state-list-components__row ${selectedComponents.has(c) ? 'selected' : ''}`}>
							<span>
								<input
									type="checkbox"
									checked={selectedComponents.has(c)}
									onChange={ev => {
										ev.stopPropagation();
										toggleSelectedComponents(c, ev.currentTarget.checked);
									}}
								/>
							</span>{' '}
							<span>{c}</span>
						</li>
					))}
				</ul>
			</>
		);
	};

	const getEditorGroup = () => {
		return (
			<div className="zeditor-options">
				<button
					title="download as a file"
					className="download-button"
					onClick={ev => {
						ev.stopPropagation();
						downloadData();
					}}>
					<DownloadFile title="download as a file" />
				</button>

				<button
					title="Edit State File"
					className="download-button"
					onClick={ev => {
						ev.stopPropagation();
						setEditView(true);
					}}>
					<EditFile title="edit file" />
				</button>
			</div>
		);
	};

	return (
		<div style={{ position: 'relative' }}>
			{editView ? (
				<>
					<button
						className="hide-editor"
						onClick={ev => {
							setEditView(false);
						}}>
						<Back />
						Back
					</button>
					{getStateComponentEditor()}
				</>
			) : (
				<>
					{currentData !== null && getEditorGroup()}
					<ZEditor options={{
						readOnly: true
					}} data={currentData === null ? 'Loading...' : currentData}/>
				</>
			)}
		</div>
	);
};

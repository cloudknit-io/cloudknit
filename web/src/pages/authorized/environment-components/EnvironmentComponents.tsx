import { ZLoaderCover } from 'components/atoms/loader/LoaderCover';
import { EnvironmentComponentCards } from 'components/molecules/cards/EnvironmentComponentCards';
import { ZSidePanel } from 'components/molecules/side-panel/SidePanel';
import { ZTablControl } from 'components/molecules/tab-control/TabControl';
import { AuditView } from 'components/organisms/audit_view/AuditView';
import { ErrorView } from 'components/organisms/error-view/ErrorView';
import { streamMapperWF } from 'helpers/streamMapper';
import { useApi } from 'hooks/use-api/useApi';
import { ApplicationWatchEvent, ZComponentSyncStatus, ZSyncStatus } from 'models/argo.models';
import { EntityStore } from 'models/entity.store';
import { CompAuditData, Component, EnvAuditData, Environment } from 'models/entity.type';
import { eventErrorColumns } from 'models/error.model';
import { LocalStorageKey } from 'models/localStorage';
import { EnvironmentItem } from 'models/projects.models';
import moment from 'moment';
import { ConfigWorkflowView } from 'pages/authorized/environment-components/config-workflow-view/ConfigWorkflowView';
import { auditColumns, ConfigParamsSet, getWorkflowLogs } from 'pages/authorized/environment-components/helpers';
import React, { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useParams } from 'react-router-dom';
import { ReactFlowProvider } from 'reactflow';
import { Subscription } from 'rxjs';
import { ArgoStreamService } from 'services/argo/ArgoStream.service';
import { ArgoWorkflowsService } from 'services/argo/ArgoWorkflows.service';
import { AuditService } from 'services/audit/audit.service';
import { subscriberWF } from 'utils/apiClient/EventClient';
import { TreeView } from '../../../components/organisms/treeview/TreeViewNew';
import { usePageHeader } from '../contexts/EnvironmentHeaderContext';
import { getCheckBoxFilters, renderSyncStatusItems } from '../environments/helpers';

const envTabs = [
	{
		id: 'Audit',
		name: 'Audit',
		show: (show: () => boolean) => true,
	},
	{
		id: 'Errors',
		name: 'Errors',
		show: (show: () => boolean) => show(),
	},
];

export const EnvironmentComponents: React.FC = () => {
	// Migrating to API
	// Get all components and start streaming
	// End Streaming whenever we change the environment or URL changes [ this would decrease load on web ]
	// Alter the data store
	const entityStore = useMemo(() => EntityStore.getInstance(), []);
	const { fetch: fetchWorkflowData } = useApi(ArgoWorkflowsService.getConfigWorkflow);
	const [syncStatusFilter, setSyncStatusFilter] = useState<Set<ZSyncStatus>>(new Set<ZSyncStatus>());
	const [filterDropDownOpen, toggleFilterDropDown] = useState(false);
	const [checkBoxFilters, setCheckBoxFilters] = useState<JSX.Element>(<></>);
	const [filterItems, setFilterItems] = useState<Array<() => JSX.Element>>([]);
	const { projectId, environmentName } = useParams<any>();
	const [isLoadingWorkflow, setIsLoadingWorkflow] = useState<boolean>();
	const showAll = environmentName === 'all' && projectId === 'all';
	const [environment, setEnvironment] = useState<Environment>();
	const { pageHeaderObservable, breadcrumbObservable } = usePageHeader();
	const [query, setQuery] = useState<string>('');
	const [showSidePanel, setShowSidePanel] = useState<boolean>(false);
	const [selectedConfig, setSelectedConfig] = useState<Component>();
	const [loading, setLoading] = useState<boolean>(true);
	const [components, setComponents] = useState<Component[]>([]);
	const [workflowData, setWorkflowData] = useState<any>();
	const [streamData2, setStreamData2] = useState<ApplicationWatchEvent | null>(null);
	const [logs, setLogs] = useState<string | null>(null);
	const [plans, setPlans] = useState<string | null>(null);
	const [workflowId, setWorkflowId] = useState<string>('');
	const [viewType, setViewType] = useState<string>(showAll ? '' : 'DAG');
	const [isEnvironmentNodeSelected, setEnvironmentNodeSelected] = useState<boolean>(false);
	const [envAuditList, setEnvAuditList] = useState<EnvAuditData[]>([]);
	const [compAuditList, setCompAuditList] = useState<CompAuditData[]>([]);
	const componentArrayRef = useRef<Component[]>([]);
	const selectedComponentRef = useRef<Component>();
	const compAuditRef = useRef<CompAuditData[]>([]);
	const envAuditRef = useRef<EnvAuditData[]>([]);
	const workflowIdRef = useRef<string>();
	const resetRefs = useCallback(() => {
		setLoading(true);
		componentArrayRef.current = [];
		selectedComponentRef.current = undefined;
		compAuditRef.current = [];
		envAuditRef.current = [];
		workflowIdRef.current = '';
		setComponents(componentArrayRef.current);
		setSelectedConfig(selectedComponentRef.current);
		setCompAuditList(compAuditRef.current);
		setEnvAuditList(envAuditRef.current);
		setWorkflowId(workflowIdRef.current);
	}, []);

	const breadcrumbItems = [
		{
			path: '/dashboard',
			name: 'All Environments',
			active: false,
		},
		{
			path: `/${projectId}`,
			name: projectId,
			active: false,
		},
		{
			path: `/${projectId}/${environmentName}`,
			name: environmentName,
			active: true,
		},
	];

	useEffect(() => {
		if (showAll) {
			setViewType('');
		} else {
			setViewType('DAG');
		}
	}, [showAll]);

	useEffect(() => {
		setFilterItems([
			renderSyncStatusItems
				.bind(null, ZComponentSyncStatus, syncStatusFilter, setSyncStatusFilter, 'Component Status')
				.bind(null, (status: string) => components.filter(e => e.status === status).length),
		]);
	}, [components, syncStatusFilter]);

	useEffect(() => {
		setCheckBoxFilters(getCheckBoxFilters(filterDropDownOpen, setToggleFilterDropDownValue, filterItems));
	}, [filterItems, filterDropDownOpen, filterDropDownOpen]);

	useEffect(() => {
		const headerTabs = showAll
			? []
			: entityStore.getAllEnvironmentsByTeamName(projectId).filter(e => e.dag.length > 0).map(environment => {
					const name: string = environment.name;
					return {
						active: environmentName === environment.name,
						name: name.charAt(0).toUpperCase() + name.slice(1),
						path: `/${projectId}/${environment.name}`,
					};
			  });

		pageHeaderObservable.next({
			breadcrumbs: breadcrumbItems,
			headerTabs,
			pageName: 'Components',
			filterTitle: 'Filter by environment',
			onSearch: setQuery,
			onViewChange: setViewType,
			buttonText: 'New Configuration',
			initialView: viewType,
			checkBoxFilters: viewType !== 'DAG' ? checkBoxFilters : null,
		});
		breadcrumbObservable.next({
			[LocalStorageKey.ENVIRONMENTS]: headerTabs,
		});
	}, [checkBoxFilters, viewType, environment]);

	useEffect(() => {
		if (!projectId || !environmentName) return;

		resetRefs();

		const subs: any[] = [];
		subs.push(
			entityStore.emitter.subscribe(data => {
				if (showAll) {
					entityStore.AllDataFetched && setComponents([...data.components]);
					componentArrayRef.current = data.components;
					setLoading(false);
				} else {
					if (data.environments.length === 0) return;
					const env = entityStore.getEnvironmentByName(projectId, environmentName);
					if (!env) return;
					setEnvironment(env);
				}
			})
		);

		subs.push(
			subscriberWF.subscribe((response: any) => {
				setStreamData2(response);
			})
		);

		if (showAll) {
			Promise.resolve(entityStore.getTeams(true));
		}

		return () => {
			subs.forEach(sub => sub.unsubscribe());
		};
	}, [projectId, environmentName]);

	useEffect(() => {
		if (!environment?.id || showAll) return;
		let subEnvAudit: Subscription | null = null;
		AuditService.getInstance()
			.getEnvironment(environment.id, environment.teamId)
			.then(data => {
				if (Array.isArray(data)) {
					envAuditRef.current = data;
					setEnvAuditList(data);
					subEnvAudit = entityStore
						.setEnvironmentAuditLister(environment.id)
						.subscribe((response: EnvAuditData) => {
							const idx = envAuditRef.current.findIndex(e => e.reconcileId === response.reconcileId);
							if (idx !== -1) {
								envAuditRef.current[idx] = response;
							} else {
								envAuditRef.current.push(response);
							}
							setEnvAuditList([...envAuditRef.current]);
						});
				}
			});

		const sub = entityStore.emitterComp.subscribe((components: Component[]) => {
			if (components.length === 0 || components[0].envId !== environment.id) return;
			componentArrayRef.current = components;
			setComponents(components);
			resetSelectedConfig(components);
			setLoading(false);
		});
		Promise.resolve(entityStore.getComponents(environment.teamId, environment.id));

		return () => {
			sub.unsubscribe();
			subEnvAudit && entityStore.removeEnvironmentAuditLister(environment.id, subEnvAudit);
		};
	}, [environment?.id]);

	useEffect(() => {
		if (!selectedConfig) return;
		let subCompAudit: Subscription | null = null;
		AuditService.getInstance()
			.getComponent(selectedConfig.id, selectedConfig.envId, selectedConfig.teamId)
			.then(data => {
				if (Array.isArray(data)) {
					compAuditRef.current = data;
					setCompAuditList(data);
					subCompAudit = entityStore
						.setComponentAuditLister(selectedConfig.id)
						.subscribe((response: CompAuditData) => {
							if (response.compId !== selectedConfig.id) return;
							const idx = compAuditRef.current.findIndex(e => e.reconcileId === response.reconcileId);
							if (idx !== -1) {
								compAuditRef.current[idx] = response;
							} else {
								compAuditRef.current.push(response);
							}
							setCompAuditList([...compAuditRef.current]);
						});
				}
			});
		return () => {
			subCompAudit && entityStore.removeComponentAuditLister(selectedConfig.id, subCompAudit);
		};
	}, [selectedConfig]);

	useEffect(() => {
		const sub = entityStore.emitterCompAudit.subscribe((auditData: CompAuditData) => {
			const idx = componentArrayRef.current.findIndex(e => e.id === auditData.compId);
			if (idx !== -1) {
				componentArrayRef.current[idx].lastAuditStatus = auditData.status;
			}
			setComponents([...componentArrayRef.current]);
		});

		return () => sub.unsubscribe();
	}, []);

	useEffect(() => {
		const newWf: any = streamMapperWF(streamData2);
		if (newWf && workflowData) {
			setWorkflowData({
				...workflowData,
				status: newWf?.object?.status,
			});
		}
	}, [streamData2]);

	const syncStatusMatch = (item: Component): boolean => {
		return syncStatusFilter.has(item.status as ZSyncStatus);
	};

	const setToggleFilterDropDownValue = (val: boolean) => {
		toggleFilterDropDown(val);
	};

	const onNodeClick = (nodeId: string): void => {
		if (nodeId === environment?.argoId) {
			setEnvironmentNodeSelected(true);
			setShowSidePanel(true);
			return;
		} else {
			setEnvironmentNodeSelected(false);
			const selectedConfig = componentArrayRef.current.find(c => c.argoId === nodeId);
			if (selectedConfig) {
				let _workflowId = selectedConfig.lastWorkflowRunId;
				if (!_workflowId) {
					_workflowId = 'initializing';
				}
				selectedComponentRef.current = selectedConfig;
				setSelectedConfig(selectedConfig);
				setShowSidePanel(true);
				if (workflowIdRef.current !== _workflowId) {
					workflowIdRef.current = _workflowId;
					setWorkflowId(_workflowId);
					setLogs(null);
					setPlans(null);
					setIsLoadingWorkflow(false);
					setWorkflowData(null);
				}
			}
		}
	};

	const getFilteredData = (): Component[] => {
		let filteredItems = [...componentArrayRef.current];
		if (syncStatusFilter.size > 0) {
			filteredItems = [...filteredItems.filter(syncStatusMatch)];
		}

		return filteredItems.filter(item => {
			return item.argoId.toLowerCase().includes(query);
		});
	};

	const resetSelectedConfig = (components: Component[]) => {
		const selectedConf = components.find((itm: any) => itm.id === selectedComponentRef.current?.id);
		if (selectedConf) {
			selectedComponentRef.current = selectedConf;
			setSelectedConfig(selectedConf);
			if (selectedConf.lastWorkflowRunId !== workflowIdRef.current) {
				workflowIdRef.current = selectedConf.lastWorkflowRunId;
				setWorkflowId(selectedConf.lastWorkflowRunId);
			}
		}
	};

	useEffect(() => {
		if (workflowId) {
			if (workflowId === 'initializing') {
				setLogs(null);
				setPlans(null);
				setIsLoadingWorkflow(false);
				setWorkflowData(null);
			} else {
				getWorkflowData(workflowId, selectedConfig?.argoId || '');
			}
		}
	}, [workflowId]);

	const getWorkflowData = (workflowId: string, configId: string) => {
		setLogs(null);
		setPlans(null);
		setIsLoadingWorkflow(true);
		fetchWorkflowData({
			projectId: projectId,
			environmentId: environmentName,
			configId: configId,
			workflowId: workflowId,
		}).then(({ data }) => {
			setIsLoadingWorkflow(false);
			setWorkflowData(data);
			const configParamsSet: ConfigParamsSet = {
				projectId,
				environmentId: environmentName,
				configId: configId,
				workflowId: workflowId,
			};
			ArgoStreamService.streamWF(configParamsSet);
			getWorkflowLogs(configParamsSet, fetchWorkflowData, setPlans, setLogs);
		});
		return;
	};

	const renderItems = (): any => {
		switch (viewType) {
			case 'DAG':
				return components.length > 0 ? (
					<ReactFlowProvider>
						<TreeView environmentItem={environment} onNodeClick={onNodeClick} />
					</ReactFlowProvider>
				) : (
					<></>
				);
			default:
				return (
					<EnvironmentComponentCards
						showAll={showAll}
						components={getFilteredData() || []}
						projectId={projectId}
						env={environment}
						selectedConfig={selectedConfig}
						workflowPhase={workflowData?.status?.phase}
						onClick={(config: Component): void => {
							onNodeClick(config.name);
						}}
					/>
				);
		}
	};

	return (
		<div className="zlifecycle-page">
			<ZLoaderCover loading={loading}>
				<section className="dashboard-content container">
					{renderItems()}
					<ZSidePanel isShown={showSidePanel} onClose={(): void => setShowSidePanel(false)}>
						{isEnvironmentNodeSelected && (
							<div className="ztab-control">
								<div className="ztab-control__tabs">
									<ZTablControl
										className="container__tabs-control"
										selected={[].length ? 'Errors' : 'Audit'}
										tabs={envTabs.filter(t => t.show(() => Boolean(environment?.errorMessage?.length)))}>
										<div id="Errors">
											{environment?.errorMessage && (
												<ErrorView
													columns={eventErrorColumns}
													dataRows={environment.errorMessage.map(e => ({
														team: EntityStore.getInstance().getTeam(environment.teamId)
															?.name,
														environment: environment.name,
														message: e,
														timestamp: moment(
															environment.lastReconcileDatetime.toString(),
															moment.ISO_8601
														).fromNow(),
													}))}
												/>
											)}
										</div>
										<div id="Audit">
											<AuditView auditData={envAuditList} auditColumns={auditColumns} />
										</div>
									</ZTablControl>
								</div>
							</div>
						)}
						<ZLoaderCover loading={isLoadingWorkflow}>
							{selectedConfig && !isEnvironmentNodeSelected && (
								<ConfigWorkflowView
									projectId={projectId}
									environmentId={environmentName}
									config={selectedConfig}
									logs={logs}
									plans={plans}
									workflowData={workflowData}
									auditData={compAuditList}
								/>
							)}
						</ZLoaderCover>
					</ZSidePanel>
				</section>
			</ZLoaderCover>
		</div>
	);
};

import { NotificationType } from 'components/argo-core';
import { ZLoaderCover } from 'components/atoms/loader/LoaderCover';
import { ZTable } from 'components/atoms/table/Table';
import { environmentName } from 'components/molecules/cards/EnvironmentCards';
import { EnvironmentComponentCards } from 'components/molecules/cards/EnvironmentComponentCards';
import { ZSidePanel } from 'components/molecules/side-panel/SidePanel';
import { AuditView } from 'components/organisms/audit_view/AuditView';
import { TreeComponent } from 'components/organisms/tree-view/TreeComponent';
import { Context } from 'context/argo/ArgoUi';
import { streamMapper, streamMapperWF } from 'helpers/streamMapper';
import { useApi } from 'hooks/use-api/useApi';
import { ApplicationWatchEvent, HealthStatusCode, ZComponentSyncStatus, ZSyncStatus } from 'models/argo.models';
import { LocalStorageKey } from 'models/localStorage';
import {
	EnvironmentComponentItem,
	EnvironmentComponentsList,
	EnvironmentItem,
	EnvironmentsList,
} from 'models/projects.models';
import { ConfigWorkflowView } from 'pages/authorized/environment-components/config-workflow-view/ConfigWorkflowView';
import {
	auditColumns,
	ConfigParamsSet,
	configTableColumns,
	getSeparatedConfigId,
	getWorkflowLogs,
} from 'pages/authorized/environment-components/helpers';
import React, { useContext, useEffect, useMemo, useRef, useState } from 'react';
import { useParams } from 'react-router-dom';
import { Subscription } from 'rxjs';
import { debounceTime } from 'rxjs/operators';
import { ArgoComponentsService } from 'services/argo/ArgoComponents.service';
import { ArgoEnvironmentsService } from 'services/argo/ArgoEnvironments.service';
import { ArgoMapper } from 'services/argo/ArgoMapper';
import { ArgoStreamService } from 'services/argo/ArgoStream.service';
import { ArgoWorkflowsService } from 'services/argo/ArgoWorkflows.service';
import { AuditService } from 'services/audit/audit.service';
import { subscriber, subscriberWF } from 'utils/apiClient/EventClient';
import { usePageHeader } from '../contexts/EnvironmentHeaderContext';
import { getCheckBoxFilters, renderSyncStatusItems } from '../environments/helpers';
import { toast } from 'react-toastify';
import { Loader } from 'components/atoms/loader/Loader';
import { ConfigWorkflowViewApplication } from './config-workflow-view/ConfigWorkflowViewApplication';
import { ErrorView } from 'components/organisms/error-view/ErrorView';
import { ZTablControl } from 'components/molecules/tab-control/TabControl';
import { ErrorStateService } from 'services/error/error-state.service';
import { eventErrorColumns } from 'models/error.model';
import { CostingService } from 'services/costing/costing.service';
import { CompAuditData, Component, EntityStore, EnvAuditData, Environment } from 'models/entity.store';

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
	const { fetch } = useApi(ArgoComponentsService.getComponents);
	const fetchEnvironments = useApi(ArgoEnvironmentsService.getEnvironments).fetch;
	const { fetch: fetchWorkflowData } = useApi(ArgoWorkflowsService.getConfigWorkflow);
	const [syncStatusFilter, setSyncStatusFilter] = useState<Set<ZSyncStatus>>(new Set<ZSyncStatus>());
	const [filterDropDownOpen, toggleFilterDropDown] = useState(false);
	const [checkBoxFilters, setCheckBoxFilters] = useState<JSX.Element>(<></>);
	const [filterItems, setFilterItems] = useState<Array<() => JSX.Element>>([]);
	const { projectId, environmentName } = useParams<any>();
	const [isLoadingWorkflow, setIsLoadingWorkflow] = useState<boolean>();
	const showAll = environmentName === 'all' && projectId === 'all';
	const notificationManager = React.useContext(Context)?.notifications;
	const [environment, setEnvironment] = useState<Environment>();
	const { pageHeaderObservable, breadcrumbObservable } = usePageHeader();
	const [query, setQuery] = useState<string>('');
	const [showSidePanel, setShowSidePanel] = useState<boolean>(false);
	const [selectedConfig, setSelectedConfig] = useState<Component>();
	const [loading, setLoading] = useState<boolean>(true);
	const [components, setComponents] = useState<Component[]>([]);
	const [workflowData, setWorkflowData] = useState<any>();
	const [streamData, setStreamData] = useState<ApplicationWatchEvent | null>(null);
	const [streamData2, setStreamData2] = useState<ApplicationWatchEvent | null>(null);
	const [logs, setLogs] = useState<string | null>(null);
	const [plans, setPlans] = useState<string | null>(null);
	const [workflowId, setWorkflowId] = useState<string>('');
	const [viewType, setViewType] = useState<string>(showAll ? '' : 'DAG');
	const [isEnvironmentNodeSelected, setEnvironmentNodeSelected] = useState<boolean>(false);
	const [envErrors, setEnvErrors] = useState<any[]>();
	const [envAuditList, setEnvAuditList] = useState<EnvAuditData[]>([]);
	const [compAuditList, setCompAuditList] = useState<CompAuditData[]>([]);
	const componentArrayRef = useRef<Component[]>([]);
	const selectedComponentRef = useRef<Component>();
	const compAuditRef = useRef<CompAuditData[]>([]);
	const envAuditRef = useRef<EnvAuditData[]>([]);
	const workflowIdRef = useRef<string>();
	const ctx = useContext(Context);

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
				.bind(null, (status: string) => components.filter(e => 'e.componentStatus' === status).length),
		]);
	}, [components, syncStatusFilter]);

	useEffect(() => {
		setCheckBoxFilters(getCheckBoxFilters(filterDropDownOpen, setToggleFilterDropDownValue, filterItems));
	}, [filterItems, filterDropDownOpen, filterDropDownOpen]);

	useEffect(() => {
		const headerTabs = entityStore.getAllEnvironmentsByTeamName(projectId).map(environment => {
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
			[LocalStorageKey.ENVIRONMENTS]: !showAll ? headerTabs : [],
		});
	}, [checkBoxFilters, viewType, environment]);

	useEffect(() => {
		if (!projectId || !environmentName) return;
		const subs: any[] = [];
		subs.push(
			entityStore.emitter.subscribe(async data => {
				if (data.environments.length === 0) return;
				const env = entityStore.getEnvironmentByName(projectId, environmentName);
				if (!env) return;
				setEnvironment(env);
			})
		);

		subs.push(
			subscriberWF.subscribe((response: any) => {
				setStreamData2(response);
			})
		);

		return () => {
			subs.forEach(sub => sub.unsubscribe());
		};
	}, [projectId, environmentName]);

	useEffect(() => {
		if (!environment?.id) return;
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
	}, [])

	// useEffect(() => {
	// 	const newEnvironments = streamMapper<EnvironmentItem>(
	// 		streamData,
	// 		environments,
	// 		ArgoMapper.parseEnvironment,
	// 		'environment',
	// 		{
	// 			projectId,
	// 		}
	// 	);
	// 	checkForFailedEnvironments(newEnvironments);
	// 	setEnvironments(newEnvironments);

	// 	const newComponents = streamMapper<EnvironmentComponentItem>(
	// 		streamData,
	// 		components,
	// 		ArgoMapper.parseComponent,
	// 		'config',
	// 		{
	// 			projectId,
	// 			environmentId,
	// 		}
	// 	);
	// 	setComponents(newComponents);
	// 	componentArrayRef.current = newComponents;
	// }, [streamData]);

	// useEffect(() => {
	// 	if (environments.length === 0 || !environmentId) return;
	// 	const env = environments.find(e => e.id === environmentId);
	// 	if (!env) return;
	// 	const sub = ctx?.failedEnvironments.subscribe(res => {
	// 		const errors = [...res.values()].filter(e => e.labels.env_name === env.labels?.env_name);
	// 		console.log(errors);
	// 	});

	// 	return () => sub?.unsubscribe();
	// }, [environmentId, environments]);

	// useEffect(() => {
	// 	const $subscription: Subscription = subscriber.subscribe((response: any) => {
	// 		setStreamData(response);
	// 	});
	// 	const $subscription2: Subscription = subscriberWF.subscribe((response: any) => {
	// 		setStreamData2(response);
	// 	});

	// 	return (): void => {
	// 		$subscription.unsubscribe();
	// 		$subscription2.unsubscribe();
	// 	};
	// }, []);

	useEffect(() => {
		const newWf: any = streamMapperWF(streamData2);
		if (newWf && workflowData) {
			setWorkflowData({
				...workflowData,
				status: newWf?.object?.status,
			});
		}
	}, [streamData2]);

	// useEffect(() => {
	// 	let subs: Subscription[] = [];
	// 	setLoading(true);
	// 	const fetchs = [
	// 		CostingService.getInstance().getEnvironmentInfo(projectId, environmentId.replace(projectId + '-', '')),
	// 		fetch(projectId, environmentId),
	// 	];
	// 	Promise.allSettled(fetchs).then(resps => {
	// 		const { data0, data1 } = resolveFetchs(resps);
	// 		if (data1) {
	// 			componentArrayRef.current = data1;
	// 			subs = setUpComponentStreams(data1);
	// 		}
	// 		if (data0) {
	// 			const newItems = componentArrayRef.current.map(nc => {
	// 				const newItem = data0.find((c: any) => c.id === nc.displayValue);
	// 				if (newItem) {
	// 					nc.componentCost = newItem.estimatedCost;
	// 					nc.componentStatus = newItem.status;
	// 					nc.costResources = newItem.costResources;
	// 					nc.syncFinishedAt = newItem.lastReconcileDatetime;
	// 					nc.isDestroy = newItem.isDestroyed;
	// 				}
	// 				return nc;
	// 			});
	// 			componentArrayRef.current = newItems;
	// 		}
	// 		setComponents(componentArrayRef.current);
	// 		setLoading(false);
	// 	});
	// 	return () => {
	// 		subs.forEach(s => s.unsubscribe());
	// 	};
	// }, ['projectId', environmentId]);

	// useEffect(() => {
	// 	if (showSidePanel === false) {
	// 		setLogs(null);
	// 		setPlans(null);
	// 		setIsLoadingWorkflow(false);
	// 		setWorkflowData(null);
	// 		setWorkflowId('');
	// 		workflowIdRef.current = '';
	// 		setSelectedConfig(undefined);
	// 	}
	// }, [showSidePanel]);

	// useEffect(() => {
	// 	resetSelectedConfig(componentArrayRef.current);
	// }, [components]);

	// const setUpComponentStreams = (newComponents: EnvironmentComponentsList) => {
	// 	return newComponents.map(e => {
	// 		const { team, component, environment } = getSeparatedConfigId(e);
	// 		return CostingService.getInstance()
	// 			.getComponentCostStream(team || '', environment || '', component || '')
	// 			.subscribe(d => {
	// 				const newItems = componentArrayRef.current.map(nc => {
	// 					if (nc.displayValue === d.id) {
	// 						nc.componentCost = d.estimatedCost;
	// 						nc.componentStatus = d.status;
	// 						nc.costResources = d.costResources;
	// 						nc.syncFinishedAt = d.lastReconcileDatetime;
	// 						nc.isDestroy = d.isDestroyed;
	// 					}
	// 					return nc;
	// 				});
	// 				componentArrayRef.current = newItems;
	// 				setComponents(newItems);
	// 			});
	// 	});
	// };

	// const resolveFetchs = (fetchResps: PromiseSettledResult<any>[]) => {
	// 	const [data0, data1] = fetchResps;
	// 	const respData: any = {
	// 		data0: null,
	// 		data1: null,
	// 	};
	// 	if (data1.status === 'fulfilled' && data1.value && 'data' in data1.value) {
	// 		const { data } = data1.value;
	// 		respData.data1 = data;
	// 	}
	// 	if (data0.status === 'fulfilled' && data0.value && 'data' in data0.value) {
	// 		const { data } = data0.value;
	// 		if ('components' in data) {
	// 			respData.data0 = data.components;
	// 		}
	// 	}
	// 	return respData;
	// };

	const syncStatusMatch = (item: EnvironmentComponentItem): boolean => {
		return syncStatusFilter.has(item.componentStatus);
	};

	const setToggleFilterDropDownValue = (val: boolean) => {
		toggleFilterDropDown(val);
	};

	// const handleVisualizationFlow = (configName: string) => {
	// 	if (!configName.includes('eks') && !configName.includes('networking') && !configName.includes('ec2')) return -1;

	// 	const toastId = `visualization_call_for_${configName}`;
	// 	notificationManager?.show({
	// 		content: (
	// 			<div style={{ display: 'flex', alignItems: 'center' }}>
	// 				<Loader height={16} width={16} color="white" />
	// 				<span style={{ marginLeft: 10 }}>Loading visualization for {configName}...</span>
	// 			</div>
	// 		),
	// 		type: NotificationType.Warning,
	// 		toastOptions: {
	// 			progress: 0,
	// 			toastId,
	// 			autoClose: false,
	// 			draggable: false,
	// 		},
	// 	});
	// 	toast.update(toastId, {
	// 		progress: 0.5,
	// 	});
	// MOCKING THE CALL
	// 	setTimeout(() => {
	// 		AuditService.getInstance()
	// 			.getVisualizationSVGDemo({
	// 				team: projectId,
	// 				environment: environmentId.replace(projectId + '-', ''),
	// 				component: configName,
	// 			})
	// 			.then(response => {
	// 				toast.update(toastId, {
	// 					progress: 0.75,
	// 				});
	// 				if (response?.data?.includes('</svg>')) {
	// 					toast.done(toastId);
	// 					return;
	// 				}
	// 				throw 'No visualization found';
	// 			})
	// 			.catch(err => {
	// 				toast.dismiss(toastId);
	// 				notificationManager?.show({
	// 					content: 'Failed to load visualization. Please contact admin.',
	// 					type: NotificationType.Error,
	// 				});
	// 			});
	// 	}, 1000);
	// };

	const onNodeClick = (configName: string, visualizationHandler?: boolean): void => {
		// if (visualizationHandler) {
		// 	const v = handleVisualizationFlow(configName);
		// 	if (v !== -1) {
		// 		return;
		// 	}
		// }
		if (configName === 'root') {
			// const rootEnv = environments.find(e => e.name === environmentName);
			// const failedEnv = ErrorStateService.getInstance().errorsInEnvironment(rootEnv?.labels?.env_name || '');

			// if (failedEnv?.length > 0) {
			// 	setEnvErrors(failedEnv);
			// }

			setEnvironmentNodeSelected(true);
			setShowSidePanel(true);
			return;
		}
		setEnvironmentNodeSelected(false);
		const selectedConfig = componentArrayRef.current.find(c => c.name === configName);
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
	};

	const labelsMatch = (labels: EnvironmentItem['labels'] = {}, query: string): boolean => {
		return Object.values(labels).some(val => val.includes(query));
	};

	const getFilteredData = (): EnvironmentComponentItem[] => {
		return [];
		// let filteredItems = [...components];
		// if (syncStatusFilter.size > 0) {
		// 	filteredItems = [...filteredItems.filter(syncStatusMatch)];
		// }

		// return filteredItems.filter(item => {
		// 	return item.name.toLowerCase().includes(query) || labelsMatch(item.labels, query);
		// });
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
			case 'list':
				return (
					<div className="zlifecycle-table">
						<ZTable table={{ columns: configTableColumns, rows: getFilteredData() }} />
					</div>
				);
			case 'DAG':
				return components.length > 0 ? (
					<TreeComponent
						environmentId={environmentName}
						nodes={components}
						environmentItem={environment}
						onNodeClick={onNodeClick}
					/>
				) : (
					<></>
				);
			default:
				return (
					<EnvironmentComponentCards
						showAll={showAll}
						components={components || []}
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

	// const checkForFailedEnvironments = (envs: EnvironmentsList) => {
	// 	if (!envs?.length) {
	// 		return;
	// 	}
	// 	const currentEnv = (envs as any).find((e: any) => e.id === environmentName);
	// 	const failedEnv = ErrorStateService.getInstance().errorsInEnvironment(currentEnv.labels?.env_name);
	// 	if (failedEnv?.length && currentEnv.labels?.env_status) {
	// 		currentEnv.labels.env_status = ZSyncStatus.ProvisionFailed;
	// 	}
	// };

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
										selected={envErrors?.length ? 'Errors' : 'Audit'}
										tabs={envTabs.filter(t => t.show(() => Boolean(envErrors?.length)))}>
										<div id="Errors">
											<ErrorView columns={eventErrorColumns} dataRows={envErrors} />
										</div>
										<div id="Audit">
											<AuditView
												auditData={envAuditList}
												auditColumns={auditColumns}
												auditId={environmentName}
											/>
										</div>
									</ZTablControl>
								</div>
							</div>
						)}
						<ZLoaderCover loading={isLoadingWorkflow}>
							{
								selectedConfig && !isEnvironmentNodeSelected && (
									// (selectedConfig.labels?.component_type === 'argocd' ? (
									// 	<ConfigWorkflowViewApplication
									// 		projectId={projectId}
									// 		environmentId={environmentName}
									// 		config={selectedConfig}
									// 	/>
									// ) : (
									<ConfigWorkflowView
										projectId={projectId}
										environmentId={environmentName}
										config={selectedConfig}
										logs={logs}
										plans={plans}
										workflowData={workflowData}
										auditData={compAuditList}
									/>
								)
								// ))
							}
						</ZLoaderCover>
					</ZSidePanel>
				</section>
			</ZLoaderCover>
		</div>
	);
};

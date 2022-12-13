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
import {
	ApplicationWatchEvent,
	HealthStatusCode,
	ZComponentSyncStatus,
	ZSyncStatus,
} from 'models/argo.models';
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
	getWorkflowLogs,
} from 'pages/authorized/environment-components/helpers';
import React, { useContext, useEffect, useRef, useState } from 'react';
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
	const { fetch } = useApi(ArgoComponentsService.getComponents);
	const fetchEnvironments = useApi(ArgoEnvironmentsService.getEnvironments).fetch;
	const { fetch: fetchWorkflowData } = useApi(ArgoWorkflowsService.getConfigWorkflow);
	const [syncStatusFilter, setSyncStatusFilter] = useState<Set<ZSyncStatus>>(new Set<ZSyncStatus>());
	const [filterDropDownOpen, toggleFilterDropDown] = useState(false);
	const [checkBoxFilters, setCheckBoxFilters] = useState<JSX.Element>(<></>);
	const [filterItems, setFilterItems] = useState<Array<() => JSX.Element>>([]);
	const { projectId, environmentId } = useParams<any>();
	const [isLoadingWorkflow, setIsLoadingWorkflow] = useState<boolean>();
	const showAll = environmentId === 'all' && projectId === 'all';
	const notificationManager = React.useContext(Context)?.notifications;
	const [environments, setEnvironments] = useState<EnvironmentsList>([]);
	const { pageHeaderObservable, breadcrumbObservable } = usePageHeader();
	const [query, setQuery] = useState<string>('');
	const [showSidePanel, setShowSidePanel] = useState<boolean>(false);
	const [selectedConfig, setSelectedConfig] = useState<EnvironmentComponentItem>();
	const [loading, setLoading] = useState<boolean>(true);
	const [components, setComponents] = useState<EnvironmentComponentsList>([]);
	const [workflowData, setWorkflowData] = useState<any>();
	const [streamData, setStreamData] = useState<ApplicationWatchEvent | null>(null);
	const [streamData2, setStreamData2] = useState<ApplicationWatchEvent | null>(null);
	const [logs, setLogs] = useState<string | null>(null);
	const [plans, setPlans] = useState<string | null>(null);
	const [workflowId, setWorkflowId] = useState<string>('');
	const [viewType, setViewType] = useState<string>(showAll ? '' : 'DAG');
	const [isEnvironmentNodeSelected, setEnvironmentNodeSelected] = useState<boolean>(false);
	const componentArrayRef = useRef<EnvironmentComponentItem[]>([]);
	const [envErrors, setEnvErrors] = useState<any[]>();
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
			path: `/${projectId}/${environmentId}`,
			name: environmentId,
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
				.bind(null, (status: string) => components.filter(e => e.componentStatus === status).length),
		]);
	}, [components, syncStatusFilter]);

	useEffect(() => {
		setCheckBoxFilters(getCheckBoxFilters(filterDropDownOpen, setToggleFilterDropDownValue, filterItems));
	}, [filterItems, filterDropDownOpen, filterDropDownOpen]);

	useEffect(() => {
		const headerTabs = [
			...environments.map(environment => {
				const name: string = environmentName(environment);
				return {
					active: environmentId === environment.id,
					name: name.charAt(0).toUpperCase() + name.slice(1),
					path: `/${projectId}/${environment.id}`,
				};
			}),
		];

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
	}, [checkBoxFilters, viewType, environments]);

	const syncStatusMatch = (item: EnvironmentComponentItem): boolean => {
		return syncStatusFilter.has(item.componentStatus);
	};

	const setToggleFilterDropDownValue = (val: boolean) => {
		toggleFilterDropDown(val);
	};

	const handleVisualizationFlow = (configName: string) => {
		if (!configName.includes('eks') && !configName.includes('networking') && !configName.includes('ec2')) return -1;

		const toastId = `visualization_call_for_${configName}`;
		notificationManager?.show({
			content: (
				<div style={{ display: 'flex', alignItems: 'center' }}>
					<Loader height={16} width={16} color="white" />
					<span style={{ marginLeft: 10 }}>Loading visualization for {configName}...</span>
				</div>
			),
			type: NotificationType.Warning,
			toastOptions: {
				progress: 0,
				toastId,
				autoClose: false,
				draggable: false,
			},
		});
		toast.update(toastId, {
			progress: 0.5,
		});
		// MOCKING THE CALL
		setTimeout(() => {
			AuditService.getInstance()
				.getVisualizationSVGDemo({
					team: projectId,
					environment: environmentId.replace(projectId + '-', ''),
					component: configName,
				})
				.then(response => {
					toast.update(toastId, {
						progress: 0.75,
					});
					if (response?.data?.includes('</svg>')) {
						toast.done(toastId);
						return;
					}
					throw 'No visualization found';
				})
				.catch(err => {
					toast.dismiss(toastId);
					notificationManager?.show({
						content: 'Failed to load visualization. Please contact admin.',
						type: NotificationType.Error,
					});
				});
		}, 1000);
	};

	const onNodeClick = (configName: string, visualizationHandler?: boolean): void => {
		if (visualizationHandler) {
			const v = handleVisualizationFlow(configName);
			if (v !== -1) {
				return;
			}
		}
		if (configName === 'root') {
			// const rootEnv = environments.find(e => e.id === environmentId);
			// const failedEnv = ErrorStateService.getInstance().errorsInEnvironment(
			// 	rootEnv?.labels?.env_name || ''
			// );

			// if (failedEnv?.length > 0) {
			// 	setEnvErrors(failedEnv);
			// }

			// setEnvironmentNodeSelected(true);
			// setShowSidePanel(true);
			return;
		}
		setEnvironmentNodeSelected(false);
		const selectedConfig = componentArrayRef.current.find(config => config.componentName === configName);
		if (selectedConfig) {
			let workflowId = selectedConfig.labels?.last_workflow_run_id || '';
			if (!selectedConfig.labels?.last_workflow_run_id) {
				workflowId = 'initializing';
			}
			setSelectedConfig(selectedConfig);
			setShowSidePanel(true);
			setWorkflowId(workflowId);
		}
	};

	useEffect(() => {
		fetchEnvironments(projectId).then(({ data }) => {
			if (data) {
				checkForFailedEnvironments(data);
				setEnvironments(data);
			}
			setLoading(false);
		});
	}, [projectId]);

	useEffect(() => {
		const newEnvironments = streamMapper<EnvironmentItem>(
			streamData,
			environments,
			ArgoMapper.parseEnvironment,
			'environment',
			{
				projectId,
			}
		);
		checkForFailedEnvironments(newEnvironments);
		setEnvironments(newEnvironments);

		const newComponents = streamMapper<EnvironmentComponentItem>(
			streamData,
			components,
			ArgoMapper.parseComponent,
			'config',
			{
				projectId,
				environmentId,
			}
		);
		setComponents(newComponents);
		componentArrayRef.current = newComponents;
		const selectedConf = newComponents.find((itm: any) => itm.id === selectedConfig?.id);
		if (selectedConf) {
			setSelectedConfig(selectedConf);
			if (selectedConf.labels?.last_workflow_run_id !== workflowId) {
				setWorkflowId(selectedConf.labels?.last_workflow_run_id || '');
			}
		}
	}, [streamData]);

	useEffect(() => {
		if (environments.length === 0 || !environmentId) return;
		const env = environments.find(e => e.id === environmentId);
		if (!env) return;
		const sub = ctx?.failedEnvironments.subscribe(res => {
			const errors = [...res.values()].filter(e => e.labels.env_name === env.labels?.env_name);
			console.log(errors);
		});

		return () => sub?.unsubscribe();
	}, [environmentId, environments]);

	useEffect(() => {
		const $subscription: Subscription = subscriber.subscribe((response: any) => {
			setStreamData(response);
		});
		const $subscription2: Subscription = subscriberWF.subscribe((response: any) => {
			setStreamData2(response);
		});

		return (): void => {
			$subscription.unsubscribe();
			$subscription2.unsubscribe();
		};
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

	useEffect(() => {
		setLoading(true);
		fetch(projectId, environmentId).then(({ data }) => {
			if (data) {
				componentArrayRef.current = data;
				setComponents(data);
			}
			setLoading(false);
		});
	}, ['projectId', environmentId]);

	useEffect(() => {
		if (showSidePanel === false) {
			setLogs(null);
			setPlans(null);
			setIsLoadingWorkflow(false);
			setWorkflowData(null);
			setWorkflowId('');
			setSelectedConfig(undefined);
		}
	}, [showSidePanel]);

	const labelsMatch = (labels: EnvironmentItem['labels'] = {}, query: string): boolean => {
		return Object.values(labels).some(val => val.includes(query));
	};

	const getFilteredData = (): EnvironmentComponentItem[] => {
		let filteredItems = [...components];
		if (syncStatusFilter.size > 0) {
			filteredItems = [...filteredItems.filter(syncStatusMatch)];
		}

		return filteredItems.filter(item => {
			return item.name.toLowerCase().includes(query) || labelsMatch(item.labels, query);
		});
	};

	useEffect(() => {
		if (workflowId) {
			if (workflowId === 'initializing') {
				setLogs(null);
				setPlans(null);
				setIsLoadingWorkflow(false);
				setWorkflowData(null);
			} else {
				getWorkflowData(workflowId, selectedConfig?.id || '');
			}
		}
	}, [workflowId]);

	const getWorkflowData = (workflowId: string, configId: string) => {
		setLogs(null);
		setPlans(null);
		setIsLoadingWorkflow(true);
		fetchWorkflowData({
			projectId: projectId,
			environmentId: environmentId,
			configId: configId,
			workflowId: workflowId,
		}).then(({ data }) => {
			setIsLoadingWorkflow(false);
			setWorkflowData(data);
			const configParamsSet: ConfigParamsSet = {
				projectId,
				environmentId,
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
				if (
					components.length > 0 &&
					environments.length > 0 &&
					components[0].labels?.environment_id === environmentId
				) {
					return (
						<TreeComponent
							environmentId={environmentId}
							nodes={components}
							environmentItem={environments.find(e => e.id === environmentId)}
							onNodeClick={onNodeClick}
						/>
					);
				}
				break;
			default:
				return (
					<EnvironmentComponentCards
						showAll={showAll}
						components={components ? getFilteredData() : []}
						projectId={projectId}
						envName={environmentName(environments.find(e => e.id === environmentId))}
						selectedConfig={selectedConfig}
						workflowPhase={workflowData?.status?.phase}
						onClick={(config: EnvironmentComponentItem): void => {
							onNodeClick(config.componentName);
						}}
					/>
				);
		}
	};

	const checkForFailedEnvironments = (envs: EnvironmentsList) => {
		if (!envs?.length) {
			return;
		}
		const currentEnv = (envs as any).find((e: any) => e.id === environmentId);
		const failedEnv = ErrorStateService.getInstance().errorsInEnvironment(currentEnv.labels?.env_name);
		if (failedEnv?.length && currentEnv.labels?.env_status) {
			currentEnv.labels.env_status = ZSyncStatus.ProvisionFailed;
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
										selected={envErrors?.length ? 'Errors' : 'Audit'}
										tabs={envTabs.filter(t => t.show(() => Boolean(envErrors?.length)))}>
										<div id="Errors">
											<ErrorView
												columns={eventErrorColumns}
												dataRows={envErrors}
											/>
										</div>
										<div id="Audit">
											<AuditView
												fetch={AuditService.getInstance().getEnvironment}
												auditColumns={auditColumns}
												auditId={environmentId}
											/>
										</div>
									</ZTablControl>
								</div>
							</div>
						)}
						<ZLoaderCover loading={isLoadingWorkflow}>
							{selectedConfig &&
								!isEnvironmentNodeSelected &&
								(selectedConfig.labels?.component_type === 'argocd' ? (
									<ConfigWorkflowViewApplication
										projectId={projectId}
										environmentId={environmentId}
										config={selectedConfig}
									/>
								) : (
									<ConfigWorkflowView
										projectId={projectId}
										environmentId={environmentId}
										config={selectedConfig}
										logs={logs}
										plans={plans}
										workflowData={workflowData}
									/>
								))}
						</ZLoaderCover>
					</ZSidePanel>
				</section>
			</ZLoaderCover>
		</div>
	);
};

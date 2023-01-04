import { NotificationType } from 'components/argo-core';
import { ZLoaderCover } from 'components/atoms/loader/LoaderCover';
import { ZTable } from 'components/atoms/table/Table';
import { EnvironmentCards } from 'components/molecules/cards/EnvironmentCards';
import { Context } from 'context/argo/ArgoUi';
import { streamMapper } from 'helpers/streamMapper';
import { useApi } from 'hooks/use-api/useApi';
import { ApplicationWatchEvent, ZEnvSyncStatus, ZSyncStatus } from 'models/argo.models';
import { LocalStorageKey } from 'models/localStorage';
import { EnvironmentItem, EnvironmentsList, PageHeaderTabs, TeamsList } from 'models/projects.models';
import {
	environmentTableColumns,
	getCheckBoxFilters,
	mockModifiedYaml,
	mockOriginalYaml,
	renderSyncStatusItems,
} from 'pages/authorized/environments/helpers';
import React, { useEffect, useMemo, useState } from 'react';
import { useParams } from 'react-router-dom';
import { Subscription } from 'rxjs';
import { ArgoEnvironmentsService } from 'services/argo/ArgoEnvironments.service';
import { ArgoMapper } from 'services/argo/ArgoMapper';
import { ArgoTeamsService } from 'services/argo/ArgoProjects.service';
import { subscriber } from 'utils/apiClient/EventClient';
import { usePageHeader } from '../contexts/EnvironmentHeaderContext';
import { DiffEditor } from '@monaco-editor/react';
import { ErrorStateService } from 'services/error/error-state.service';
import AuthStore from 'auth/AuthStore';
import { EntityStore, Environment, Team } from 'models/entity.store';

type CompareEnv = {
	env: EnvironmentItem | null;
};

type CompareEnvs = {
	a: CompareEnv | null;
	b: CompareEnv | null;
};

export const Environments: React.FC = () => {
	const { projectId } = useParams();
	const entityStore = useMemo(() => EntityStore.getInstance(), []);
	const { fetch } = useApi(ArgoEnvironmentsService.getEnvironments);
	const [filterDropDownOpen, toggleFilterDropDown] = useState(false);
	const [query, setQuery] = useState<string>('');
	const [syncStatusFilter, setSyncStatusFilter] = useState<Set<ZEnvSyncStatus>>(new Set<ZEnvSyncStatus>());
	const [loading, setLoading] = useState<boolean>(true);
	const [environments, setEnvironments] = useState<Environment[]>([]);
	// Need to discuss about these
	const [failedEnvironments, setFailedEnvironments] = useState<EnvironmentsList>([]);
	const [viewType, setViewType] = useState<string>('');
	// Filters [status]
	const [checkBoxFilters, setCheckBoxFilters] = useState<JSX.Element>(<></>);
	const [teams, setTeams] = useState<Team[]>([]);
	const [filterItems, setFilterItems] = useState<Array<() => JSX.Element>>([]);
	const { pageHeaderObservable, breadcrumbObservable } = usePageHeader();
	const ctx = React.useContext(Context);
	const errorStateService = ErrorStateService.getInstance();
	const [compareMode, setCompareMode] = useState<boolean>(false);
	const [compareEnvs, setCompareEnvs] = useState<CompareEnvs>({
		a: null,
		b: null,
	});

	const breadcrumbItems = [
		{
			path: '/dashboard',
			name: 'All Environments',
			active: false,
		},
		{
			path: `/${projectId}`,
			name: projectId,
			active: projectId !== undefined,
		},
	];

	useEffect(() => {
		const $subscription: Subscription[] = [];
		$subscription.push(
			entityStore.emitter.subscribe(data => {
				if (data.length === 0) return;
				setEnvironments(
					projectId ? entityStore.getAllEnvironmentsByTeamName(projectId) : entityStore.Environments
				);
				setLoading(false);
			})
		);
		// $subscription.push(
		// 	errorStateService.updates.subscribe(() => {
		// 		environments.forEach(e => checkForFailedEnvironments(e));
		// 	})
		// );
		return (): void => $subscription.forEach(e => e.unsubscribe());
	}, [projectId]);

	// useEffect(() => {
	// fetch(projectId).then(({ data }) => {
	// 	if (data) {
	// 		data.forEach(e => checkForFailedEnvironments(e));
	// 		setEnvironments(data);
	// 	}
	// 	setLoading(false);
	// });
	// if (!projectId) return;

	// setLoading(true);
	// const $subscription: Subscription[] = [];
	// const envSub = entityStore
	// .getEnvironmentsEmitter(projectId)?.subscribe(data => {
	// 	if (data.length === 0) return;
	// 	setEnvironments(data);
	// 	setLoading(false);
	// });
	// envSub && $subscription.push(envSub);
	// const setFailedEnv = (envs: any) => {
	// 	const list: any = envs
	// 		.filter((e: any) => !environments.some(en => en.labels?.env_name === e.environment))
	// 		.map((e: any) => ({
	// 			id: `${e.team}-${e.environment}`,
	// 			name: e.environment,
	// 			labels: {
	// 				project_id: e.team,
	// 				env_name: e.environment,
	// 				env_status: ZSyncStatus.ProvisionFailed,
	// 				failed_environment: true,
	// 			},
	// 		}));
	// 	if (projectId) {
	// 		setFailedEnvironments(list.filter((e: any) => e.labels?.project_id === projectId));
	// 	} else {
	// 		setFailedEnvironments(list);
	// 	}
	// };
	// setFailedEnv(errorStateService.ErrorsEnvs || []);
	// $subscription.push(
	// 	errorStateService.updates.subscribe(() => {
	// 		setFailedEnv(errorStateService.ErrorsEnvs);
	// 	})
	// );
	// return (): void => $subscription.forEach(e => e.unsubscribe());
	// }, [projectId]);

	// useEffect(() => {
	// 	const newItems = streamMapper<EnvironmentItem>(
	// 		streamData,
	// 		environments,
	// 		ArgoMapper.parseEnvironment,
	// 		'environment'
	// 	);
	// 	setEnvironments(newItems);
	// }, [streamData, environments]);

	// const checkForFailedEnvironments = (currentEnv: EnvironmentItem) => {
	// 	if (!currentEnv) {
	// 		return;
	// 	}
	// 	const failedEnv = errorStateService.errorsInEnvironment(currentEnv.labels?.env_name || '');
	// 	if (failedEnv?.length && currentEnv.labels?.env_status) {
	// 		currentEnv.labels.env_status = ZSyncStatus.ProvisionFailed;
	// 	}
	// };

	const setToggleFilterDropDownValue = (val: boolean) => {
		toggleFilterDropDown(val);
	};

	const labelsMatch = (labels: EnvironmentItem['labels'] = {}, query: string): boolean => {
		return Object.values(labels).some(val => val.toString().includes(query));
	};

	const syncStatusMatch = (item: EnvironmentItem): boolean => {
		return syncStatusFilter.has(item.labels?.env_status as ZEnvSyncStatus);
	};

	useEffect(() => {
		setFilterItems([
			renderSyncStatusItems
				.bind(null, ZEnvSyncStatus, syncStatusFilter, setSyncStatusFilter, 'Environment Status')
				.bind(
					null,
					(status: string) =>
						[...environments, ...failedEnvironments].filter(e => 'Unknown' === status).length
				),
		]);
	}, [query, environments, syncStatusFilter, failedEnvironments]);

	useEffect(() => {
		setCheckBoxFilters(getCheckBoxFilters(filterDropDownOpen, setToggleFilterDropDownValue, filterItems));
	}, [filterItems, filterDropDownOpen, filterDropDownOpen]);

	useEffect(() => {
		const headerTabs: PageHeaderTabs = [
			{ name: 'All', path: '/dashboard', active: projectId === undefined },
			...entityStore.Teams.map(team => {
				const teamId = team.name;
				return {
					active: projectId === teamId,
					name: teamId.charAt(0).toUpperCase() + teamId.slice(1),
					path: `/${teamId}`,
				};
			}),
		];
		pageHeaderObservable.next({
			breadcrumbs: breadcrumbItems,
			headerTabs,
			pageName: 'Environments',
			filterTitle: 'Filter by team:',
			onSearch: setQuery,
			onViewChange: setViewType,
			buttonText: 'New environment',
			checkBoxFilters: checkBoxFilters,
			diffChecker: {
				setter: setCompareMode,
				getEnvs: getCompareEnvs,
			},
		});

		breadcrumbObservable.next({ [LocalStorageKey.TEAMS]: headerTabs });
	}, [teams, checkBoxFilters]);

	const getCompareEnvs = () => {
		return compareEnvs;
	};

	const getFilteredData = (): Environment[] => {
		return [];
		// let filteredItems = [...environments, ...failedEnvironments];
		// if (syncStatusFilter.size > 0) {
		// 	filteredItems = [...filteredItems.filter(syncStatusMatch)];
		// }

		// return filteredItems.filter(item => {
		// 	return item.name.toLowerCase().includes(query) || labelsMatch(item.labels, query);
		// });
	};

	const renderDiffEditor = () => {
		setTimeout(() => {
			const div = document.querySelector('section.dashboard-content div.context-menu') as HTMLDivElement;
			div.style.opacity = '1';
		}, 100);
		return (
			<div
				className="context-menu"
				style={{ display: 'flex', height: '80vh', opacity: 0, transition: 'all 0.3s' }}>
				<header className="diff-editor-header">
					<h4>Comparison</h4>{' '}
					<button
						onClick={e => {
							setCompareMode(false);
						}}>
						x
					</button>
				</header>
				<DiffEditor language="yaml" original={mockOriginalYaml} modified={mockModifiedYaml} />
			</div>
		);
	};

	const setSelectedEnvironmentToCompare = (environment: EnvironmentItem, setSelected: any) => {
		if (setSelected) {
			if (compareEnvs.a?.env && compareEnvs.b?.env) {
				ctx?.notifications?.show({
					content: 'Deselect an environment to select this one!',
					type: NotificationType.Warning,
				});
				return false;
			}
			if (!compareEnvs.a?.env) {
				compareEnvs.a = {
					env: environment,
				};
				setCompareEnvs({ ...compareEnvs });
				return true;
			}

			if (!compareEnvs.b?.env) {
				compareEnvs.b = {
					env: environment,
				};
				setCompareEnvs({ ...compareEnvs });
				return true;
			}
		} else {
			if (environment === compareEnvs.a?.env) {
				compareEnvs.a.env = null;
				setCompareEnvs({ ...compareEnvs });
				return false;
			} else if (environment === compareEnvs.b?.env) {
				compareEnvs.b.env = null;
				setCompareEnvs({ ...compareEnvs });
				return false;
			}
		}
	};

	return (
		<div className="zlifecycle-page">
			<ZLoaderCover loading={loading}>
				<section className="dashboard-content">
					{viewType === 'list' ? (
						<div className="zlifecycle-table">
							<ZTable table={{ columns: environmentTableColumns, rows: getFilteredData() }} />
						</div>
					) : (
						<EnvironmentCards
							// environments={environments ? getFilteredData() : []}
							environments={environments}
							compareEnabled={{
								compareMode,
								setSelectedEnvironmentToCompare,
							}}
						/>
					)}
					{compareEnvs.a?.env && compareEnvs.b?.env ? compareMode && renderDiffEditor() : null}
				</section>
			</ZLoaderCover>
		</div>
	);
};

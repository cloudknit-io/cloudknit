import { DiffEditor } from '@monaco-editor/react';
import { NotificationType, NotificationsApi } from 'components/argo-core';
import { ZLoaderCover } from 'components/atoms/loader/LoaderCover';
import { EnvironmentCards } from 'components/molecules/cards/EnvironmentCards';
import { ZModalPopup } from 'components/molecules/modal/ZModalPopup';
import { Context } from 'context/argo/ArgoUi';
import { ZEnvSyncStatus } from 'models/argo.models';
import { EntityStore } from 'models/entity.store';
import { Environment } from 'models/entity.type';
import { LocalStorageKey } from 'models/localStorage';
import { EnvironmentItem, PageHeaderTabs } from 'models/projects.models';
import {
	getCheckBoxFilters,
	mockModifiedYaml,
	mockOriginalYaml,
	renderSyncStatusItems,
} from 'pages/authorized/environments/helpers';
import React, { useEffect, useMemo, useState } from 'react';
import { useParams } from 'react-router-dom';
import { Subscription } from 'rxjs';
import { EntityService } from 'services/entity/entity.service';
import { usePageHeader } from '../contexts/EnvironmentHeaderContext';
import { playgroundFeatureVisible } from '../feature_toggle';

type CompareEnv = {
	env: EnvironmentItem | null;
};

type CompareEnvs = {
	a: CompareEnv | null;
	b: CompareEnv | null;
};

export const Environments: React.FC = () => {
	const { projectId } = useParams() as any;
	const entityStore = useMemo(() => EntityStore.getInstance(), []);
	const [filterDropDownOpen, toggleFilterDropDown] = useState(false);
	const [query, setQuery] = useState<string>('');
	const [syncStatusFilter, setSyncStatusFilter] = useState<Set<ZEnvSyncStatus>>(new Set<ZEnvSyncStatus>());
	const [loading, setLoading] = useState<boolean>(true);
	const [environments, setEnvironments] = useState<Environment[]>([]);
	const [viewType, setViewType] = useState<string>('');
	const [pushingCommit, setPushingCommit] = useState<boolean | null>(null);
	const [commitInfo, setCommitInfo] = useState<string | null>(null);
	const [checkBoxFilters, setCheckBoxFilters] = useState<JSX.Element>(<></>);
	const [filterItems, setFilterItems] = useState<Array<() => JSX.Element>>([]);
	const { pageHeaderObservable, breadcrumbObservable } = usePageHeader();
	const ctx = React.useContext(Context);
	const [compareMode, setCompareMode] = useState<boolean>(false);
	const [compareEnvs, setCompareEnvs] = useState<CompareEnvs>({
		a: null,
		b: null,
	});
	const nm = React.useContext(Context)?.notifications as NotificationsApi;

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
				if (data.environments.length === 0) return;
				setEnvironments([
					...(projectId ? entityStore.getAllEnvironmentsByTeamName(projectId) : entityStore.Environments),
				]);
				setLoading(false);
			})
		);
		return (): void => $subscription.forEach(e => e.unsubscribe());
	}, [projectId]);

	const setToggleFilterDropDownValue = (val: boolean) => {
		toggleFilterDropDown(val);
	};

	const labelsMatch = (labels: EnvironmentItem['labels'] = {}, query: string): boolean => {
		return Object.values(labels).some(val => val.toString().includes(query));
	};

	const syncStatusMatch = (item: Environment): boolean => {
		return syncStatusFilter.has(item.status as ZEnvSyncStatus);
	};

	useEffect(() => {
		setFilterItems([
			renderSyncStatusItems
				.bind(null, ZEnvSyncStatus, syncStatusFilter, setSyncStatusFilter, 'Environment Status')
				.bind(null, (status: string) => [...environments].filter(e => e.status === status).length),
		]);
	}, [query, environments, syncStatusFilter]);

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
	}, [checkBoxFilters]);

	const getCompareEnvs = () => {
		return compareEnvs;
	};

	const getFilteredData = (): Environment[] => {
		let filteredItems = [...environments];
		if (syncStatusFilter.size > 0) {
			filteredItems = [...filteredItems.filter(syncStatusMatch)];
		}

		return filteredItems.filter(item => {
			return item.name.toLowerCase().includes(query);
		});
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
			<ZLoaderCover loading={loading || pushingCommit === true}>
				<section className="dashboard-content">
					{viewType === 'list' ? (
						<div className="zlifecycle-table">
							{/* <ZTable table={{ columns: environmentTableColumns, rows: getFilteredData() }} /> */}
						</div>
					) : (
						<EnvironmentCards
							environments={environments ? getFilteredData() : []}
							compareEnabled={{
								compareMode,
								setSelectedEnvironmentToCompare,
							}}
						/>
					)}
					{compareEnvs.a?.env && compareEnvs.b?.env ? compareMode && renderDiffEditor() : null}
				</section>
				{!playgroundFeatureVisible() && (
					<ZModalPopup
						header={<div className="d-flex align-center">Provison an Environment</div>}
						isShown={
							!loading &&
							environments?.length > 0 &&
							environments[0].status === ZEnvSyncStatus.Destroyed &&
							pushingCommit !== false
						}
						onClose={() => {}}>
						<div className="d-flex flex-dir-column">
							<div style={{ display: 'block' }}>
								<small>
									Clicking on this button will push a commit to our repository and cloudknit will
									start provisioning your environment.
								</small>
							</div>
							<div className="d-flex flex-dir-row align-center justify-between mt-10">
								<div
									className="d-flex flex-dir-column"
									style={{ background: '#eee', padding: '10px', borderRadius: '5px' }}>
									<small>
										<em>This is the repository where you will commit.</em>
									</small>
									<small>
										<a
											style={{ color: 'teal' }}
											href="https://github.com/zlab-tech/hooli-config"
											target="_blank">
											<i>https://github.com/zlab-tech/hooli-config</i>
										</a>
									</small>
								</div>
								<button
									onClick={() => {
										setPushingCommit(true);
										EntityService.getInstance()
											.gitCommit(environments[0].teamId, environments[0].id)
											.then(({ status, html_url }: any) => {
												if (status === 'error') {
													nm.show({
														content: 'There was an error provisioning the environment',
														type: NotificationType.Error,
													});
													setCommitInfo(null);
												} else {
													nm.show({
														content: 'Well Done! Provisioning your environment...',
														type: NotificationType.Success,
													});
													setCommitInfo(html_url);
												}
												setPushingCommit(false);
											});
									}}
									className="btn shadowy-input btn__update">
									Create a Commit
								</button>
							</div>
						</div>
					</ZModalPopup>
				)}
				{!playgroundFeatureVisible() && (
					<ZModalPopup
						header={<div className="d-flex align-center">Your commit was successful.</div>}
						isShown={commitInfo !== null}
						onClose={() => {}}>
						<div className="d-flex flex-dir-column align-center">
							<div style={{ display: 'block' }}>
								<small>You can see your commit by clicking on the below link.</small>
							</div>
							<div>
								<small>
									<i>
										<a style={{ color: 'teal' }} href={commitInfo as string} target="_blank">
											{commitInfo}
										</a>
									</i>
								</small>
							</div>
							<div>
								<button
									className="btn shadowy-input btn__update mt-10"
									onClick={() => {
										setCommitInfo(null);
									}}>
									Close
								</button>
							</div>
						</div>
					</ZModalPopup>
				)}
			</ZLoaderCover>
		</div>
	);
};

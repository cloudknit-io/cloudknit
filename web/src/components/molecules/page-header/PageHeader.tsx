import './style.scss';

import { ReactComponent as DAGViewIcon } from 'assets/images/icons/DAG-view.svg';
import { ReactComponent as ClearInputIcon } from 'assets/images/icons/field/Backspace.svg';
import { ReactComponent as GridViewIcon } from 'assets/images/icons/grid-view-solid.svg';
import { ReactComponent as ListViewIcon } from 'assets/images/icons/list-view-solid.svg';
import { ReactComponent as SearchIcon } from 'assets/images/icons/search.svg';
import { ReactComponent as OutOfSyncIcon } from 'assets/images/icons/card-status/sync/Not Sync.svg';
import classNames from 'classnames';
import { Button } from 'components/atoms/button/Button';
import { ZText } from 'components/atoms/text/Text';
import { PageHeaderTabs } from 'models/projects.models';
import React, { ChangeEvent, FC, useEffect, useState } from 'react';
import { useHistory, useParams } from 'react-router-dom';
import { ReactComponent as Compare } from 'assets/images/icons/compare.svg';
import { BradAdarshFeatureVisible, FeatureKeys, featureToggled } from 'pages/authorized/feature_toggle';
import { ZSidePanel } from '../side-panel/SidePanel';
import { ErrorView } from 'components/organisms/error-view/ErrorView';
import { ErrorStateService } from 'services/error/error-state.service';
import { eventErrorColumns, EventMessage } from 'models/error.model';
import AuthStore from 'auth/AuthStore';

type Props = {
	breadcrumbs: {
		path: string;
		name: string;
		active: boolean;
	}[];
	initialView?: string;
	buttonText: string;
	onClickButton: () => void;
	filterTitle: string;
	showBtn?: boolean;
	pageName: string;
	headerTabs: PageHeaderTabs;
	onSearch: (query: string) => void;
	onViewChange?: (viewType: string) => void;
	checkBoxFilters: JSX.Element;
	diffChecker: any;
};

export const ZPageHeader: FC<Props> = ({
	breadcrumbs,
	buttonText,
	initialView,
	onClickButton,
	onSearch,
	pageName,
	showBtn = false,
	headerTabs,
	filterTitle,
	onViewChange,
	checkBoxFilters,
	diffChecker,
}: Props) => {
	const history = useHistory();
	const view = initialView || 'grid';
	const [viewType, setViewType] = useState<string>(view);
	const [errors, setErrors] = useState<EventMessage[]>([]);
	const [showErrors, toggleShowErrors] = useState<boolean>(false);
	const { projectId } = useParams() as any;
	const noViewType = ['Environment Builder', 'Settings', 'Resource View', null];
	const handleOnViewChange = (type: string) => {
		setViewType(type);
		onViewChange && onViewChange(type);
	};

	useEffect(() => {
		setViewType(view);
	}, [initialView]);

	useEffect(() => {
		setViewType(view);
		setQuery('');
		onSearch('');
		onViewChange && onViewChange(view === 'grid' ? '' : view);
	}, [onViewChange, onSearch]);

	const [query, setQuery] = useState<string>('');

	useEffect(() => {
		if (!AuthStore.getOrganization()) {
			return;
		}
		const sub = ErrorStateService.getInstance().updates.subscribe(() => {
			const errors = ErrorStateService.getInstance().Errors;
			if (projectId) {
				setErrors(errors.filter(e => e.team === projectId));
			} else {
				setErrors(errors);
			}
		});
		return () => sub?.unsubscribe();
	}, [projectId]);

	const renderDiffChecker = () => {
		return (
			<button
				id="diff-checker"
				className="diff-checker"
				title="Compare YAML"
				onClick={e => {
					// const compareEnvs = diffChecker.getEnvs();
					// if (compareEnvs.a?.env && compareEnvs.b?.env) {
					diffChecker.setter(true);
					// }
				}}>
				<Compare viewBox="0 0 24 24" height={24} width={24} />
				<span style={{ marginLeft: 5 }}>Compare</span>
			</button>
		);
	};

	const renderErrors = () => {
		return (
			<div
				onClick={() => {
					toggleShowErrors(true);
				}}
				style={{
					fontWeight: 'bold',
					cursor: 'pointer',
					color: '#FE0000',
					display: 'flex',
					alignItems: 'center',
					minWidth: '70px',
					justifyContent: 'space-between',
					marginRight: '45px',
				}}>
				<OutOfSyncIcon /> {errors.length} {errors.length > 1 ? 'Errors' : 'Error'}
			</div>
		);
	};

	return (
		<div className="zlifecycle-page-header">
			<div className="zlifecycle-page-header__left__breadcrumbs">
				{breadcrumbs.length > 0 &&
					breadcrumbs.map((item, index) => (
						<div key={item.path} onClick={(): void => history.push(item.path)}>
							<ZText.Body
								className={item.active ? 'zlifecycle-page-header__left__breadcrumbs--active' : ''}
								size="14"
								weight="bold"
								lineHeight="18">
								{item.name}{' '}
								<span style={{ margin: '0 8px' }}>{index !== breadcrumbs.length - 1 && '/'}</span>
							</ZText.Body>
						</div>
					))}
			</div>
			<ZSidePanel
				isShown={showErrors}
				onClose={() => {
					toggleShowErrors(false);
				}}>
				<ErrorView
					columns={eventErrorColumns}
					dataRows={errors}
				/>
			</ZSidePanel>
			<div className="zlifecycle-page-header__search-filter">
				<div className="zlifecycle-page-header__left">
					<div className="zlifecycle-page-header__left__filters">
						<ZText.Body className="page-offset" size="36" weight="bold">
							{pageName}
						</ZText.Body>
						<div className="op-container">
							{diffChecker && errors.length > 0 && renderErrors()}
							{diffChecker && BradAdarshFeatureVisible() && renderDiffChecker()}
							{checkBoxFilters ? checkBoxFilters : null}
						</div>
						{viewType !== 'DAG' && !noViewType.includes(pageName) && (
							<div className="zlifecycle-page-header__left__search">
								<SearchIcon className={'zlifecycle-page-header__left__search__icon'} />
								<input
									className="zlifecycle-page-header__left__search__input shadowy-input"
									placeholder="Search"
									value={query}
									onChange={(e: ChangeEvent<HTMLInputElement>): void => {
										setQuery(e.target.value);
										onSearch(e.target.value.toLowerCase());
									}}
								/>
								{query !== '' && (
									<ClearInputIcon
										className={'zlifecycle-page-header__left__search__clear'}
										onClick={() => {
											setQuery('');
											onSearch('');
										}}
									/>
								)}
							</div>
						)}
						{headerTabs.length > 0 && (
							<div>
								<ZText.Body
									className="zlifecycle-page-header__tabs-container__title"
									size="14"
									lineHeight="18">
									{filterTitle}
								</ZText.Body>
								<div className="zlifecycle-page-header__tabs-container">
									{headerTabs.map(tab => (
										<div
											className={`${
												tab.active ? 'zlifecycle-page-header__tabs-container__tab--active' : ''
											} zlifecycle-page-header__tabs-container__tab`}
											key={`header-tab${tab.name}`}
											onClick={() => {
												history.push(tab.path);
											}}>
											<ZText.Body
												size="14"
												key={`header-tab-body-${tab.name}`}
												weight="bold"
												lineHeight="18">
												{tab.name}
											</ZText.Body>
											)
										</div>
									))}
								</div>
							</div>
						)}
					</div>
				</div>

				{!noViewType.includes(pageName) && (
					<div className="zlifecycle-page-header__right">
						{showBtn && <Button onClick={onClickButton}>{buttonText}</Button>}
						<div className="zlifecycle-page-header__right__options">
							<ListViewIcon
								className={classNames({ inactive: viewType !== 'list' })}
								onClick={() => handleOnViewChange('list')}
							/>
							<GridViewIcon
								className={classNames({ inactive: viewType !== 'grid' })}
								onClick={() => handleOnViewChange('grid')}
							/>
							{pageName === 'Components' && projectId !== 'all' ? (
								<DAGViewIcon
									className={classNames('z-hover', { inactive: viewType !== 'DAG' })}
									onClick={(): void => handleOnViewChange('DAG')}
								/>
							) : (
								''
							)}
						</div>
					</div>
				)}
			</div>
		</div>
	);
};

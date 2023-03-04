import { ReactComponent as ChevronRight } from 'assets/images/icons/chevron-right.svg';
import { BreadcrumbFilter, BreadcrumbInfo } from 'models/breadcrumb.model';
import { LocalStorageKey } from 'models/localStorage';
import React, { FC, useEffect, useState } from 'react';
import { useHistory, useParams } from 'react-router';
import { LocalStorage } from 'utils/localStorage/localStorage';

import { BreadcrumbObservableValue, usePageHeader } from '../contexts/EnvironmentHeaderContext';
import { FeatureKeys, FeatureRoutes } from '../feature_toggle';
import { Breadcrumb } from './Breadcrumb';

enum FilterRegion {
	TEAMS = 1,
	ENVIRONMENT = 2,
	COMPONENTS = 3,
}

export const Breadcrumbs: FC = () => {
	const [breadcrumbInfo, setBreadcrumbInfo] = useState<BreadcrumbInfo[]>([]);
	const { breadcrumbObservable } = usePageHeader();
	const history = useHistory();
	const { projectId, environmentId } = useParams() as any;
	const currentPath = history.location.pathname;

	const getFilter = (filterRegion: FilterRegion): BreadcrumbFilter[] => {
		let filter: any = [],
			type = '',
			tokens = history.location.pathname.split('/');
		switch (filterRegion) {
			case FilterRegion.TEAMS:
				filter = LocalStorage.getItem<[]>(LocalStorageKey.TEAMS) || [];
				break;
			case FilterRegion.ENVIRONMENT:
				filter = LocalStorage.getItem<[]>(LocalStorageKey.ENVIRONMENTS) || [];
				break;
			default:
				filter = [];
		}

		return filter.map((item: any) => ({
			title: item.name,
			routePath: item.path + (type ? '/' + type : ''),
		}));
	};

	const setBreadcrumbMap = (pathname: string) => {
		if ('/dashboard' === pathname) {
			setBreadcrumbInfo([
				{
					title: 'Home',
					path: '/dashboard',
				},
				{
					title: 'All',
					path: '/dashboard',
					filters: getFilter(FilterRegion.TEAMS),
				},
			]);
		} else {
			const tokens = pathname.split('/');
			setBreadcrumbInfo(
				tokens.filter((_e, i) => i < 3).map(
					(t, _i): BreadcrumbInfo => {
						if (_i === 0) {
							return {
								title: 'Home',
								path: '/dashboard',
							};
						} else {
							return {
								title: t === environmentId ? t.replace(projectId + '-', '') : t,
								path: '/' + tokens.slice(1, _i + 1).join('/'),
								filters: getFilter(_i) || [],
							};
						}
					}
				)
			);
		}
	};

	useEffect(() => {
		setBreadcrumbMap(currentPath);
	}, [currentPath]);

	useEffect(() => {
		const s = breadcrumbObservable.subscribe((filterInfo: BreadcrumbObservableValue | boolean) => {
			const info = filterInfo as BreadcrumbObservableValue;
			if (info[LocalStorageKey.ENVIRONMENTS]) {
				LocalStorage.setItem(LocalStorageKey.ENVIRONMENTS, info[LocalStorageKey.ENVIRONMENTS]);
				setBreadcrumbMap(history.location.pathname);
			} else if (info[LocalStorageKey.TEAMS]) {
				LocalStorage.setItem(LocalStorageKey.TEAMS, info[LocalStorageKey.TEAMS]);
				setBreadcrumbMap(history.location.pathname);
			} else if (filterInfo === false) {
				setBreadcrumbInfo([]);
			}
		});

		return () => s.unsubscribe();
	});

	return (
		<>
			{breadcrumbInfo.map((breadcrumb, _i) => (
				<React.Fragment key={`breadcrumb-${breadcrumb.title}-${_i}`}>
					<Breadcrumb title={breadcrumb.title} filters={breadcrumb.filters} path={breadcrumb.path} />
					{_i !== breadcrumbInfo.length - 1 ? <ChevronRight style={{ margin: '0 8px' }} /> : null}
				</React.Fragment>
			))}
		</>
	);
};

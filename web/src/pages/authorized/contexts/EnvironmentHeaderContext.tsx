import { PageHeaderTabs } from 'models/projects.models';
import React from 'react';
import { Subject } from 'rxjs';

export type BreadcrumbObservableValue = { [key: string]: PageHeaderTabs };

// TODO: Apply a better type to page header Observable, [Use a model for PageHeaderOptions]
export type EnvironmentHeaderContext = {
	breadcrumbObservable: Subject<BreadcrumbObservableValue | boolean>;
	pageHeaderObservable: Subject<any>;
};

export const breadcrumbObservable = new Subject<BreadcrumbObservableValue | boolean>();
export const pageHeaderObservable = new Subject<any>();

export const EnvironmentPageHeaderCtx = React.createContext<EnvironmentHeaderContext>({
	pageHeaderObservable,
	breadcrumbObservable,
});

export const usePageHeader = () => React.useContext(EnvironmentPageHeaderCtx);

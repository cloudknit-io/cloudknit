import { ZPageHeader } from 'components/molecules/page-header/PageHeader';
import React, { useContext, useEffect, useState } from 'react';

import { Breadcrumbs } from '../breadcrumb/Breadcrumbs';
import { usePageHeader } from '../contexts/EnvironmentHeaderContext';

export const EnvironmentHeader: React.FC = () => {
	const initialHeaderData = {
		breadcrumbs: [],
		headerTabs: [],
		pageName: '',
		filterTitle: '',
		onSearch: () => {},
		onViewChange: () => {},
		buttonText: '',
		initialView: '',
		checkBoxFilters: <></>,
		diffChecker: null
	};
	const { pageHeaderObservable } = usePageHeader();
	const [headerData, setHeaderData] = useState(initialHeaderData);
	useEffect(() => {
		const s = pageHeaderObservable.subscribe((data: any) => {
			if (data) setHeaderData(data);
		});
		return () => s.unsubscribe();
	});

	return (
		<>
			<div className="breadcrumbs-container">
				<Breadcrumbs />
			</div>
			<ZPageHeader
				breadcrumbs={[]}
				buttonText={headerData.buttonText}
				headerTabs={[]}
				pageName={headerData.pageName}
				filterTitle={''}
				initialView={headerData.initialView}
				onClickButton={(): void => {
					return;
				}}
				onViewChange={headerData.onViewChange}
				onSearch={headerData.onSearch}
				checkBoxFilters={headerData.checkBoxFilters}
				diffChecker={headerData.diffChecker}
			/>
		</>
	);
};

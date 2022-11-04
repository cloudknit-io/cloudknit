import './styles.scss';
import React, { useEffect } from 'react';
import { useState } from 'react';
import { FC } from 'react';

type Props = {
	id?: string;
	columns: any[];
	dataRows?: any[];
};

export const ErrorView: FC<Props> = ({ id, columns, dataRows = [] }) => {
	const [rows, setRows] = useState<any[]>(dataRows);

	useEffect(() => {
		setRows(dataRows);
	}, [dataRows]);

	return (
		<div className="zlifecycle-error-table">
			<div className="zlifecycle-error-table-column">
				{columns.map(c => (
					<span key={c.name}>{c.name}</span>
				))}
			</div>
			<div className="zlifecycle-error-table-rows">
				{rows.map(r => (
					<div className="zlifecycle-error-table-row" key={r.message}>
						{columns.map(c => (
							<span key={c.id}>{r[c.id]}</span>
						))}
					</div>
				))}
			</div>
		</div>
	);
};

import './styles.scss';
import 'react-virtualized/styles.css';

import classNames from 'classnames';
import { ZSyncStatus } from 'models/argo.models';
import React, { CSSProperties, FC, ReactNode } from 'react';
import { AutoSizer, Column, Table, TableRowProps, WindowScroller } from 'react-virtualized';

export interface TableColumn {
	id: string;
	name: string;
	combine?: boolean;
	render?: (data: any, index?: number) => string | number | ReactNode;
	textAlign?: string;
	width?: number;
}

interface Props {
	height?: number | string;
	table: {
		columns: TableColumn[];
		rows: any[];
		columnSections?: ReactNode;
		selectedRows?: number[];
	};
	data?: any;
	noData?: string | ReactNode;
	rowHeight?: number;
	rowClassName?: string;
	rowConditionalClass?: (data: any) => string;
	renderCustomRow?: (rowProps: TableRowProps) => ReactNode;
	cellClassName?: string;
	onRowClick?: (data: any) => void;
	onRowMouseOver?: () => void;
	onRowMouseOut?: () => void;
	scrollToIndex?: number;
}

export const transformAlignment = (align: string | undefined): string => {
	switch (align) {
		case 'left':
			return 'justify-content-start text-left';
		case 'right':
			return 'justify-content-end text-right';
		case 'center':
			return 'justify-content-center text-center';
		default:
			return 'justify-content-start text-left';
	}
};

export const ZTable: FC<Props> = (props: Props) => {
	const {
		table,
		data,
		noData,
		height,
		rowHeight,
		rowClassName,
		cellClassName,
		rowConditionalClass,
		onRowClick,
		onRowMouseOver,
		onRowMouseOut,
		scrollToIndex,
		renderCustomRow,
	} = props;

	const getWidth = (width: number, tableCols: TableColumn[]): number => {
		let customWidths = 0;
		tableCols.forEach(col => (customWidths += col.width || 0));

		return (width - customWidths) / tableCols.filter(col => col.width === undefined).length;
	};

	const rowRenderer = (rowProps: TableRowProps): ReactNode => {
		const {
			columns,
			rowData,
			index,
			onRowClick: onRowClickCustom,
			onRowDoubleClick,
			onRowMouseOver: onRowMouseOverCustom,
			onRowMouseOut: onRowMouseOutCustom,
			style,
		} = rowProps;

		if (renderCustomRow) {
			return renderCustomRow(rowProps);
		}

		return (
			<div
				onClick={(event): void => onRowClickCustom && onRowClickCustom({ rowData, index, event })}
				onDoubleClickCapture={(event): void => onRowDoubleClick && onRowDoubleClick({ rowData, index, event })}
				onMouseEnter={(event): void => onRowMouseOverCustom && onRowMouseOverCustom({ rowData, index, event })}
				onMouseLeave={(event): void => onRowMouseOutCustom && onRowMouseOutCustom({ rowData, index, event })}>
				<div
					{...rowProps}
					style={{
						...style,
						height: style.height - 8,
					}}>
					{columns}
				</div>
			</div>
		);
	};

	return (
		<WindowScroller>
			{({ height: heightScroller }): ReactNode => (
				<div style={{ height: height || heightScroller - 170 }}>
					<AutoSizer>
						{({ width, height: heightAuto }): ReactNode => (
							<>
								<Table
									scrollToIndex={scrollToIndex}
									data={data}
									height={table.rows.length === 0 ? 0 : heightAuto}
									width={width}
									headerHeight={40}
									overscanRowCount={10}
									rowClassName={({ index }): string =>
										classNames(
											rowClassName,
											rowConditionalClass && rowConditionalClass(table?.rows[index]),
											{
												'active-row': table?.selectedRows?.includes(table.rows[index]?.id),
												destroyed:
													table?.rows[index]?.componentStatus === ZSyncStatus.Destroyed,
											}
										)
									}
									rowStyle={({ index }): CSSProperties => ({
										backgroundColor: index === scrollToIndex ? '#ebfeff' : '',
									})}
									rowRenderer={rowRenderer}
									rowHeight={({ index }): number =>
										table.rows[index].section ? 50 : rowHeight || 72
									}
									onRowClick={({ rowData }): void => onRowClick && onRowClick(rowData)}
									onRowMouseOver={onRowMouseOver}
									onRowMouseOut={onRowMouseOut}
									rowGetter={({ index }): { [key: string]: string } => table.rows[index]}
									rowCount={table.rows.length}>
									{table.columns.map((column, columnIndex) => (
										<Column
											className={classNames(
												cellClassName,
												'flex',
												transformAlignment(column?.textAlign)
											)}
											headerClassName={`text-${column.textAlign}`}
											key={`table-col-${columnIndex}`}
											dataKey={column.id}
											label={column.name}
											cellRenderer={({ cellData, rowData, rowIndex }): ReactNode =>
												column.render
													? column.render(column.combine ? rowData : cellData, rowIndex)
													: cellData
											}
											width={column.width || getWidth(width, table.columns)}
										/>
									))}
								</Table>
								{table.rows.length === 0 && (
									<div
										className="text-center d-flex align-items-center justify-content-center mb-3 zlifecycle-table__no-data"
										style={{ width: width, height: heightAuto - 170 }}>
										<div>
											{noData ? <div>{noData}</div> : <p className="mt-5">No data available</p>}
										</div>
									</div>
								)}
							</>
						)}
					</AutoSizer>
				</div>
			)}
		</WindowScroller>
	);
};

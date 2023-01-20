import { DataSource, DataSourceOptions } from 'typeorm';
import { dbConfig } from '.';

export const AppDataSource: DataSource = new DataSource(
  dbConfig as DataSourceOptions
);

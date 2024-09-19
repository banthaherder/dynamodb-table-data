import { DataSourceJsonData } from '@grafana/data';
import { DataQuery } from '@grafana/schema';

export interface MyQuery extends DataQuery {
  tableName?: string;
  region?: string;
}

export interface DataSourceResponse {
  datapoints: [];
}

/**
 * These are options configured for each DataSource instance
 */
export interface MyDataSourceOptions extends DataSourceJsonData {
  defaultRegion?: string;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface MySecureJsonData {
  accessKeyId?: string;
  secretAccessKey?: string;
  sessionToken?: string;
}

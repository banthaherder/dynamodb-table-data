import { CoreApp, DataSourceInstanceSettings, ScopedVars } from '@grafana/data';
import { MyQuery, MyDataSourceOptions } from './types';
import { getTemplateSrv, DataSourceWithBackend } from '@grafana/runtime';

export class DataSource extends DataSourceWithBackend<MyQuery, MyDataSourceOptions> {
  jsonData: MyDataSourceOptions;

  constructor(instanceSettings: DataSourceInstanceSettings<MyDataSourceOptions>) {
    super(instanceSettings);

    this.jsonData = instanceSettings.jsonData;
  }

  getDefaultQuery(_: CoreApp): Partial<MyQuery> {
    return {
      region: this.jsonData.defaultRegion,
      tableName: '',
    };
  }

  applyTemplateVariables(query: MyQuery, scopedVars: ScopedVars) {
    return {
      ...query,
      // tableName: getTemplateSrv().replace(query.tableName, scopedVars),
    };
  }

  filterQuery(query: MyQuery): boolean {
    // Prevent the query from being executed if no tableName is provided
    return !!query.tableName;
  }
}

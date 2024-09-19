import React, { ChangeEvent } from 'react';
import { InlineField, Input, Stack } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from '../datasource';
import { MyDataSourceOptions, MyQuery } from '../types';

type Props = QueryEditorProps<DataSource, MyQuery, MyDataSourceOptions>;

export function QueryEditor({ query, onChange, onRunQuery, datasource }: Props) {

  const onTableNameChange = (event: ChangeEvent<HTMLInputElement>) => {
    onChange({ ...query, tableName: event.target.value });
  };

  const onRegionChange = (event: ChangeEvent<HTMLInputElement>) => {
    onChange({ ...query, region: event.target.value });
  };

  const runQuery = () => {
    // executes the query
    onRunQuery();
  };

  const { tableName, region } = query;
  const { defaultRegion } = datasource.jsonData;

  return (
    <Stack gap={0}>
      <InlineField label="Region">
        <Input
          id="query-editor-region"
          onChange={onRegionChange}
          onBlur={runQuery}
          value={region ? region : defaultRegion}
          width={16}
          type="string"
        />
      </InlineField>
      <InlineField label="Table Name">
        <Input
          id="query-editor-table-name"
          onChange={onTableNameChange}
          onBlur={runQuery}
          value={tableName}
          width={32}
          type="string"
        />
      </InlineField>
    </Stack>
  );
}

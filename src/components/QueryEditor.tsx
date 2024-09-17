import React, { ChangeEvent } from 'react';
import { InlineField, Input, Stack } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from '../datasource';
import { MyDataSourceOptions, MyQuery } from '../types';

type Props = QueryEditorProps<DataSource, MyQuery, MyDataSourceOptions>;

export function QueryEditor({ query, onChange, onRunQuery }: Props) {
  const onTableNameChange = (event: ChangeEvent<HTMLInputElement>) => {
    onChange({ ...query, tableName: event.target.value });
    // executes the query
    onRunQuery();
  };

  const { tableName } = query;

  return (
    <Stack gap={0}>
      <InlineField label="Table Name">
        <Input
          id="query-editor-table-name"
          onChange={onTableNameChange}
          value={tableName}
          width={32}
          type="string"
        />
      </InlineField>
    </Stack>
  );
}

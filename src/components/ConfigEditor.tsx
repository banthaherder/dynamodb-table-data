import React, { ChangeEvent } from 'react';
import { InlineField, SecretInput } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { MyDataSourceOptions, MySecureJsonData } from '../types';

interface Props extends DataSourcePluginOptionsEditorProps<MyDataSourceOptions, MySecureJsonData> {}

export function ConfigEditor(props: Props) {
  const { onOptionsChange, options } = props;
  const { secureJsonFields, secureJsonData } = options;

  // const onPathChange = (event: ChangeEvent<HTMLInputElement>) => {
  //   onOptionsChange({
  //     ...options,
  //     jsonData: {
  //       ...jsonData,
  //       path: event.target.value,
  //     },
  //   });
  // };

  // Secure field (only sent to the backend)
  const onAccessKeyIdChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      secureJsonData: {
        ...options.secureJsonData,
        accessKeyId: event.target.value,
      },
    });
  };

  const onResetAccessKeyId = () => {
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        accessKeyId: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        accessKeyId: '',
      },
    });
  };

  const onSecretAccessKeyChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      secureJsonData: {
        ...options.secureJsonData,
        secretAccessKey: event.target.value,
      },
    });
  };

  const onResetSecretAccessKey = () => {
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        secretAccessKey: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        secretAccessKey: '',
      },
    });
  };

  const onSessionTokenChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      secureJsonData: {
        ...options.secureJsonData,
        sessionToken: event.target.value,
      },
    });
  };

  const onResetSessionToken = () => {
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        sessionToken: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        sessionToken: '',
      },
    });
  };

  return (
    <>
      {/* <InlineField label="Table Name" labelWidth={25} interactive tooltip={'Json field returned to frontend'}>
        <Input
          id="config-editor-path"
          onChange={onPathChange}
          value={jsonData.tableName}
          placeholder="Enter the path, e.g. /api/v1"
          width={40}
        />
      </InlineField> */}
      <InlineField label="AWS Access Key ID" labelWidth={25} interactive tooltip={'Secure json field (backend only)'}>
        <SecretInput
          required
          id="config-editor-aws-access-key-id"
          isConfigured={secureJsonFields.accessKeyId}
          value={secureJsonData?.accessKeyId}
          placeholder="Enter your AWS Access Key ID"
          width={40}
          onReset={onResetAccessKeyId}
          onChange={onAccessKeyIdChange}
        />
      </InlineField>
      <InlineField label="AWS Secret Access Key" labelWidth={25} interactive tooltip={'Secure json field (backend only)'}>
        <SecretInput
          required
          id="config-editor-aws-secret-access-key"
          isConfigured={secureJsonFields.secretAccessKey}
          value={secureJsonData?.secretAccessKey}
          placeholder="Enter your AWS Secret Access Key"
          width={40}
          onReset={onResetSecretAccessKey}
          onChange={onSecretAccessKeyChange}
        />
      </InlineField>
      <InlineField label="AWS Session Token" labelWidth={25} interactive tooltip={'Secure json field (backend only)'}>
        <SecretInput
          required
          id="config-editor-aws-session-token"
          isConfigured={secureJsonFields.sessionToken}
          value={secureJsonData?.sessionToken}
          placeholder="Enter your AWS Session Token"
          width={40}
          onReset={onResetSessionToken}
          onChange={onSessionTokenChange}
        />
      </InlineField>
    </>
  );
}

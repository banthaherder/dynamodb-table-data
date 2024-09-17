package plugin

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/banthaherder/dynamodb-table-data/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

// Make sure Datasource implements required interfaces. This is important to do
// since otherwise we will only get a not implemented error response from plugin in
// runtime. In this example datasource instance implements backend.QueryDataHandler,
// backend.CheckHealthHandler interfaces. Plugin should not implement all these
// interfaces - only those which are required for a particular task.
var (
	_ backend.QueryDataHandler      = (*Datasource)(nil)
	_ backend.CheckHealthHandler    = (*Datasource)(nil)
	_ instancemgmt.InstanceDisposer = (*Datasource)(nil)
)

// QueryModel defines the structure of the query parameters from the front end
type QueryModel struct {
	TableName string `json:"tableName"`
	// TODO: add more fields here if needed (e.g., Filters, ProjectionExpression)
}

// NewDatasource creates a new datasource instance.
func NewDatasource(_ context.Context, _ backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	return &Datasource{}, nil
}

// Datasource is an example datasource which can respond to data queries, reports
// its health and has streaming skills.
type Datasource struct{}

// Dispose here tells plugin SDK that plugin wants to clean up resources when a new instance
// created. As soon as datasource settings change detected by SDK old datasource instance will
// be disposed and a new one will be created using NewSampleDatasource factory function.
func (d *Datasource) Dispose() {
	// Clean up datasource instance resources.
}

// createAWSSession creates an AWS session using the plugin settings
func createAWSSession(settings *models.PluginSettings) (*session.Session, error) {
	creds := credentials.NewStaticCredentials(
		settings.Secrets.AccessKeyId,
		settings.Secrets.SecretAccessKey,
		settings.Secrets.SessionToken,
	)

	awsConfig := &aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: creds,
	}

	return session.NewSession(awsConfig)
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifier).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (d *Datasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	// create response struct
	response := backend.NewQueryDataResponse()

	// loop over queries and execute them individually.
	for _, q := range req.Queries {
		res := d.query(ctx, req.PluginContext, q)

		// save the response in a hashmap
		// based on with RefID as identifier
		response.Responses[q.RefID] = res
	}

	return response, nil
}

type queryModel struct {
	TableName string `json:"tableName"`
}

// Field represents a structure holding field name, type, and values.
type Field struct {
	Name string
	// Type   string
	Values []string
}

// getFieldType determines the type of a DynamoDB attribute value.
func getFieldType(attr *dynamodb.AttributeValue) string {
	if attr.S != nil {
		return "S" // String
	} else if attr.N != nil {
		return "N" // Number
	} else if attr.B != nil {
		return "B" // Binary
	} else if attr.BOOL != nil {
		return "BOOL" // Boolean
	} else if attr.M != nil {
		return "M" // Map
	} else if attr.L != nil {
		return "L" // List
	} else if attr.NULL != nil && *attr.NULL {
		return "NULL" // Null
	}
	return "UNKNOWN"
}

// ExtractFields generates a slice of fields with their name, type, and values.
func ExtractFields(data []map[string]*dynamodb.AttributeValue) []Field {
	// Extract unique field names
	uniqueColumns := ExtractUniqueColumns(data)

	// Create a slice of Field structs
	fields := make([]Field, 0)

	for fieldName := range uniqueColumns {
		// Collect all values for this field
		values := make([]string, 0)
		for _, item := range data {
			if value, exists := item[fieldName]; exists {
				// Add the value as a string (for simplicity here; conversion may vary)
				if value.S != nil {
					values = append(values, *value.S)
				} else if value.N != nil {
					values = append(values, *value.N)
				} else if value.B != nil {
					values = append(values, string(value.B))
				} else {
					// If the value is not S, N, or B, use the placeholder
					values = append(values, "N/A")
				}
			} else {
				// If the field doesn't exist, use the placeholder
				values = append(values, "N/A")
			}
		}

		// Append the field information
		fields = append(fields, Field{
			Name: fieldName,
			// Type:   uniqueColumns[fieldName],
			Values: values,
		})
	}

	return fields
}

// Example ExtractUniqueColumns function (placeholder).
func ExtractUniqueColumns(data []map[string]*dynamodb.AttributeValue) map[string]string {
	uniqueColumns := make(map[string]string)

	for _, record := range data {
		for key, value := range record {
			if _, exists := uniqueColumns[key]; !exists {
				// Determine the type of the AttributeValue and store it.
				if value.S != nil {
					uniqueColumns[key] = "S"
				} else if value.N != nil {
					uniqueColumns[key] = "N"
				} else if value.B != nil {
					uniqueColumns[key] = "B"
				} else if value.BOOL != nil {
					uniqueColumns[key] = "BOOL"
				} else if value.M != nil {
					uniqueColumns[key] = "M"
				} else if value.L != nil {
					uniqueColumns[key] = "L"
				} else if value.NULL != nil && *value.NULL {
					uniqueColumns[key] = "NULL"
				}
			}
		}
	}

	return uniqueColumns
}

func (d *Datasource) query(_ context.Context, pCtx backend.PluginContext, query backend.DataQuery) backend.DataResponse {
	var response backend.DataResponse

	// Unmarshal the JSON into our queryModel.
	var qm queryModel

	err := json.Unmarshal(query.JSON, &qm)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("json unmarshal: %v", err.Error()))
	}

	// Load plugin settings
	settings, err := models.LoadPluginSettings(*pCtx.DataSourceInstanceSettings)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("failed to load plugin settings: %w", err.Error()))
	}

	// Create AWS session
	awsSession, err := createAWSSession(settings)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("failed to create AWS session: %w", err.Error()))
	}

	// Create DynamoDB client
	dynamoClient := dynamodb.New(awsSession)

	if err := json.Unmarshal(query.JSON, &qm); err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("failed to unmarshal query JSON: %w", err.Error()))
	}

	if qm.TableName == "" {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("tableName is required in the query", err.Error()))
	}

	// Perform DynamoDB scan with the table name from the query
	scanInput := &dynamodb.ScanInput{
		TableName: aws.String(qm.TableName),
	}

	scanOutput, err := dynamoClient.Scan(scanInput)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("failed to scan DynamoDB table %s: %w", qm.TableName, err.Error()))
	}

	fields := ExtractFields(scanOutput.Items)

	frame := data.NewFrame(query.RefID)

	// add the frames to the response.
	response.Frames = append(response.Frames, frame)

	// create data frame response.
	// For an overview on data frames and how grafana handles them:
	// https://grafana.com/developers/plugin-tools/introduction/data-frames
	for _, field := range fields {
		frame.Fields = append(frame.Fields,
			data.NewField(field.Name, nil, field.Values),
		)
	}

	// add the frames to the response.
	response.Frames = append(response.Frames, frame)

	return response
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (d *Datasource) CheckHealth(_ context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	res := &backend.CheckHealthResult{}
	config, err := models.LoadPluginSettings(*req.PluginContext.DataSourceInstanceSettings)

	if err != nil {
		res.Status = backend.HealthStatusError
		res.Message = "Unable to load settings"
		return res, nil
	}

	if config.Secrets.AccessKeyId == "" {
		res.Status = backend.HealthStatusError
		res.Message = "AWS Access Key ID is missing"
		return res, nil
	}

	if config.Secrets.SecretAccessKey == "" {
		res.Status = backend.HealthStatusError
		res.Message = "AWS Secret Access Key is missing"
		return res, nil
	}

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "Data source is working",
	}, nil
}

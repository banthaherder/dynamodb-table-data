package models

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type PluginSettings struct {
	// Region  string                `json:"region"`
	Secrets *SecretPluginSettings `json:"-"`
}

type SecretPluginSettings struct {
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	SessionToken    string `json:"sessionToken"`
}

func LoadPluginSettings(source backend.DataSourceInstanceSettings) (*PluginSettings, error) {
	settings := PluginSettings{}
	err := json.Unmarshal(source.JSONData, &settings)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal PluginSettings json: %w", err)
	}

	settings.Secrets = loadSecretPluginSettings(source.DecryptedSecureJSONData)

	return &settings, nil
}

func loadSecretPluginSettings(source map[string]string) *SecretPluginSettings {
	return &SecretPluginSettings{
		AccessKeyId:     source["accessKeyId"],
		SecretAccessKey: source["secretAccessKey"],
		SessionToken:    source["sessionToken"],
	}
}

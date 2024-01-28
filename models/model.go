package models

/*
 *  https://ibm.github.io/watson-machine-learning-sdk/_modules/ibm_watson_machine_learning/foundation_models/model.html#Model
 */

import (
	"fmt"
	"net/http"
)

type Model struct {
	url        string
	region     IBMCloudRegion
	apiVersion string

	ibmCloudAPIKey IBMCloudAPIKey
	projectID      WatsonxProjectID

	modelType ModelTypes

	token IAMToken

	httpClient Doer
}

func NewModel(apiKey, projectID string, options ...ModelOption) (*Model, error) {

	opts := defaulModelOptions()
	for _, opt := range options {
		if opt != nil {
			opt(opts)
		}
	}

	if opts.URL == "" {
		// User did not specify a URL, build it from the region
		opts.URL = buildBaseURL(opts.Region)
	}

	m := &Model{
		url:        opts.URL,
		region:     opts.Region,
		apiVersion: opts.APIVersion,

		ibmCloudAPIKey: apiKey,
		projectID:      projectID,

		modelType: opts.Model,

		// token: set below

		httpClient: &http.Client{},
	}

	err := m.RefreshToken()
	if err != nil {
		return nil, err
	}

	return m, nil
}

// CheckAndRefreshToken checks the IAM token if it expired; if it did, it refreshes it; nothing if not
func (m *Model) CheckAndRefreshToken() error {
	if m.token.Expired() {
		return m.RefreshToken()
	}
	return nil
}

// RefreshToken generates and sets the model with a new token
func (m *Model) RefreshToken() error {
	token, err := GenerateToken(m.httpClient, m.ibmCloudAPIKey)
	if err != nil {
		return err
	}
	m.token = token
	return nil
}

func buildBaseURL(region IBMCloudRegion) string {
	return fmt.Sprintf(BaseURLFormatStr, region)
}

func defaulModelOptions() *ModelOptions {
	return &ModelOptions{
		URL:        "",
		Region:     DefaultRegion,
		APIVersion: DefaultAPIVersion,
		Model:      DefaultModelType,
	}
}

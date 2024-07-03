package models

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

const (
	IAMCloudHost = "iam.cloud.ibm.com"
)

type Client struct {
	url        string
	iam        string
	region     IBMCloudRegion
	apiVersion string

	token     IAMToken
	apiKey    WatsonxAPIKey
	projectID WatsonxProjectID

	httpClient Doer
}

func NewClient(options ...ClientOption) (*Client, error) {

	opts := defaulClientOptions()
	for _, opt := range options {
		if opt != nil {
			opt(opts)
		}
	}

	if opts.URL == "" {
		// User did not specify a URL, build it from the region
		opts.URL = buildBaseURL(opts.Region)
	}

	if opts.IAM == "" {
		// User did not specify a IAM, use the default IAM cloud host
		opts.IAM = IAMCloudHost
	}

	if opts.apiKey == "" {
		return nil, errors.New("no watsonx API key provided")
	}

	if opts.projectID == "" {
		return nil, errors.New("no watsonx project ID provided")
	}

	m := &Client{
		url:        opts.URL,
		iam:        opts.IAM,
		region:     opts.Region,
		apiVersion: opts.APIVersion,

		// token: set below
		apiKey:    opts.apiKey,
		projectID: opts.projectID,

		httpClient: &http.Client{},
	}

	err := m.RefreshToken()
	if err != nil {
		return nil, err
	}

	return m, nil
}

// CheckAndRefreshToken checks the IAM token if it expired; if it did, it refreshes it; nothing if not
func (m *Client) CheckAndRefreshToken() error {
	if m.token.Expired() {
		return m.RefreshToken()
	}
	return nil
}

// RefreshToken generates and sets the model with a new token
func (m *Client) RefreshToken() error {
	token, err := GenerateToken(m.httpClient, m.apiKey, m.iam)
	if err != nil {
		return err
	}
	m.token = token
	return nil
}

func buildBaseURL(region IBMCloudRegion) string {
	return fmt.Sprintf(BaseURLFormatStr, region)
}

func defaulClientOptions() *ClientOptions {
	return &ClientOptions{
		URL:        os.Getenv(WatsonxURLEnvVarName),
		IAM:        os.Getenv(WatsonxIAMEnvVarName),
		Region:     DefaultRegion,
		APIVersion: DefaultAPIVersion,

		apiKey:    os.Getenv(WatsonxAPIKeyEnvVarName),
		projectID: os.Getenv(WatsonxProjectIDEnvVarName),
	}
}

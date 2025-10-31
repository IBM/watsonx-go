package models

import (
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/IBM/watsonx-go/pkg/auth"
	"github.com/IBM/watsonx-go/pkg/constants"
	"github.com/IBM/watsonx-go/pkg/http"
	"github.com/IBM/watsonx-go/pkg/types"
)

type Client struct {
	url        string
	apiVersion string

	token *auth.AuthToken
	auth  auth.Authenticator

	projectID constants.WatsonxProjectID
	spaceID   constants.WatsonxSpaceID

	httpClient types.Doer
}

func NewClient(options ...ClientOption) (*Client, error) {

	opts := defaultClientOptions()
	for _, opt := range options {
		if opt != nil {
			opt(opts)
		}
	}

	if opts.URL == "" {
		if opts.CPD != "" {
			opts.URL = opts.CPD
		} else {
			// User did not specify a URL, build it from the region
			opts.URL = buildBaseURL(opts.Region)
		}
	}

	if opts.IAM == "" && opts.CPD == "" {
		// User did not specify a IAM, use the default IAM cloud host
		opts.IAM = constants.DefaultIAMCloudHost
	}

	if opts.apiKey == "" && opts.cpdPassword == "" && opts.cpdAPIKey == "" {
		return nil, errors.New("no API key or password provided")
	}

	if opts.projectID == "" && opts.spaceID == "" {
		return nil, errors.New("no watsonx project ID or space ID provided")
	}

	if opts.projectID != "" && opts.spaceID != "" {
		return nil, errors.New("either project ID or space ID should be provided, not both")
	}

	m := &Client{
		url:        opts.URL,
		apiVersion: opts.APIVersion,

		// token and auth set below
		projectID: opts.projectID,
		spaceID:   opts.spaceID,

		httpClient: http.NewHttpClient(opts.disableTLSVerification),
	}

	var a auth.Authenticator
	var err error
	if opts.CPD != "" {
		a, err = auth.NewCPDAuthenticator(m.httpClient, opts.CPD, opts.CPDUsername, opts.cpdPassword, opts.cpdAPIKey)
	} else {
		a, err = auth.NewIAMAuthenticator(m.httpClient, opts.apiKey, opts.IAM)
	}
	if err != nil {
		return nil, err
	}
	m.auth = a

	if err := m.RefreshToken(); err != nil {
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
	token, err := m.auth.GenerateToken()
	if err != nil {
		return err
	}
	m.token = token
	return nil
}

// generateUrlFromEndpoint generates a URL from the endpoint and the client's configuration
func (m *Client) generateUrlFromEndpoint(endpoint string) string {
	params := url.Values{
		"version": {m.apiVersion},
	}

	generateTextURL := url.URL{
		Scheme:   "https",
		Host:     m.url,
		Path:     endpoint,
		RawQuery: params.Encode(),
	}

	return generateTextURL.String()
}

func buildBaseURL(region constants.IBMCloudRegion) string {
	return fmt.Sprintf(constants.BaseURLFormatStr, region)
}

func defaultClientOptions() *ClientOptions {
	return &ClientOptions{
		URL:        os.Getenv(constants.WatsonxURLEnvVarName),
		IAM:        os.Getenv(constants.WatsonxIAMEnvVarName),
		Region:     constants.DefaultRegion,
		APIVersion: constants.DefaultAPIVersion,

		cpdPassword: os.Getenv(constants.CPDPasswordEnvVarName),
		CPD:         os.Getenv(constants.CPDHostEnvVarName),
		CPDUsername: os.Getenv(constants.CPDUsernameEnvVarName),
		cpdAPIKey:   os.Getenv(constants.CPDAPIKeyEnvVarName),

		apiKey:    os.Getenv(constants.WatsonxAPIKeyEnvVarName),
		projectID: os.Getenv(constants.WatsonxProjectIDEnvVarName),
		spaceID:   os.Getenv(constants.WatsonxSpaceIDEnvVarName),

		disableTLSVerification: false,
	}
}

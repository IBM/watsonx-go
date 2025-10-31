package models

import (
	"github.com/IBM/watsonx-go/pkg/constants"
)

type ClientOption func(*ClientOptions)

type ClientOptions struct {
	URL string

	IAM    string
	Region constants.IBMCloudRegion

	CPD         string
	CPDUsername string
	cpdPassword string
	cpdAPIKey   string

	APIVersion string

	apiKey constants.WatsonxAPIKey

	projectID constants.WatsonxProjectID
	spaceID   constants.WatsonxSpaceID

	disableTLSVerification bool
}

func WithURL(url string) ClientOption {
	return func(o *ClientOptions) {
		o.URL = url
	}
}

func WithIAM(iamHost string) ClientOption {
	return func(o *ClientOptions) {
		o.IAM = iamHost
	}
}

func WithCPD(cpdHost string, cpdUsername string) ClientOption {
	return func(o *ClientOptions) {
		o.CPD = cpdHost
		o.CPDUsername = cpdUsername
	}
}

func WithCPDPassword(password string) ClientOption {
	return func(o *ClientOptions) {
		o.cpdPassword = password
	}
}

func WithCPDAPIKey(apiKey string) ClientOption {
	return func(o *ClientOptions) {
		o.cpdAPIKey = apiKey
	}
}

func WithRegion(region constants.IBMCloudRegion) ClientOption {
	return func(o *ClientOptions) {
		o.Region = region
	}
}

func WithAPIVersion(apiVersion string) ClientOption {
	return func(o *ClientOptions) {
		o.APIVersion = apiVersion
	}
}

func WithWatsonxAPIKey(watsonxAPIKey constants.WatsonxAPIKey) ClientOption {
	return func(o *ClientOptions) {
		o.apiKey = watsonxAPIKey
	}
}

func WithWatsonxProjectID(projectID constants.WatsonxProjectID) ClientOption {
	return func(o *ClientOptions) {
		o.projectID = projectID
	}
}

func WithWatsonxSpaceID(spaceID constants.WatsonxSpaceID) ClientOption {
	return func(o *ClientOptions) {
		o.spaceID = spaceID
	}
}

func WithDisableTLSVerification() ClientOption {
	return func(o *ClientOptions) {
		o.disableTLSVerification = true
	}
}

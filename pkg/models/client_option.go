package models

type ClientOption func(*ClientOptions)

type ClientOptions struct {
	URL        string
	IAM        string
	Region     IBMCloudRegion
	APIVersion string

	apiKey    WatsonxAPIKey
	projectID WatsonxProjectID

	retryConfig *RetryConfig
}

func WithURL(url string) ClientOption {
	return func(o *ClientOptions) {
		o.URL = url
	}
}

func WithIAM(iam string) ClientOption {
	return func(o *ClientOptions) {
		o.IAM = iam
	}
}

func WithRegion(region IBMCloudRegion) ClientOption {
	return func(o *ClientOptions) {
		o.Region = region
	}
}

func WithAPIVersion(apiVersion string) ClientOption {
	return func(o *ClientOptions) {
		o.APIVersion = apiVersion
	}
}

func WithWatsonxAPIKey(watsonxAPIKey WatsonxAPIKey) ClientOption {
	return func(o *ClientOptions) {
		o.apiKey = watsonxAPIKey
	}
}

func WithWatsonxProjectID(projectID WatsonxProjectID) ClientOption {
	return func(o *ClientOptions) {
		o.projectID = projectID
	}
}

func WithClientRetryConfig(retryConfig *RetryConfig) ClientOption {
	return func(o *ClientOptions) {
		o.retryConfig = retryConfig
	}
}

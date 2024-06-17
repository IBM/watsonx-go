package models

type ClientOption func(*ClientOptions)

type ClientOptions struct {
	URL        string
	Region     IBMCloudRegion
	APIVersion string

	apiKey    WatsonxAPIKey
	projectID WatsonxProjectID
}

func WithURL(url string) ClientOption {
	return func(o *ClientOptions) {
		o.URL = url
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

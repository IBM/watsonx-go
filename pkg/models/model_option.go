package models

type ModelOption func(*ModelOptions)

type ModelOptions struct {
	URL        string
	Region     IBMCloudRegion
	APIVersion string

	watsonxAPIKey WatsonxAPIKey
	projectID     WatsonxProjectID
}

func WithURL(url string) ModelOption {
	return func(o *ModelOptions) {
		o.URL = url
	}
}

func WithRegion(region IBMCloudRegion) ModelOption {
	return func(o *ModelOptions) {
		o.Region = region
	}
}

func WithAPIVersion(apiVersion string) ModelOption {
	return func(o *ModelOptions) {
		o.APIVersion = apiVersion
	}
}

func WithWatsonxAPIKey(watsonxAPIKey WatsonxAPIKey) ModelOption {
	return func(o *ModelOptions) {
		o.watsonxAPIKey = watsonxAPIKey
	}
}

func WithWatsonxProjectID(projectID WatsonxProjectID) ModelOption {
	return func(o *ModelOptions) {
		o.projectID = projectID
	}
}

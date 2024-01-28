package models

type ModelOption func(*ModelOptions)

type ModelOptions struct {
	URL        string
	Region     IBMCloudRegion
	APIVersion string

	ibmCloudAPIKey IBMCloudAPIKey
	projectID      WatsonxProjectID

	Model      ModelTypes
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

func WithIBMCloudAPIKey(ibmCloudAPIKey IBMCloudAPIKey) ModelOption {
	return func(o *ModelOptions) {
		o.ibmCloudAPIKey = ibmCloudAPIKey
	}
}

func WithWatsonxProjectID(projectID WatsonxProjectID) ModelOption {
	return func(o *ModelOptions) {
		o.projectID = projectID
	}
}

func WithModel(model ModelTypes) ModelOption {
	return func(o *ModelOptions) {
		o.Model = model
	}
}

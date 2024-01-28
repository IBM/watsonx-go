package models

type ModelOption func(*ModelOptions)

type ModelOptions struct {
	URL        string
	Region     IBMCloudRegion
	APIVersion string
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

func WithModel(model ModelTypes) ModelOption {
	return func(o *ModelOptions) {
		o.Model = model
	}
}

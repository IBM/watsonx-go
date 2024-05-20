package models

import (
	"net/http"
)

/*
 *  https://ibm.github.io/watson-machine-learning-sdk/model.html#ibm_watson_machine_learning.foundation_models.utils.enums.ModelTypes
 */

type (
	WatsonxAPIKey    = string
	WatsonxProjectID = string
	IBMCloudRegion   = string
	ModelType        = string
)

const (
	WatsonxAPIKeyEnvVarName    = "WATSONX_API_KEY"
	WatsonxProjectIDEnvVarName = "WATSONX_PROJECT_ID"

	US_South  IBMCloudRegion = "us-south"
	Dallas    IBMCloudRegion = US_South
	EU_DE     IBMCloudRegion = "eu-de"
	Frankfurt IBMCloudRegion = EU_DE
	JP_TOK    IBMCloudRegion = "jp-tok"
	Tokyo     IBMCloudRegion = JP_TOK

	DefaultRegion     = US_South
	BaseURLFormatStr  = "%s.ml.cloud.ibm.com" // Need to call SPrintf on it with region
	DefaultAPIVersion = "2024-05-20"
)

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

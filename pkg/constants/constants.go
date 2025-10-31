package constants

type (
	WatsonxAPIKey    = string
	WatsonxProjectID = string
	WatsonxSpaceID   = string
	IBMCloudRegion   = string
	ModelType        = string
)

const (
	WatsonxURLEnvVarName = "WATSONX_URL_HOST" // Override the default URL host '*.ml.cloud.ibm.com'
	WatsonxIAMEnvVarName = "WATSONX_IAM_HOST" // Override the default IAM host 'iam.cloud.ibm.com'

	CPDHostEnvVarName     = "WATSONX_CPD_HOST"
	CPDUsernameEnvVarName = "WATSONX_CPD_USERNAME"
	CPDPasswordEnvVarName = "WATSONX_CPD_PASSWORD"
	CPDAPIKeyEnvVarName   = "WATSONX_CPD_API_KEY"

	WatsonxAPIKeyEnvVarName    = "WATSONX_API_KEY"
	WatsonxProjectIDEnvVarName = "WATSONX_PROJECT_ID"
	WatsonxSpaceIDEnvVarName   = "WATSONX_SPACE_ID"

	US_South  IBMCloudRegion = "us-south"
	Dallas    IBMCloudRegion = US_South
	EU_DE     IBMCloudRegion = "eu-de"
	Frankfurt IBMCloudRegion = EU_DE
	JP_TOK    IBMCloudRegion = "jp-tok"
	Tokyo     IBMCloudRegion = JP_TOK

	DefaultRegion       = US_South
	BaseURLFormatStr    = "%s.ml.cloud.ibm.com" // Need to call SPrintf on it with region
	DefaultAPIVersion   = "2024-05-20"
	DefaultIAMCloudHost = "iam.cloud.ibm.com"
)

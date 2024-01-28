package models

import (
	"net/http"
)

/*
 *  https://ibm.github.io/watson-machine-learning-sdk/model.html#ibm_watson_machine_learning.foundation_models.utils.enums.ModelTypes
 */

type (
	IBMCloudAPIKey   = string
	WatsonxProjectID = string
	IBMCloudRegion   = string

	ModelTypes      = string
	DecodingMethods = string
)

const (
	IBMCloudAPIKeyEnvVarName   = "IBMCLOUD_API_KEY"
	WatsonxProjectIDEnvVarName = "WATSONX_PROJECT_ID"

	US_South  IBMCloudRegion = "us-south"
	Dallas    IBMCloudRegion = US_South
	EU_DE     IBMCloudRegion = "eu-de"
	Frankfurt IBMCloudRegion = EU_DE
	JP_TOK    IBMCloudRegion = "jp-tok"
	Tokyo     IBMCloudRegion = JP_TOK

	DefaultRegion     = US_South
	BaseURLFormatStr  = "%s.ml.cloud.ibm.com" // Need to call SPrintf on it with region
	DefaultAPIVersion = "2023-05-02"

	// https://ibm.github.io/watson-machine-learning-sdk/_modules/ibm_watson_machine_learning/foundation_models/utils/enums.html#ModelTypes
	FLAN_T5_XXL             ModelTypes = "google/flan-t5-xxl"
	FLAN_UL2                ModelTypes = "google/flan-ul2"
	MT0_XXL                 ModelTypes = "bigscience/mt0-xxl"
	GPT_NEOX                ModelTypes = "eleutherai/gpt-neox-20b"
	MPT_7B_INSTRUCT2        ModelTypes = "ibm/mpt-7b-instruct2"
	STARCODER               ModelTypes = "bigcode/starcoder"
	LLAMA_2_70B_CHAT        ModelTypes = "meta-llama/llama-2-70b-chat"
	LLAMA_2_13B_CHAT        ModelTypes = "meta-llama/llama-2-13b-chat"
	GRANITE_13B_INSTRUCT    ModelTypes = "ibm/granite-13b-instruct-v1"
	GRANITE_13B_CHAT        ModelTypes = "ibm/granite-13b-chat-v1"
	FLAN_T5_XL              ModelTypes = "google/flan-t5-xl"
	GRANITE_13B_CHAT_V2     ModelTypes = "ibm/granite-13b-chat-v2"
	GRANITE_13B_INSTRUCT_V2 ModelTypes = "ibm/granite-13b-instruct-v2"

	// https://ibm.github.io/watson-machine-learning-sdk/_modules/ibm_watson_machine_learning/foundation_models/utils/enums.html#DecodingMethods
	Sample DecodingMethods = "sample"
	Greedy DecodingMethods = "greedy"

	DefaultModelType      = FLAN_T5_XL
	DefaultDecodingMethod = Greedy
)

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

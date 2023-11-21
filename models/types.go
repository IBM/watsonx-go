package models

import (
	"fmt"
)

/*
 *  https://ibm.github.io/watson-machine-learning-sdk/model.html#ibm_watson_machine_learning.foundation_models.utils.enums.ModelTypes
 */

type ModelType string

type ModelTypes struct {
	FLAN_T5_XXL          ModelType
	FLAN_UL2             ModelType
	MT0_XXL              ModelType
	GPT_NEOX             ModelType
	MPT_7B_INSTRUCT2     ModelType
	STARCODER            ModelType
	LLAMA_2_70B_CHAT     ModelType
	GRANITE_13B_INSTRUCT ModelType
	GRANITE_13B_CHAT     ModelType
}

var ModelTypesEnum = ModelTypes{
	FLAN_T5_XXL:          "google/flan-t5-xxl",
	FLAN_UL2:             "google/flan-ul2",
	MT0_XXL:              "bigscience/mt0-xxl",
	GPT_NEOX:             "eleutherai/gpt-neox-20b",
	MPT_7B_INSTRUCT2:     "ibm/mpt-7b-instruct2",
	STARCODER:            "bigcode/starcoder",
	LLAMA_2_70B_CHAT:     "meta-llama/llama-2-70b-chat",
	GRANITE_13B_INSTRUCT: "ibm/granite-13b-instruct-v1",
	GRANITE_13B_CHAT:     "ibm/granite-13b-chat-v1",
}

func GetDefaultModelType() ModelType {
	return ModelTypesEnum.FLAN_UL2
}

/*
 *  https://ibm.github.io/watson-machine-learning-sdk/model.html#ibm_watson_machine_learning.foundation_models.utils.enums.DecodingMethods
 */

type DecodingMethod string

type DecodingMethods struct {
	SAMPLE DecodingMethod
	GREEDY DecodingMethod
}

var DecodingMethodsEnum = DecodingMethods{
	SAMPLE: "sample",
	GREEDY: "greedy",
}

/*
 *  https://ibm.github.io/watson-machine-learning-sdk/model.html#metanames.GenTextParamsMetaNames
 */

type LengthPenalty struct {
	DecayFactor float64 `json:"decay_factor"`
	StartIndex  uint    `json:"start_index"`
}

type GenParams struct {
	DecodingMethod      *DecodingMethod `json:"decoding_methods,omitempty"`
	LengthPenalty       *LengthPenalty  `json:"length_penalty,omitempty"`
	Temperature         *float64        `json:"temperature,omitempty"`
	TopP                *float64        `json:"top_p,omitempty"`
	TopK                *uint           `json:"top_k,omitempty"`
	RandomSeed          *uint           `json:"random_seed,omitempty"`
	RepetitionPenalty   *uint           `json:"repetition_penalty,omitempty"`
	MinNewTokens        *uint           `json:"min_new_tokens,omitempty"`
	MaxNewTokens        *uint           `json:"max_new_tokens,omitempty"`
	StopSequences       *[]string       `json:"stop_sequences,omitempty"`
	TimeLimit           *uint           `json:"time_limit,omitempty"`
	TruncateInputTokens *uint           `json:"truncate_input_tokens,omitempty"`
	// ReturnOptions       string`json:"return_options,omitempty"`
}

func DefaultGenParams() *GenParams {
	return &GenParams{
		DecodingMethod:      nil,
		LengthPenalty:       nil,
		Temperature:         nil,
		TopP:                nil,
		TopK:                nil,
		RandomSeed:          nil,
		RepetitionPenalty:   nil,
		MinNewTokens:        nil,
		MaxNewTokens:        nil,
		StopSequences:       nil,
		TimeLimit:           nil,
		TruncateInputTokens: nil,
		// returnOptions:       nil,
	}
}

func printValueOrDefault[T any](ptr *T) string {
	if ptr == nil {
		return "Default"
	}
	return fmt.Sprintf("%v", *ptr)
}

func (gp GenParams) String() string {
	return fmt.Sprintf(""+
		"DecodingMethod:        %v\n"+
		"LengthPenalty:         %v\n"+
		"Temperature:           %v\n"+
		"TopP:                  %v\n"+
		"TopK:                  %v\n"+
		"RandomSeed:            %v\n"+
		"RepetitionPenalty:     %v\n"+
		"MinNewTokens:          %v\n"+
		"MaxNewTokens:          %v\n"+
		"StopSequences:         %v\n"+
		"TimeLimit:             %v\n"+
		"TruncateInputTokens:   %v\n",
		printValueOrDefault(gp.DecodingMethod),
		printValueOrDefault(gp.LengthPenalty),
		printValueOrDefault(gp.Temperature),
		printValueOrDefault(gp.TopP),
		printValueOrDefault(gp.TopK),
		printValueOrDefault(gp.RandomSeed),
		printValueOrDefault(gp.RepetitionPenalty),
		printValueOrDefault(gp.MinNewTokens),
		printValueOrDefault(gp.MaxNewTokens),
		printValueOrDefault(gp.StopSequences),
		printValueOrDefault(gp.TimeLimit),
		printValueOrDefault(gp.TruncateInputTokens),
	)
}

type Credentials struct {
	ApiKey      string
	Url         string
	AccessToken string `json:"access_token"`
}

func (credentials *Credentials) setDefaultUrl() {
	credentials.Url = "https://us-south.ml.cloud.ibm.com"
}

package foundation_models

/*
https://ibm.github.io/watson-machine-learning-sdk/model.html#ibm_watson_machine_learning.foundation_models.utils.enums.ModelTypes
*/
type ModelType string

const (
	FLAN_T5_XXL          ModelType = "google/flan-t5-xxl"
	FLAN_UL2             ModelType = "google/flan-ul2"
	MT0_XXL              ModelType = "bigscience/mt0-xxl"
	GPT_NEOX             ModelType = "eleutherai/gpt-neox-20b"
	MPT_7B_INSTRUCT2     ModelType = "ibm/mpt-7b-instruct2"
	STARCODER            ModelType = "bigcode/starcoder"
	LLAMA_2_70B_CHAT     ModelType = "meta-llama/llama-2-70b-chat"
	GRANITE_13B_INSTRUCT ModelType = "ibm/granite-13b-instruct-v1"
	GRANITE_13B_CHAT     ModelType = "ibm/granite-13b-chat-v1"
)

func GetDefaultModelType() ModelType {
	return FLAN_UL2
}

/*
https://ibm.github.io/watson-machine-learning-sdk/model.html#ibm_watson_machine_learning.foundation_models.utils.enums.DecodingMethods
*/
type DecodingMethods string

const (
	SAMPLE DecodingMethods = "sample"
	GREEDY DecodingMethods = "greedy"
)

/*
https://ibm.github.io/watson-machine-learning-sdk/model.html#metanames.GenTextParamsMetaNames
*/
type GenParams struct {
	DecodingMethod      DecodingMethods
	LengthPenalty       map[string]interface{}
	Temperature         float64
	TopP                float64
	TopK                uint
	RandomSeed          uint
	RepetitionPenalty   uint
	MinNewTokens        uint
	MaxNewTokens        uint
	StopSequences       []string
	TimeLimit           uint
	TruncateInputTokens uint
	// ReturnOptions       string
}

func NewGenParams(config map[string]interface{}) *GenParams {
	metaNames := &GenParams{
		DecodingMethod:      SAMPLE,
		LengthPenalty:       map[string]interface{}{"decay_factor": 2.5, "start_index": 5},
		Temperature:         0.5,
		TopP:                0.2,
		TopK:                1,
		RandomSeed:          33,
		RepetitionPenalty:   2,
		MinNewTokens:        50,
		MaxNewTokens:        200,
		StopSequences:       []string{"fail"},
		TimeLimit:           600000,
		TruncateInputTokens: 200,
		// returnOptions: [],
	}

	if config == nil {
		return metaNames
	}

	// Check and set values from the configuration
	if value, ok := config["DecodingMethod"]; ok {
		metaNames.DecodingMethod = value.(DecodingMethods)
	}
	if value, ok := config["LengthPenalty"]; ok {
		metaNames.LengthPenalty = value.(map[string]interface{})
	}
	if value, ok := config["Temperature"]; ok {
		metaNames.Temperature = value.(float64)
	}
	if value, ok := config["TopP"]; ok {
		metaNames.TopP = value.(float64)
	}
	if value, ok := config["TopK"]; ok {
		metaNames.TopK = value.(uint)
	}
	if value, ok := config["RandomSeed"]; ok {
		metaNames.RandomSeed = value.(uint)
	}
	if value, ok := config["RepetitionPenalty"]; ok {
		metaNames.RepetitionPenalty = value.(uint)
	}
	if value, ok := config["MinNewTokens"]; ok {
		metaNames.MinNewTokens = value.(uint)
	}
	if value, ok := config["MaxNewTokens"]; ok {
		metaNames.MaxNewTokens = value.(uint)
	}
	if value, ok := config["StopSequences"]; ok {
		// Assuming the configuration provides an array or slice of strings
		if stopSequences, ok := value.([]string); ok {
			metaNames.StopSequences = stopSequences
		}
	}
	if value, ok := config["TimeLimit"]; ok {
		metaNames.TimeLimit = value.(uint)
	}
	if value, ok := config["TruncateInputTokens"]; ok {
		metaNames.TruncateInputTokens = value.(uint)
	}
	// if value, ok := config["ReturnOptions"]; ok {
	// 	// Assuming the configuration provides a map of strings to booleans
	// 	if returnOptions, ok := value.(map[string]interface{}); ok {
	// 		metaName.ReturnOptions = returnOptions
	// 	}
	// }

	return metaNames
}

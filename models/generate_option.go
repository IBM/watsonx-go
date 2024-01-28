package models

import (
	"fmt"
)

type GenerateOption func(*GenerateOptions)

type LengthPenalty struct {
	DecayFactor float64 `json:"decay_factor"`
	StartIndex  uint    `json:"start_index"`
}

type ReturnOptions struct {
	InputText       bool `json:"input_text"`
	GeneratedTokens bool `json:"generated_tokens"`
	InputTokens     bool `json:"input_tokens"`
	TokenLogProbs   bool `json:"token_logprobs"`
	TokenRanks      bool `json:"token_ranks"`
	TopNTokens      int  `json:"top_n_tokens"`
}

type GenerateOptions struct {
	// https://ibm.github.io/watson-machine-learning-sdk/_modules/metanames.html#GenTextParamsMetaNames
	DecodingMethod      *string        `json:"decoding_method,omitempty"`
	LengthPenalty       *LengthPenalty `json:"length_penalty,omitempty"`
	Temperature         *float64       `json:"temperature,omitempty"`
	TopP                *float64       `json:"top_p,omitempty"`
	TopK                *uint          `json:"top_k,omitempty"`
	RandomSeed          *uint          `json:"random_seed,omitempty"`
	RepetitionPenalty   *float64       `json:"repetition_penalty,omitempty"`
	MinNewTokens        *uint          `json:"min_new_tokens,omitempty"`
	MaxNewTokens        *uint          `json:"max_new_tokens,omitempty"`
	StopSequences       *[]string      `json:"stop_sequences,omitempty"`
	TimeLimit           *uint          `json:"time_limit,omitempty"`
	TruncateInputTokens *uint          `json:"truncate_input_tokens,omitempty"`
	ReturnOptions       *ReturnOptions `json:"return_options,omitempty"`
}

func WithDecodingMethod(decodingMethod string) GenerateOption {
	return func(opts *GenerateOptions) {
		opts.DecodingMethod = &decodingMethod
	}
}

func WithLengthPenalty(decayFactor float64, startIndex uint) GenerateOption {
	return func(opts *GenerateOptions) {
		opts.LengthPenalty = &LengthPenalty{decayFactor, startIndex}
	}
}

func WithTemperature(temperature float64) GenerateOption {
	return func(opts *GenerateOptions) {
		opts.Temperature = &temperature
	}
}

func WithTopP(topP float64) GenerateOption {
	return func(opts *GenerateOptions) {
		opts.TopP = &topP
	}
}

func WithTopK(topK uint) GenerateOption {
	return func(opts *GenerateOptions) {
		opts.TopK = &topK
	}
}

func WithRandomSeed(randomSeed uint) GenerateOption {
	return func(opts *GenerateOptions) {
		opts.RandomSeed = &randomSeed
	}
}

func WithRepetitionPenalty(repetitionPenalty uint) GenerateOption {
	return func(opts *GenerateOptions) {
		opts.RepetitionPenalty = &repetitionPenalty
	}
}

func WithMinNewTokens(minNewTokens uint) GenerateOption {
	return func(opts *GenerateOptions) {
		opts.MinNewTokens = &minNewTokens
	}
}

func WithMaxNewTokens(maxNewTokens uint) GenerateOption {
	return func(opts *GenerateOptions) {
		opts.MaxNewTokens = &maxNewTokens
	}
}

func WithStopSequences(stopSequences []string) GenerateOption {
	return func(opts *GenerateOptions) {
		opts.StopSequences = &stopSequences
	}
}

func WithTimeLimit(timeLimit uint) GenerateOption {
	return func(opts *GenerateOptions) {
		opts.TimeLimit = &timeLimit
	}
}

func WithTruncateInputTokens(truncateInputTokens uint) GenerateOption {
	return func(opts *GenerateOptions) {
		opts.TruncateInputTokens = &truncateInputTokens
	}
}

func WithReturnOptions(inputText, generatedTokens, inputTokens, tokenLogProbs, tokenRanks bool, topNTokens int) GenerateOption {
	return func(opts *GenerateOptions) {
		opts.ReturnOptions = &ReturnOptions{inputText, generatedTokens, inputTokens, tokenLogProbs, tokenRanks, topNTokens}
	}
}

func (gp *GenerateOptions) String() string {
	return fmt.Sprintf(
		"decodingMethod: %v\n"+
			"lengthPenalty: %v\n"+
			"temperature: %v\n"+
			"topP: %v\n"+
			"topK: %v\n"+
			"randomSeed: %v\n"+
			"repetitionPenalty: %v\n"+
			"minNewTokens: %v\n"+
			"maxNewTokens: %v\n"+
			"stopSequences: %v\n"+
			"timeLimit: %v\n"+
			"truncateInputTokens: %v\n"+
			"returnOptions: %v",
		gp.DecodingMethod,
		gp.LengthPenalty,
		gp.Temperature,
		gp.TopP,
		gp.TopK,
		gp.RandomSeed,
		gp.RepetitionPenalty,
		gp.MinNewTokens,
		gp.MaxNewTokens,
		gp.StopSequences,
		gp.TimeLimit,
		gp.TruncateInputTokens,
		gp.ReturnOptions,
	)
}

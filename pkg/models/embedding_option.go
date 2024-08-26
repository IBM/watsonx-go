package models

import "fmt"

type EmbeddingOption func(*EmbeddingOptions)

type EmbeddingOptions struct {
	TruncateInputTokens *uint                   `json:"truncate_input_tokens,omitempty"`
	ReturnOptions       *EmbeddingReturnOptions `json:"return_options,omitempty"`
}

type EmbeddingReturnOptions struct {
	InputText bool `json:"input_text"`
}

func WithEmbeddingTruncateInputTokens(truncateInputTokens uint) EmbeddingOption {
	return func(opts *EmbeddingOptions) {
		opts.TruncateInputTokens = &truncateInputTokens
	}
}

func WithEmbeddingReturnOptions(inputText bool) EmbeddingOption {
	return func(opts *EmbeddingOptions) {
		opts.ReturnOptions = &EmbeddingReturnOptions{inputText}
	}
}

func (ep *EmbeddingOptions) String() string {
	return fmt.Sprintf(
		"truncateInputTokens: %v\n"+
			"returnOptions: %v\n",
		ep.TruncateInputTokens,
		ep.ReturnOptions,
	)
}

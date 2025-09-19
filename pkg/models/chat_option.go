package models

// ChatOption defines the function signature for chat configuration options
type ChatOption func(*ChatOptions)

// ChatOptions holds all the configurable parameters for chat requests
type ChatOptions struct {
	Tools               []ChatTool          `json:"tools,omitempty"`
	ToolChoiceOption    *string             `json:"tool_choice_option,omitempty"`
	ToolChoice          *ChatToolChoice     `json:"tool_choice,omitempty"`
	Context             *string             `json:"context,omitempty"`
	MaxTokens           *uint               `json:"max_tokens,omitempty"`
	MaxCompletionTokens *uint               `json:"max_completion_tokens,omitempty"`
	Temperature         *float64            `json:"temperature,omitempty"`
	TopP                *float64            `json:"top_p,omitempty"`
	FrequencyPenalty    *float64            `json:"frequency_penalty,omitempty"`
	PresencePenalty     *float64            `json:"presence_penalty,omitempty"`
	Stop                []string            `json:"stop,omitempty"`
	N                   *uint               `json:"n,omitempty"`
	ResponseFormat      *ChatResponseFormat `json:"response_format,omitempty"`
	Seed                *int                `json:"seed,omitempty"`
	TimeLimit           *uint               `json:"time_limit,omitempty"`
	LogitBias           map[string]float64  `json:"logit_bias,omitempty"`
	LogProbs            *bool               `json:"logprobs,omitempty"`
	TopLogProbs         *uint               `json:"top_logprobs,omitempty"`
}

// WithChatTools sets the tools available for the chat completion
func WithChatTools(tools ...ChatTool) ChatOption {
	return func(opts *ChatOptions) {
		opts.Tools = tools
	}
}

// WithChatToolChoice sets how the model should use tools
func WithChatToolChoice(choice string) ChatOption {
	return func(opts *ChatOptions) {
		opts.ToolChoiceOption = &choice
	}
}

// WithChatToolChoiceFunction forces the model to use a specific function
func WithChatToolChoiceFunction(name string) ChatOption {
	return func(opts *ChatOptions) {
		opts.ToolChoice = &ChatToolChoice{
			Type:     "function",
			Function: &ChatToolChoiceFunction{Name: name},
		}
	}
}

// WithChatContext sets additional context for the chat completion
func WithChatContext(context string) ChatOption {
	return func(opts *ChatOptions) {
		opts.Context = &context
	}
}

// WithChatMaxTokens sets the maximum number of tokens to generate
func WithChatMaxTokens(maxTokens uint) ChatOption {
	return func(opts *ChatOptions) {
		opts.MaxTokens = &maxTokens
	}
}

// WithChatMaxCompletionTokens sets the maximum number of completion tokens
func WithChatMaxCompletionTokens(maxTokens uint) ChatOption {
	return func(opts *ChatOptions) {
		opts.MaxCompletionTokens = &maxTokens
	}
}

// WithChatTemperature sets the sampling temperature
func WithChatTemperature(temperature float64) ChatOption {
	return func(opts *ChatOptions) {
		opts.Temperature = &temperature
	}
}

// WithChatTopP sets the nucleus sampling parameter
func WithChatTopP(topP float64) ChatOption {
	return func(opts *ChatOptions) {
		opts.TopP = &topP
	}
}

// WithChatFrequencyPenalty sets the frequency penalty
func WithChatFrequencyPenalty(penalty float64) ChatOption {
	return func(opts *ChatOptions) {
		opts.FrequencyPenalty = &penalty
	}
}

// WithChatPresencePenalty sets the presence penalty
func WithChatPresencePenalty(penalty float64) ChatOption {
	return func(opts *ChatOptions) {
		opts.PresencePenalty = &penalty
	}
}

// WithChatStop sets the stop sequences
func WithChatStop(stop ...string) ChatOption {
	return func(opts *ChatOptions) {
		opts.Stop = stop
	}
}

// WithChatN sets the number of chat completion choices to generate
func WithChatN(n uint) ChatOption {
	return func(opts *ChatOptions) {
		opts.N = &n
	}
}

// WithChatJSONMode enables JSON mode for the response
func WithChatJSONMode() ChatOption {
	return func(opts *ChatOptions) {
		opts.ResponseFormat = &ChatResponseFormat{Type: "json_object"}
	}
}

// WithChatJSONSchema sets a JSON schema for structured output
func WithChatJSONSchema(schema interface{}) ChatOption {
	return func(opts *ChatOptions) {
		opts.ResponseFormat = &ChatResponseFormat{
			Type:       "json_schema",
			JSONSchema: schema,
		}
	}
}

// WithChatSeed sets the random seed for reproducible outputs
func WithChatSeed(seed int) ChatOption {
	return func(opts *ChatOptions) {
		opts.Seed = &seed
	}
}

// WithChatTimeLimit sets the time limit for the request
func WithChatTimeLimit(timeLimit uint) ChatOption {
	return func(opts *ChatOptions) {
		opts.TimeLimit = &timeLimit
	}
}

// WithChatLogitBias sets the logit bias for token selection
func WithChatLogitBias(logitBias map[string]float64) ChatOption {
	return func(opts *ChatOptions) {
		opts.LogitBias = logitBias
	}
}

// WithChatLogProbs enables or disables log probabilities in the response
func WithChatLogProbs(logProbs bool) ChatOption {
	return func(opts *ChatOptions) {
		opts.LogProbs = &logProbs
	}
}

// WithChatTopLogProbs sets the number of most likely tokens to return at each position
func WithChatTopLogProbs(topLogProbs uint) ChatOption {
	return func(opts *ChatOptions) {
		opts.TopLogProbs = &topLogProbs
	}
}

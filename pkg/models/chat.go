package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

// API endpoints
const (
	ChatEndpoint       string = "/ml/v1/text/chat"
	ChatStreamEndpoint string = "/ml/v1/text/chat_stream"
)

// Chat response body
const (
	errorChatBodyBufferSize = 1024
)

// Message roles
const (
	RoleSystem    = "system"
	RoleUser      = "user"
	RoleAssistant = "assistant"
	RoleTool      = "tool"
)

// Chat message content types
type ChatMessageContent struct {
	Type     string        `json:"type"` // "text", "image_url", etc.
	Text     *string       `json:"text,omitempty"`
	ImageURL *ChatImageURL `json:"image_url,omitempty"`
}

type ChatImageURL struct {
	URL    string `json:"url"`
	Detail string `json:"detail,omitempty"` // "low", "high", "auto"
}

// Chat message structure that handles both string and array content formats
type ChatMessage struct {
	Role       string                  `json:"role"`
	Content    ChatMessageContentUnion `json:"content"`
	Name       *string                 `json:"name,omitempty"`
	ToolCalls  []ChatToolCall          `json:"tool_calls,omitempty"`
	ToolCallID *string                 `json:"tool_call_id,omitempty"`
}

// ChatMessageContentUnion handles both string and array formats
type ChatMessageContentUnion struct {
	StringContent *string              `json:"-"`
	ArrayContent  []ChatMessageContent `json:"-"`
}

// UnmarshalJSON implements custom unmarshaling to handle both string and array content
func (c *ChatMessageContentUnion) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as string first
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		c.StringContent = &str
		return nil
	}

	// Try to unmarshal as array
	var arr []ChatMessageContent
	if err := json.Unmarshal(data, &arr); err == nil {
		c.ArrayContent = arr
		return nil
	}

	return fmt.Errorf("content must be either string or array of ChatMessageContent")
}

// MarshalJSON implements custom marshaling
func (c *ChatMessageContentUnion) MarshalJSON() ([]byte, error) {
	if c.StringContent != nil {
		return json.Marshal(*c.StringContent)
	}
	if c.ArrayContent != nil {
		return json.Marshal(c.ArrayContent)
	}
	return json.Marshal([]ChatMessageContent{})
}

// GetText returns the text content regardless of format
func (c *ChatMessageContentUnion) GetText() string {
	if c.StringContent != nil {
		return *c.StringContent
	}
	if len(c.ArrayContent) > 0 && c.ArrayContent[0].Text != nil {
		return *c.ArrayContent[0].Text
	}
	return ""
}

// ToArray converts the content to array format
func (c *ChatMessageContentUnion) ToArray() []ChatMessageContent {
	if c.ArrayContent != nil {
		return c.ArrayContent
	}
	if c.StringContent != nil {
		return []ChatMessageContent{
			{
				Type: "text",
				Text: c.StringContent,
			},
		}
	}
	return []ChatMessageContent{}
}

// Tool definitions
type ChatTool struct {
	Type     string           `json:"type"` // "function"
	Function ChatToolFunction `json:"function"`
}

type ChatToolFunction struct {
	Name        string      `json:"name"`
	Description *string     `json:"description,omitempty"`
	Parameters  interface{} `json:"parameters,omitempty"` // JSON schema
}

// Tool calls in responses
type ChatToolCall struct {
	ID       string               `json:"id"`
	Type     string               `json:"type"` // "function"
	Function ChatToolCallFunction `json:"function"`
}

type ChatToolCallFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"` // JSON string
}

// Tool choice options
type ChatToolChoice struct {
	Type     string                  `json:"type"` // "function"
	Function *ChatToolChoiceFunction `json:"function,omitempty"`
}

type ChatToolChoiceFunction struct {
	Name string `json:"name"`
}

// Response format options
type ChatResponseFormat struct {
	Type       string      `json:"type"` // "text", "json_object", "json_schema"
	JSONSchema interface{} `json:"json_schema,omitempty"`
}

// Chat completion request
type ChatRequest struct {
	ModelID             string              `json:"model_id"`
	Messages            []ChatMessage       `json:"messages"`
	ProjectID           *string             `json:"project_id,omitempty"`
	SpaceID             *string             `json:"space_id,omitempty"`
	Tools               []ChatTool          `json:"tools,omitempty"`
	ToolChoiceOption    *string             `json:"tool_choice_option,omitempty"` // "auto", "none", "required"
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
	Stream              *bool               `json:"stream,omitempty"`
	ResponseFormat      *ChatResponseFormat `json:"response_format,omitempty"`
	Seed                *int                `json:"seed,omitempty"`
	TimeLimit           *uint               `json:"time_limit,omitempty"`
}

// Chat completion response
type ChatResponse struct {
	ID           string         `json:"id"`
	ModelID      string         `json:"model_id"`
	Created      int64          `json:"created"`
	CreatedAt    *time.Time     `json:"created_at,omitempty"`
	Choices      []ChatChoice   `json:"choices"`
	Usage        *ChatUsage     `json:"usage,omitempty"`
	ModelVersion *string        `json:"model_version,omitempty"`
	System       *SystemDetails `json:"system,omitempty"`
}

type ChatChoice struct {
	Index        int           `json:"index"`
	Message      *ChatMessage  `json:"message,omitempty"`
	Delta        *ChatMessage  `json:"delta,omitempty"` // For streaming
	FinishReason *string       `json:"finish_reason,omitempty"`
	LogProbs     *ChatLogProbs `json:"logprobs,omitempty"`
}

type ChatUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ChatLogProbs struct {
	Content []ChatContentLogProbs `json:"content,omitempty"`
	Refusal []ChatContentLogProbs `json:"refusal,omitempty"`
}

type ChatContentLogProbs struct {
	Token       string           `json:"token"`
	LogProb     float64          `json:"logprob"`
	Bytes       []int            `json:"bytes,omitempty"`
	TopLogProbs []ChatTopLogProb `json:"top_logprobs,omitempty"`
}

type ChatTopLogProb struct {
	Token   string  `json:"token"`
	LogProb float64 `json:"logprob"`
	Bytes   []int   `json:"bytes,omitempty"`
}

// SystemDetails represents system information from the response
type SystemDetails struct {
	Warnings interface{} `json:"warnings,omitempty"`
}

const ChatMessageTypeText = "text"

// CreateChatMessage creates a text message with the specified role and content
func CreateChatMessage(role string, content ChatMessageContentUnion) ChatMessage {
	return ChatMessage{
		Role:    role,
		Content: content,
	}
}

// CreateFunction creates a tool function definition
func CreateFunction(name, description string, parameters interface{}) ChatTool {
	return ChatTool{
		Type: "function",
		Function: ChatToolFunction{
			Name:        name,
			Description: &description,
			Parameters:  parameters,
		},
	}
}

// CreateToolMessage creates a tool response message
func CreateToolMessage(toolCallID, content string) ChatMessage {
	return ChatMessage{
		Role: RoleTool,
		Content: ChatMessageContentUnion{
			ArrayContent: []ChatMessageContent{
				{
					Type: ChatMessageTypeText,
					Text: &content,
				},
			},
		},
		ToolCallID: &toolCallID,
	}
}

// CreateSystemMessage creates a system message
func CreateSystemMessage(content string) ChatMessage {
	stringContent := ChatMessageContentUnion{
		StringContent: &content,
	}
	return CreateChatMessage(RoleSystem, stringContent)
}

// CreateUserMessage creates a user message
func CreateUserMessage(content string) ChatMessage {
	arrayContent := ChatMessageContentUnion{
		ArrayContent: []ChatMessageContent{
			{
				Type: ChatMessageTypeText,
				Text: &content,
			},
		},
	}

	return CreateChatMessage(RoleUser, arrayContent)
}

// CreateAssistantMessage creates an assistant message
func CreateAssistantMessage(content string) ChatMessage {
	arrayContent := ChatMessageContentUnion{
		ArrayContent: []ChatMessageContent{
			{
				Type: ChatMessageTypeText,
				Text: &content,
			},
		},
	}
	return CreateChatMessage(RoleAssistant, arrayContent)
}

// Chat generates a text chat based on messages and parameters
func (c *Client) Chat(modelID string, messages []ChatMessage, options ...ChatOption) (ChatResponse, error) {
	// Validate input
	if modelID == "" {
		return ChatResponse{}, errors.New("modelID cannot be empty")
	}

	if len(messages) == 0 {
		return ChatResponse{}, errors.New("messages cannot be empty")
	}

	// Apply options
	opts := &ChatOptions{}
	for _, opt := range options {
		if opt != nil {
			opt(opts)
		}
	}

	// Build the request payload
	payload := c.BuildChatRequest(modelID, messages, opts)

	// Make the API request
	response, err := c.generateChatRequest(payload)
	if err != nil {
		return ChatResponse{}, err
	}

	// Validate response
	if len(response.Choices) == 0 {
		return ChatResponse{}, errors.New("no choices received in response")
	}

	return response, nil
}

// SimpleChat provides a simple interface for single-turn text chat conversations
func (c *Client) SimpleChat(modelID, prompt string, options ...ChatOption) (string, error) {
	messages := []ChatMessage{
		CreateUserMessage(prompt),
	}

	response, err := c.Chat(modelID, messages, options...)
	if err != nil {
		return "", err
	}

	if len(response.Choices) == 0 || response.Choices[0].Message == nil {
		return "", errors.New("no response received")
	}

	choice := response.Choices[0]
	if choice.Message == nil {
		return "", errors.New("no message in response")
	}

	text := choice.Message.Content.GetText()
	if text == "" {
		return "", errors.New("no text content in response")
	}

	return text, nil
}

// BuildChatRequest constructs the ChatRequest payload
func (c *Client) BuildChatRequest(modelID string, messages []ChatMessage, opts *ChatOptions) ChatRequest {
	// Use the project ID from the client (already configured during client creation)
	projectID := string(c.projectID)

	payload := ChatRequest{
		ModelID:             modelID,
		Messages:            messages,
		ProjectID:           &projectID,
		Tools:               opts.Tools,
		ToolChoiceOption:    opts.ToolChoiceOption,
		ToolChoice:          opts.ToolChoice,
		Context:             opts.Context,
		MaxTokens:           opts.MaxTokens,
		MaxCompletionTokens: opts.MaxCompletionTokens,
		Temperature:         opts.Temperature,
		TopP:                opts.TopP,
		FrequencyPenalty:    opts.FrequencyPenalty,
		PresencePenalty:     opts.PresencePenalty,
		Stop:                opts.Stop,
		N:                   opts.N,
		ResponseFormat:      opts.ResponseFormat,
		Seed:                opts.Seed,
		TimeLimit:           opts.TimeLimit,
	}

	return payload
}

// generateChatRequest sends a request to the chat endpoint
func (c *Client) generateChatRequest(payload ChatRequest) (ChatResponse, error) {
	// Ensure we have a valid token
	err := c.CheckAndRefreshToken()
	if err != nil {
		return ChatResponse{}, fmt.Errorf("failed to refresh token: %w", err)
	}

	// Build the URL using the client's configuration
	chatURL := c.generateUrlFromEndpoint(ChatEndpoint)

	// Marshal the payload to JSON
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return ChatResponse{}, fmt.Errorf("failed to marshal request payload: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest(http.MethodPost, chatURL, bytes.NewReader(payloadJSON))
	if err != nil {
		return ChatResponse{}, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token.value)

	// Execute the request using the client's HTTP client with retry
	res, err := c.httpClient.DoWithRetry(req)
	if err != nil {
		return ChatResponse{}, err
	}
	defer func() {
		if cerr := res.Body.Close(); cerr != nil {
			log.Println("error closing response body: ", cerr)
		}
	}()

	// Check for successful status code
	if res.StatusCode != http.StatusOK {
		//Get the detailed error response from Watson X API
		return ChatResponse{}, DecodeWatsonxError(res)
	}

	// Decode the response
	var chatRes ChatResponse
	if err := json.NewDecoder(res.Body).Decode(&chatRes); err != nil {
		return ChatResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return chatRes, nil
}

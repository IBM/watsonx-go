package test

import (
	"strings"
	"testing"

	wx "github.com/IBM/watsonx-go/pkg/models"
)

const (
	modelChatLlama3Integration  = "meta-llama/llama-3-3-70b-instruct"
	modelChatInvalidIntegration = "invalid-chat-model-test"
)

func TestChatSingleMessage(t *testing.T) {
	client := getClient(t)

	messages := []wx.ChatMessage{
		wx.CreateUserMessage("What is the capital of France?"),
	}

	response, err := client.Chat(modelChatLlama3Integration, messages)

	if err != nil {
		t.Fatalf("Expected no error for chat request, but got %v", err)
	}

	if len(response.Choices) == 0 {
		t.Fatal("Expected at least one choice in response")
	}

	if response.Choices[0].Message == nil {
		t.Fatal("Expected message in first choice")
	}

	content := response.Choices[0].Message.Content.GetText()
	if content == "" {
		t.Fatal("Expected non-empty response content")
	}

	if response.ModelID == "" {
		t.Fatal("Expected model ID to be set in response")
	}

	if response.Usage == nil {
		t.Fatal("Expected usage information in response")
	}
}

func TestChatMultipleMessages(t *testing.T) {
	client := getClient(t)

	messages := []wx.ChatMessage{
		wx.CreateSystemMessage("You are a helpful assistant that answers in one word."),
		wx.CreateUserMessage("What is 2+2?"),
	}

	response, err := client.Chat(
		modelChatLlama3Integration,
		messages,
		wx.WithChatTemperature(0.3),
		wx.WithChatMaxTokens(10),
	)

	if err != nil {
		t.Fatalf("Expected no error for chat request, but got %v", err)
	}

	if len(response.Choices) == 0 {
		t.Fatal("Expected at least one choice in response")
	}

	content := response.Choices[0].Message.Content.GetText()
	if content == "" {
		t.Fatal("Expected non-empty response content")
	}

	if response.Usage.TotalTokens == 0 {
		t.Fatal("Expected non-zero total tokens")
	}

	if response.Usage.CompletionTokens > 3 {
		t.Fatal("Expected completion tokens to be less than or equal to 3")
	}

	lowerContent := strings.ToLower(content)
	if !strings.Contains(lowerContent, "four") && !strings.Contains(lowerContent, "4") {
		t.Fatal("Expected response to contain 'four' or '4'")
	}
}

func TestSimpleChat(t *testing.T) {
	client := getClient(t)

	response, err := client.SimpleChat(
		modelChatLlama3Integration,
		"Say 'hello' in French in one word",
		wx.WithChatTemperature(0.1),
	)

	if err != nil {
		t.Fatalf("Expected no error for simple chat, but got %v", err)
	}

	if response == "" {
		t.Fatal("Expected non-empty response")
	}

	lowerResponse := strings.ToLower(response)
	if !strings.Contains(lowerResponse, "bonjour") {
		t.Fatal("Expected response to contain 'bonjour'")
	}
}

func TestChatWithJSONMode(t *testing.T) {
	client := getClient(t)

	messages := []wx.ChatMessage{
		wx.CreateSystemMessage("You are a helpful assistant that responds in JSON format."),
		wx.CreateUserMessage("What is the capital of France? Respond with {\"capital\": \"city_name\"}"),
	}

	response, err := client.Chat(
		modelChatLlama3Integration,
		messages,
		wx.WithChatJSONMode(),
		wx.WithChatTemperature(0.1),
	)

	if err != nil {
		t.Fatalf("Expected no error for JSON mode chat, but got %v", err)
	}

	if len(response.Choices) == 0 {
		t.Fatal("Expected at least one choice in response")
	}

	content := response.Choices[0].Message.Content.GetText()
	if content == "" {
		t.Fatal("Expected non-empty response content")
	}
}

func TestChatInvalidModel(t *testing.T) {
	client := getClient(t)

	messages := []wx.ChatMessage{
		wx.CreateUserMessage("Hello"),
	}

	_, err := client.Chat(modelChatInvalidIntegration, messages)

	if err == nil {
		t.Fatal("Expected error for invalid model but got nil")
	}
}

func TestChatInvalidParameter(t *testing.T) {
	client := getClient(t)

	messages := []wx.ChatMessage{
		wx.CreateUserMessage("Hello"),
	}

	_, err := client.Chat(
		modelChatLlama3Integration,
		messages,
		wx.WithChatTemperature(100), // Invalid temperature - should be 0-2
	)

	if err == nil {
		t.Fatal("Expected error for invalid parameter but got nil")
	}
}

func TestChatInvalidModelReturnsWatsonxError(t *testing.T) {
	client := getClient(t)

	messages := []wx.ChatMessage{
		wx.CreateUserMessage("Hello"),
	}

	_, err := client.Chat(modelChatInvalidIntegration, messages)
	if err == nil {
		t.Fatal("Expected error for invalid model but got nil")
	}

	wxErr, ok := err.(*wx.WatsonxError)
	if !ok {
		t.Fatalf("Expected error type *WatsonxError, got %T", err)
	}

	if wxErr.StatusCode == 0 {
		t.Fatal("Expected StatusCode to be set in WatsonxError")
	}

	if len(wxErr.Errors) == 0 {
		t.Fatal("Expected error details to be present")
	}

	if wxErr.Errors[0].Code == "" {
		t.Fatal("Expected error code to be populated")
	}

	if wxErr.Errors[0].Message == "" {
		t.Fatal("Expected error message to be populated")
	}
}

func TestChatInvalidParameterReturnsWatsonxError(t *testing.T) {
	client := getClient(t)

	messages := []wx.ChatMessage{
		wx.CreateUserMessage("Hello"),
	}

	_, err := client.Chat(
		modelChatLlama3Integration,
		messages,
		wx.WithChatTemperature(100), // Invalid temperature - should be 0-2
	)

	if err == nil {
		t.Fatal("Expected error for invalid parameter but got nil")
	}

	wxErr, ok := err.(*wx.WatsonxError)
	if !ok {
		t.Fatalf("Expected error type *WatsonxError, got %T", err)
	}

	if wxErr.StatusCode == 0 {
		t.Fatal("Expected StatusCode to be set")
	}

	if len(wxErr.Errors) == 0 {
		t.Fatal("Expected structured error details")
	}
}

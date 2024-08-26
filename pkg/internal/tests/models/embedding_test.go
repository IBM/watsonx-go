package test

import (
	wx "github.com/IBM/watsonx-go/pkg/models"
	"reflect"
	"testing"
)

const (
	EmbeddingModelId        = "ibm/slate-30m-english-rtrvr"
	EmbeddingModelDimension = 384
)

func TestEmbeddingSingleQuery(t *testing.T) {
	client := getClient(t)

	text := "Hello, world!"

	response, err := client.EmbedQuery(EmbeddingModelId, text)

	if err != nil {
		t.Fatalf("Expected no error for embedding query, but got %v", err)
	}

	if len(response.Results) != 1 {
		t.Fatalf("Expected 1 embedding in response, but got %d", len(response.Results))
	}

	if len(response.Results[0].Embedding) != EmbeddingModelDimension {
		t.Fatalf("Expected dimension of %d, but got %d", EmbeddingModelDimension, len(response.Results[0].Embedding))
	}

	if response.Results[0].Input != "" {
		t.Fatalf("Expected input to be empty, but got %s", response.Results[0].Input)
	}

	if response.Model != EmbeddingModelId {
		t.Fatalf("Expected model to be %s, but got %s", EmbeddingModelId, response.Model)
	}
}

func TestEmbeddingSingleQueryWithOptions(t *testing.T) {
	client := getClient(t)

	text := "Hello, world!"

	response, err := client.EmbedQuery(
		EmbeddingModelId,
		text,
		wx.WithEmbeddingTruncateInputTokens(2),
		wx.WithEmbeddingReturnOptions(true),
	)

	if err != nil {
		t.Fatalf("Expected no error for embedding query, but got %v", err)
	}

	if len(response.Results) != 1 {
		t.Fatalf("Expected 1 embedding in response, but got %d", len(response.Results))
	}

	if len(response.Results[0].Embedding) != EmbeddingModelDimension {
		t.Fatalf("Expected dimension of %d, but got %d", EmbeddingModelDimension, len(response.Results[0].Embedding))
	}

	if response.Results[0].Input != text {
		t.Fatalf("Expected input to be %s, but got %s", text, response.Results[0].Input)
	}

	if response.Model != EmbeddingModelId {
		t.Fatalf("Expected model to be %s, but got %s", EmbeddingModelId, response.Model)
	}

	// the embedding must be different from the one without truncate options
	responseNoOptions, err := client.EmbedQuery(EmbeddingModelId, text, wx.WithEmbeddingReturnOptions(true))
	if err != nil {
		t.Fatalf("Expected no error for embedding query, but got %v", err)
	}

	if reflect.DeepEqual(response.Results[0].Embedding, responseNoOptions.Results[0].Embedding) {
		t.Fatalf("Expected different embeddings with and without options, but got the same")
	}
}

func TestEmbeddingMultipleQueries(t *testing.T) {
	client := getClient(t)

	texts := []string{"Hello, world!", "How are you?"}

	response, err := client.EmbedDocuments(EmbeddingModelId, texts, wx.WithEmbeddingReturnOptions(true))

	if err != nil {
		t.Fatalf("Expected no error for embedding queries, but got %v", err)
	}

	if len(response.Results) != len(texts) {
		t.Fatalf("Expected %d embeddings in response, but got %d", len(texts), len(response.Results))
	}

	if len(response.Results) != 2 {
		t.Fatalf("Expected 2 embeddings in response, but got %d", len(response.Results))
	}

	for i, result := range response.Results {
		if len(result.Embedding) != EmbeddingModelDimension {
			t.Fatalf("Expected dimension of %d, but got %d for query %d", EmbeddingModelDimension, len(result.Embedding), i)
		}

		if result.Input != texts[i] {
			t.Fatalf("Expected input to be %s, but got %s for query %d", texts[i], result.Input, i)
		}
	}

	if response.Model != EmbeddingModelId {
		t.Fatalf("Expected model to be %s, but got %s", EmbeddingModelId, response.Model)
	}
}

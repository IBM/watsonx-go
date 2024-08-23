package test

import (
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

	if response.Model != EmbeddingModelId {
		t.Fatalf("Expected model to be %s, but got %s", EmbeddingModelId, response.Model)
	}
}

func TestEmbeddingMultipleQueries(t *testing.T) {
	client := getClient(t)

	texts := []string{"Hello, world!", "How are you?"}

	response, err := client.EmbedDocuments(EmbeddingModelId, texts)

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
	}

	if response.Model != EmbeddingModelId {
		t.Fatalf("Expected model to be %s, but got %s", EmbeddingModelId, response.Model)
	}
}

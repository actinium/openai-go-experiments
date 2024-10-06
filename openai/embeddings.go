package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
)

const (
	embeddingsModel      = "text-embedding-ada-002"
	embeddingsDimensions = 1536
)

type EmbeddingsClient struct {
	openaiClient *OpenAIClient
}

func NewEmbeddingsClient(client *OpenAIClient) *EmbeddingsClient {
	return &EmbeddingsClient{
		openaiClient: client,
	}
}

type Embedding struct {
	Text      string
	Embedding [embeddingsDimensions]float32
}

type embeddingsRequestPayload struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type embeddingsResponsePayload struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string                        `json:"object"`
		Embedding [embeddingsDimensions]float32 `json:"embedding"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens uint `json:"prompt_tokens"`
		TotalTokens  uint `json:"total_tokens"`
	} `json:"usage"`
}

func (client *EmbeddingsClient) CreateWithContext(ctx context.Context, input string) (Embedding, error) {
	requestPayload, err := json.Marshal(embeddingsRequestPayload{
		Model: embeddingsModel,
		Input: input,
	})
	if err != nil {
		return Embedding{}, err
	}

	body, err := client.openaiClient.post(ctx, "/v1/embeddings", bytes.NewReader(requestPayload))
	if err != nil {
		return Embedding{}, err
	}

	defer body.Close()
	content, err := io.ReadAll(body)
	if err != nil {
		return Embedding{}, err
	}

	var responsePayload embeddingsResponsePayload
	err = json.Unmarshal(content, &responsePayload)
	if err != nil {
		return Embedding{}, err
	}

	return Embedding{
		Text:      input,
		Embedding: responsePayload.Data[0].Embedding,
	}, nil

}

func (client *EmbeddingsClient) Create(input string) (Embedding, error) {
	ctx := context.Background()

	return client.CreateWithContext(ctx, input)
}

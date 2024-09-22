package openai

import (
	"context"
	"encoding/json"
	"io"
)

type modelsResponsePayload struct {
	Object string `json:"object"`
	Data   []struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Created int64  `json:"created"`
		OwnedBy string `json:"owned_by"`
	} `json:"data"`
}

func (openAIClient *OpenAIClient) Models() ([]string, error) {
	var models []string

	body, err := openAIClient.get(context.Background(), "/models")
	if err != nil {
		return []string{}, err
	}

	defer body.Close()
	content, err := io.ReadAll(body)
	if err != nil {
		return []string{}, err
	}

	var responsePayload modelsResponsePayload
	err = json.Unmarshal(content, &responsePayload)
	if err != nil {
		return []string{}, err
	}

	for _, model := range responsePayload.Data {
		models = append(models, model.ID)
	}

	return models, nil
}

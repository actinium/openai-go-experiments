package openai

import (
	"context"
	"io"
	"net/http"
	"time"
)

const (
	baseUrl        = "https://api.openai.com/v1"
	requestTimeout = 1 * time.Minute
)

type OpenAIClient struct {
	apiKey  string
	baseUrl string
	timeout time.Duration
}

func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{
		apiKey:  apiKey,
		baseUrl: baseUrl,
		timeout: requestTimeout,
	}
}

func (client *OpenAIClient) http() *http.Client {
	return &http.Client{
		Timeout: client.timeout,
	}
}

func (client *OpenAIClient) post(
	ctx context.Context,
	path string,
	payload io.Reader,
) (io.ReadCloser, error) {
	req, err := http.NewRequest("POST", client.baseUrl+path, payload)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	req.Header.Add("Authorization", "Bearer "+client.apiKey)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	res, err := client.http().Do(req)
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}

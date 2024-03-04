package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
)

const (
	ImageSize256  = "256x256"
	ImageSize512  = "512x512"
	ImageSize1024 = "1024x1024"
)

const (
	ResponseFormatUrl    = "url"
	ResponseFormatBase64 = "b64_json"
)

type DalleOptions struct {
	Model          string
	N              uint
	Size           string
	ResponseFormat string
}

var DefaultDalleOptions = DalleOptions{
	Model:          "dall-e-2",
	N:              1,
	Size:           ImageSize1024,
	ResponseFormat: ResponseFormatUrl,
}

type DalleClient struct {
	openaiClient *OpenAIClient
	options      DalleOptions
}

func NewDalleClient(client *OpenAIClient, options DalleOptions) *DalleClient {
	return &DalleClient{
		openaiClient: client,
		options:      options,
	}
}

type dalleRequestPayload struct {
	Model          string `json:"model"`
	Prompt         string `json:"prompt"`
	N              uint   `json:"n"`
	Size           string `json:"size"`
	ResponseFormat string `json:"response_format"`
}

type dalleUrlResponsePayload struct {
	CreatedAt uint `json:"created"`
	Data      []struct {
		Url string `json:"url"`
	} `json:"data"`
}

type dalleBase64ResponsePayload struct {
	CreatedAt uint `json:"created"`
	Data      []struct {
		Base64 string `json:"b64_json"`
	} `json:"data"`
}

func (dalle *DalleClient) GenerateImage(ctx context.Context, prompt string) ([]string, error) {
	requestPayload, err := json.Marshal(dalleRequestPayload{
		Model:          dalle.options.Model,
		Prompt:         prompt,
		N:              dalle.options.N,
		Size:           dalle.options.Size,
		ResponseFormat: dalle.options.ResponseFormat,
	})
	if err != nil {
		return nil, err
	}

	body, err := dalle.openaiClient.post(ctx, "/images/generations", bytes.NewReader(requestPayload))
	if err != nil {
		return nil, err
	}

	defer body.Close()
	content, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	if dalle.options.ResponseFormat == ResponseFormatUrl {
		var responsePayload dalleUrlResponsePayload
		err = json.Unmarshal(content, &responsePayload)
		if err != nil {
			return nil, err
		}

		urls := make([]string, len(responsePayload.Data))
		for i, data := range responsePayload.Data {
			urls[i] = data.Url
		}

		return urls, nil
	} else {
		var responsePayload dalleBase64ResponsePayload
		err = json.Unmarshal(content, &responsePayload)
		if err != nil {
			return nil, err
		}

		base64s := make([]string, len(responsePayload.Data))
		for i, data := range responsePayload.Data {
			base64s[i] = data.Base64
		}

		return base64s, nil
	}

}

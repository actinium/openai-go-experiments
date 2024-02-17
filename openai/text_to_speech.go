package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
)

var TextToSpeechVoices = [...]string{"alloy", "echo", "fable", "onyx", "nova", "shimmer"}
var TextToSpeechFormats = [...]string{"mp3", "opus", "aac", "flac"}

type TextToSpeechOption struct {
	Voice  string
	HD     bool
	Speed  float32
	Format string
}

var DefaultTextToSpeechOption = TextToSpeechOption{
	Voice:  "alloy",
	HD:     false,
	Speed:  1.0,
	Format: "mp3",
}

type TextToSpeechClient struct {
	openaiClient *OpenAIClient
	options      TextToSpeechOption
}

func NewTextToSpeechClient(openaiClient *OpenAIClient, options TextToSpeechOption) *TextToSpeechClient {
	return &TextToSpeechClient{
		openaiClient: openaiClient,
		options:      options,
	}
}

func (client *TextToSpeechClient) SetVoice(voice string) {
	client.options.Voice = voice
}

func (client *TextToSpeechClient) SetHD(hd bool) {
	client.options.HD = hd
}

func (options TextToSpeechOption) model() string {
	if options.HD {
		return "tts-1-hd"
	}
	return "tts-1"

}

type Audio struct {
	Data   []byte
	Format string
}

type textToSpeechRequestPayload struct {
	Input string  `json:"input"`
	Model string  `json:"model"`
	Voice string  `json:"voice"`
	Speed float32 `json:"speed"`
}

func (client *TextToSpeechClient) GenerateAudioWithContext(ctx context.Context, text string) (*Audio, error) {
	requestPayload, err := json.Marshal(textToSpeechRequestPayload{
		Input: text,
		Model: client.options.model(),
		Voice: client.options.Voice,
		Speed: float32(client.options.Speed),
	})
	if err != nil {
		return nil, err
	}

	body, err := client.openaiClient.post(ctx, "/audio/speech", bytes.NewReader(requestPayload))
	if err != nil {
		return nil, err
	}
	defer body.Close()

	content, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	return &Audio{
		Data:   content,
		Format: client.options.Format,
	}, nil
}

func (client *TextToSpeechClient) GenerateAudio(text string) (*Audio, error) {
	ctx := context.Background()

	return client.GenerateAudioWithContext(ctx, text)
}

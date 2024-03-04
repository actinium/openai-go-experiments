package openai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strings"
)

type ChatOptions struct {
	Model       string
	Temperature float32
	MaxTokens   uint
}

var DefaultChatOptions = ChatOptions{
	Model:       "gpt-3.5-turbo",
	Temperature: 1.0,
	MaxTokens:   1024,
}

type ChatClient struct {
	openaiClient *OpenAIClient
	options      ChatOptions
}

func NewChatClient(client *OpenAIClient, options ChatOptions) *ChatClient {
	return &ChatClient{
		openaiClient: client,
		options:      options,
	}
}

func (client *ChatClient) UseModel(model string) {
	client.options.Model = model
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Chat struct {
	openaiClient *OpenAIClient
	options      ChatOptions
	Messages     []Message
}

func (client *ChatClient) NewChat() *Chat {
	return &Chat{
		openaiClient: client.openaiClient,
		options:      client.options,
	}
}

func (client *ChatClient) NewChatWithSystemPrompt(systemPrompt string) *Chat {
	chat := Chat{
		openaiClient: client.openaiClient,
		options:      client.options,
	}

	chat.AppendMessage("system", systemPrompt)

	return &chat
}

func (chat *Chat) AppendMessage(role string, message string) {
	chat.Messages = append(chat.Messages, Message{Role: role, Content: message})
}

type chatRequestPayload struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float32   `json:"temperature"`
	MaxTokens   uint      `json:"max_tokens"`
	Stream      bool      `json:"stream"`
}

type chatResponsePayload struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	CreatedAt uint   `json:"created"`
	Model     string `json:"model"`
	Choices   []struct {
		Index        uint    `json:"index"`
		Message      Message `json:"message"`
		FinishReason string  `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     uint `json:"prompt_tokens"`
		CompletionTokens uint `json:"completion_tokens"`
		TotalTokens      uint `json:"total_tokens"`
	} `json:"usage"`
}

func (chat *Chat) MakeRequest(ctx context.Context) (Message, error) {
	requestPayload, err := json.Marshal(chatRequestPayload{
		Model:       chat.options.Model,
		Messages:    chat.Messages,
		Temperature: chat.options.Temperature,
		MaxTokens:   uint(chat.options.MaxTokens),
		Stream:      false,
	})
	if err != nil {
		return Message{}, err
	}

	body, err := chat.openaiClient.post(ctx, "/chat/completions", bytes.NewReader(requestPayload))
	if err != nil {
		return Message{}, err
	}

	defer body.Close()
	content, err := io.ReadAll(body)
	if err != nil {
		return Message{}, err
	}

	var responsePayload chatResponsePayload
	err = json.Unmarshal(content, &responsePayload)
	if err != nil {
		return Message{}, err
	}

	response := responsePayload.Choices[0].Message
	chat.AppendMessage(response.Role, response.Content)

	return response, nil
}

func (chat *Chat) SendWithContext(ctx context.Context, message string) (string, error) {
	chat.AppendMessage("user", message)

	resp, err := chat.MakeRequest(ctx)
	if err != nil {
		return "", err
	}

	return resp.Content, nil
}

func (chat *Chat) Send(message string) (string, error) {
	ctx := context.Background()

	return chat.SendWithContext(ctx, message)
}

type chatStreamingResponsePayload struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	CreatedAt uint   `json:"created"`
	Model     string `json:"model"`
	Choices   []struct {
		Index        uint    `json:"index"`
		Delta        Message `json:"delta"`
		FinishReason string  `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     uint `json:"prompt_tokens"`
		CompletionTokens uint `json:"completion_tokens"`
		TotalTokens      uint `json:"total_tokens"`
	} `json:"usage"`
}

type StreamEvent struct {
	Content string
	Error   error
}

func (chat *Chat) SendStreaming(message string) (<-chan StreamEvent, error) {
	chat.AppendMessage("user", message)

	requestPayload, err := json.Marshal(chatRequestPayload{
		Model:       chat.options.Model,
		Messages:    chat.Messages,
		Temperature: chat.options.Temperature,
		MaxTokens:   uint(chat.options.MaxTokens),
		Stream:      true,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	body, err := chat.openaiClient.post(ctx, "/chat/completions", bytes.NewReader(requestPayload))
	if err != nil {
		return nil, err
	}

	respChan := make(chan StreamEvent)

	go func() {
		defer body.Close()
		scanner := bufio.NewScanner(body)

		role := ""
		response := ""
		var data string
		for scanner.Scan() {
			line := scanner.Text()

			switch {
			case strings.HasPrefix(line, "data:"):
				data += strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			case line == "":
				if data == "[DONE]" {
					break
				}
				var payload chatStreamingResponsePayload
				json.Unmarshal([]byte(data), &payload)
				if err != nil {
					respChan <- StreamEvent{Error: err}
				}
				delta := payload.Choices[0].Delta
				if delta.Role != "" {
					role = delta.Role
				}
				response += delta.Content
				respChan <- StreamEvent{Content: delta.Content}
				data = ""
			}
		}
		close(respChan)
		chat.AppendMessage(role, response)
	}()

	return respChan, nil
}

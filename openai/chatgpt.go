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

	chat.appendMessage("system", systemPrompt)

	return &chat
}

func (chat *Chat) appendMessage(role string, message string) {
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

func (chat *Chat) Send(message string) (string, error) {
	chat.appendMessage("user", message)

	requestPayload, err := json.Marshal(chatRequestPayload{
		Model:       chat.options.Model,
		Messages:    chat.Messages,
		Temperature: chat.options.Temperature,
		MaxTokens:   uint(chat.options.MaxTokens),
		Stream:      false,
	})
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	body, err := chat.openaiClient.post(ctx, "/chat/completions", bytes.NewReader(requestPayload))
	if err != nil {
		return "", err
	}

	defer body.Close()
	content, err := io.ReadAll(body)
	if err != nil {
		return "", err
	}

	var responsePayload chatResponsePayload
	json.Unmarshal(content, &responsePayload)
	if err != nil {
		return "", err
	}

	role := responsePayload.Choices[len(responsePayload.Choices)-1].Message.Role
	response := responsePayload.Choices[len(responsePayload.Choices)-1].Message.Content
	chat.appendMessage(role, response)

	return response, nil
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
	chat.appendMessage("user", message)

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
				delta := payload.Choices[len(payload.Choices)-1].Delta
				if delta.Role != "" {
					role = delta.Role
				}
				response += delta.Content
				respChan <- StreamEvent{Content: delta.Content}
				data = ""
			}
		}
		close(respChan)
		chat.appendMessage(role, response)
	}()

	return respChan, nil
}
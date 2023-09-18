package colorpicker

import (
	"context"
	_ "embed"
	"openaigo/openai"
)

//go:embed system_prompt.txt
var systemPrompt string

type ColorPicker struct {
	chatClient *openai.ChatClient
}

func New(chatClient *openai.ChatClient) *ColorPicker {
	return &ColorPicker{
		chatClient: chatClient,
	}
}

func (c *ColorPicker) GenerateColorWithContext(ctx context.Context, description string) (string, error) {
	chat := c.chatClient.NewChatWithSystemPrompt(systemPrompt)

	return chat.SendWithContext(ctx, description)
}

func (c *ColorPicker) GenerateColor(description string) (string, error) {
	chat := c.chatClient.NewChatWithSystemPrompt(systemPrompt)

	return chat.Send(description)
}

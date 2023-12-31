package colorpicker

import (
	"context"
	_ "embed"

	"github.com/actinium/openai-go-experiments/openai"
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

func (c *ColorPicker) GenerateColor(ctx context.Context, description string) (string, error) {
	chat := c.chatClient.NewChatWithSystemPrompt(systemPrompt)

	return chat.SendWithContext(ctx, description)
}

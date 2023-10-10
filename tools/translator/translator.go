package translator

import (
	"context"
	"fmt"

	"github.com/actinium/openai-go-experiments/openai"
)

type Translator struct {
	chatClient *openai.ChatClient
}

func New(chatClient *openai.ChatClient) *Translator {
	return &Translator{
		chatClient: chatClient,
	}
}

func (t *Translator) Translate(ctx context.Context, fromLanguage string, toLanguage string, text string) (string, error) {
	systemPrompt := translationPrompt(fromLanguage, toLanguage)
	chat := t.chatClient.NewChatWithSystemPrompt(systemPrompt)

	return chat.SendWithContext(ctx, text)
}

func translationPrompt(fromLanguage string, toLanguage string) string {
	return fmt.Sprintf(
		"You will be given a text in %s. Your task is to translate it into %s.",
		fromLanguage,
		toLanguage,
	)
}

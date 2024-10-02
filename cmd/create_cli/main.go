package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/actinium/openai-go-experiments/openai"
	"github.com/actinium/openai-go-experiments/setup"
)

const model string = "o1-mini"

func main() {
	setup.LoadEnv()
	openaiClient := openai.NewOpenAIClient(os.Getenv("OPENAI_APIKEY"))
	chatClient := openai.NewChatClient(openaiClient, openai.ChatOptions{
		Model:       model,
		Temperature: 1.0,
		MaxTokens:   0,
	})

	prompt := strings.Join(os.Args[1:], " ")
	if prompt == "" {
		log.Fatal("Error: empty prompt")
	}
	prompt = "Create " + prompt

	chat := chatClient.NewChat()

	reply, err := chat.Send(prompt)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(reply)
}

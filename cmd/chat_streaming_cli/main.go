package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/actinium/openai-go-experiments/openai"
	"github.com/actinium/openai-go-experiments/setup"
)

func main() {
	chatClient, _, _ := setup.Clients()

	var model string
	flag.StringVar(&model, "model", "gpt-3.5-turbo", "a ChatGPT model")

	flag.Parse()

	if model != "" {
		chatClient.UseModel(model)
	}

	startChat(chatClient.NewChat())
}

func startChat(chat *openai.Chat) {
	fmt.Print("> ")
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		userInput := s.Text()

		responseStream, err := chat.SendStreaming(userInput)
		if err != nil {
			log.Fatal(err.Error())
		}
		for part := range responseStream {
			if part.Error != nil {
				log.Fatal(part.Error.Error())
			}
			fmt.Print(part.Content)
		}

		fmt.Print("\n> ")
	}
	fmt.Println()
}

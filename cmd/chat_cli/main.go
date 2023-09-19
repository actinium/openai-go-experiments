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
	chatClient := setup.ChatClient()

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

		response, err := chat.Send(userInput)
		if err != nil {
			log.Fatal(err.Error())
		}

		fmt.Println(response)
		fmt.Print("> ")
	}
}

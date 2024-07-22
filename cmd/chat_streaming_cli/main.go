package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/actinium/openai-go-experiments/openai"
	"github.com/actinium/openai-go-experiments/setup"
	"github.com/fatih/color"
)

func main() {
	chatClient := setup.Clients().Chat()

	var model string
	flag.StringVar(&model, "model", "gpt-4o-mini", "a ChatGPT model")

	flag.Parse()

	if model != "" {
		chatClient.UseModel(model)
	}

	startChat(chatClient.NewChat())
}

func startChat(chat *openai.Chat) {
	blue := color.New(color.FgHiBlue).SprintFunc()
	fmt.Printf("%s ", blue(">"))
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

		fmt.Printf("\n%s ", blue(">"))
	}
	fmt.Println()
}

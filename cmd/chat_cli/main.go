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

		response, err := chat.Send(userInput)
		if err != nil {
			log.Fatal(err.Error())
		}

		fmt.Println(response)
		fmt.Printf("%s ", blue(">"))
	}
}

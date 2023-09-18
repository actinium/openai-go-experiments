package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/actinium/openai-go-experiments/openai"
	"github.com/actinium/openai-go-experiments/setup"
)

func main() {
	chatClient, _, _ := setup.Clients()

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

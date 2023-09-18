package main

import (
	"bufio"
	"fmt"
	"log"
	"openaigo/openai"
	"openaigo/setup"
	"os"
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

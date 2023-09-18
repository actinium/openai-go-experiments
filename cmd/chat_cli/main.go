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
	chatClient, _ := setup.Clients()

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

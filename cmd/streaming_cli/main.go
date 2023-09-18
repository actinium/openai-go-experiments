package main

import (
	"bufio"
	"fmt"
	"log"
	"openaigo/openai"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	openAIClient := openai.NewOpenAIClient(os.Getenv("OPENAI_APIKEY"))
	chatClient := openai.NewChatClient(openAIClient, openai.DefaultChatOptions)

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

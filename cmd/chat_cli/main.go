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

		response, err := chat.Send(userInput)
		if err != nil {
			log.Fatal(err.Error())
		}

		fmt.Println(response)
		fmt.Print("> ")
	}
}

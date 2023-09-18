package setup

import (
	"log"
	"openaigo/openai"
	"os"

	"github.com/joho/godotenv"
)

func Clients() (*openai.ChatClient, *openai.DalleClient) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	openAIClient := openai.NewOpenAIClient(os.Getenv("OPENAI_APIKEY"))
	chatClient := openai.NewChatClient(openAIClient, openai.DefaultChatOptions)
	imageClient := openai.NewDalleClient(openAIClient, openai.DefaultDalleOptions)

	return chatClient, imageClient
}
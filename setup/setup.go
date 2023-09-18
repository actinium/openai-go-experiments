package setup

import (
	"log"
	"os"

	"github.com/actinium/openai-go-experiments/openai"
	"github.com/joho/godotenv"
)

func Clients() (*openai.ChatClient, *openai.EmbeddingsClient, *openai.DalleClient) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	openAIClient := openai.NewOpenAIClient(os.Getenv("OPENAI_APIKEY"))
	chatClient := openai.NewChatClient(openAIClient, openai.DefaultChatOptions)
	imageClient := openai.NewDalleClient(openAIClient, openai.DefaultDalleOptions)
	embeddingsClient := openai.NewEmbeddingsClient(openAIClient)

	return chatClient, embeddingsClient, imageClient
}

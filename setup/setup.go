package setup

import (
	"log"
	"os"
	"sync"

	"github.com/actinium/openai-go-experiments/openai"
	"github.com/joho/godotenv"
)

var once sync.Once

func LoadEnv() {
	once.Do(func() {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatal("Error loading .env file")
		}
	})
}

func Clients() (*openai.ChatClient, *openai.EmbeddingsClient, *openai.DalleClient) {
	LoadEnv()

	openAIClient := openai.NewOpenAIClient(os.Getenv("OPENAI_APIKEY"))

	chatClient := openai.NewChatClient(openAIClient, openai.DefaultChatOptions)
	imageClient := openai.NewDalleClient(openAIClient, openai.DefaultDalleOptions)
	embeddingsClient := openai.NewEmbeddingsClient(openAIClient)

	return chatClient, embeddingsClient, imageClient
}

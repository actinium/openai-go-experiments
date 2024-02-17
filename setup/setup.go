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

type ClientFactory struct {
	openAIClient *openai.OpenAIClient
}

func (cf *ClientFactory) Chat() *openai.ChatClient {
	return openai.NewChatClient(cf.openAIClient, openai.DefaultChatOptions)
}

func (cf *ClientFactory) Dalle() *openai.DalleClient {
	return openai.NewDalleClient(cf.openAIClient, openai.DefaultDalleOptions)
}

func (cf *ClientFactory) Embeddigs() *openai.EmbeddingsClient {
	return openai.NewEmbeddingsClient(cf.openAIClient)
}

func (cf *ClientFactory) TTS() *openai.TextToSpeechClient {
	return openai.NewTextToSpeechClient(cf.openAIClient, openai.DefaultTextToSpeechOption)
}

func Clients() *ClientFactory {
	LoadEnv()

	openAIClient := openai.NewOpenAIClient(os.Getenv("OPENAI_APIKEY"))

	return &ClientFactory{
		openAIClient: openAIClient,
	}
}

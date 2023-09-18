package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/actinium/openai-go-experiments/openai"
	"github.com/actinium/openai-go-experiments/tools/colorpicker"
	"github.com/joho/godotenv"
)

func main() {
	prompt := strings.Join(os.Args[1:], " ")
	if prompt == "" {
		log.Fatal("Error: empty prompt")
	}

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error: couldn't load .env file")
	}
	openAIClient := openai.NewOpenAIClient(os.Getenv("OPENAI_APIKEY"))
	chatClient := openai.NewChatClient(openAIClient, openai.DefaultChatOptions)
	colorPicker := colorpicker.New(chatClient)

	color, err := colorPicker.GenerateColor(prompt)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(color)
}

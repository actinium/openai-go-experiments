package main

import (
	"fmt"
	"log"
	"openaigo/openai"
	"openaigo/tools/colorpicker"
	"os"
	"strings"

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

	color, err := colorPicker.Color(prompt)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(color)
}

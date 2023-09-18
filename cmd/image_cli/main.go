package main

import (
	"fmt"
	"log"
	"openaigo/openai"
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
	imageClient := openai.NewDalleClient(openAIClient, openai.DefaultDalleOptions)

	urls, err := imageClient.GenerateImage(prompt)
	if err != nil {
		log.Fatal("Error: couldn't generate image")
	}

	fmt.Println(urls[0])
}

package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/actinium/openai-go-experiments/openai"
	"github.com/actinium/openai-go-experiments/setup"
)

func main() {
	setup.LoadEnv()

	openAIClient := openai.NewOpenAIClient(os.Getenv("OPENAI_APIKEY"))
	models, err := openAIClient.Models()
	if err != nil {
		log.Fatal(err.Error())
	}

	sort.Slice(models, func(i, j int) bool { return models[i] < models[j] })

	for _, model := range models {
		fmt.Println(model)
	}
}

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/actinium/openai-go-experiments/setup"
)

func main() {
	imageClient := setup.Clients().Dalle()

	prompt := strings.Join(os.Args[1:], " ")
	if prompt == "" {
		log.Fatal("Error: empty prompt")
	}

	ctx := context.Background()
	urls, err := imageClient.GenerateImage(ctx, prompt)
	if err != nil {
		log.Fatal("Error: couldn't generate image")
	}

	fmt.Println(urls[0])
}

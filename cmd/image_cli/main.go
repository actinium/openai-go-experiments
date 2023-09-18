package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/actinium/openai-go-experiments/setup"
)

func main() {
	_, _, imageClient := setup.Clients()

	prompt := strings.Join(os.Args[1:], " ")
	if prompt == "" {
		log.Fatal("Error: empty prompt")
	}

	urls, err := imageClient.GenerateImage(prompt)
	if err != nil {
		log.Fatal("Error: couldn't generate image")
	}

	fmt.Println(urls[0])
}

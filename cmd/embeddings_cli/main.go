package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/actinium/openai-go-experiments/setup"
)

func main() {
	embeddings := setup.Clients().Embeddigs()

	input := strings.Join(os.Args[1:], " ")
	if input == "" {
		log.Fatal("Error: empty input")
	}

	embedding, err := embeddings.Create(input)
	if err != nil {
		log.Fatal("Error: couldn't create embedding")
	}

	fmt.Printf("%q\n%v\n", embedding.Text, embedding.Embedding)
}

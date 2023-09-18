package main

import (
	"fmt"
	"log"
	"openaigo/setup"
	"os"
	"strings"
)

func main() {
	_, embeddings, _ := setup.Clients()

	input := strings.Join(os.Args[1:], " ")
	if input == "" {
		log.Fatal("Error: empty input")
	}

	embedding, err := embeddings.Create(input)
	if err != nil {
		log.Fatal("Error: couldn't create embedding")
	}

	fmt.Println(embedding.Text)
	fmt.Println(embedding.Embedding)
}

package main

import (
	_ "embed"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/actinium/openai-go-experiments/openai"
	"github.com/actinium/openai-go-experiments/setup"
	"github.com/fatih/color"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed index.html
var page []byte

func indexHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Write(page)
	}
}

//go:embed system_prompt.txt
var systemPrompt string

func chatHandler(chatClient *openai.ChatClient) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var messages []openai.Message

		err := json.NewDecoder(req.Body).Decode(&messages)
		if err != nil {
			http.Error(w, "Could not parse request body", http.StatusBadRequest)
			return
		}

		chat := chatClient.NewChatWithSystemPrompt(systemPrompt)

		for _, message := range messages {
			if message.Role == "" || message.Content == "" {
				http.Error(w, "Empty message or role", http.StatusBadRequest)
				return
			}
			chat.AppendMessage(message.Role, message.Content)
		}

		response, err := chat.MakeRequest(req.Context())
		if err != nil {
			http.Error(w, "Could not generate response", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "Could not generate response", http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	setup.LoadEnv()

	chatClient := setup.Clients().Chat()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", indexHandler())
	r.Post("/chat", chatHandler(chatClient))

	link := color.New(color.Underline, color.FgHiBlue).SprintFunc()
	log.Printf("Listening on %s\n", link(os.Getenv("HTTP_ADDR")))
	http.ListenAndServe(os.Getenv("HTTP_ADDR"), r)
}

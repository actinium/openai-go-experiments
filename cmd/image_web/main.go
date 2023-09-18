package main

import (
	_ "embed"
	"encoding/json"
	"net/http"
	"os"

	"github.com/actinium/openai-go-experiments/openai"
	"github.com/actinium/openai-go-experiments/setup"
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

func imagineHandler(imageClient *openai.DalleClient) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if !req.URL.Query().Has("prompt") {
			http.Error(w, "Missing prompt", 400)
			return
		}

		prompt := req.URL.Query().Get("prompt")
		if prompt == "" {
			http.Error(w, "Prompt is empty", 400)
			return
		}

		ctx := req.Context()
		urls, err := imageClient.GenerateImageWithContext(ctx, prompt)
		if err != nil {
			http.Error(w, "Could not generate image", 500)
			return
		}

		b, err := json.Marshal(&urls)
		if err != nil {
			http.Error(w, "Could not generate JSON response", 500)
			return
		}

		w.Write(b)
	}
}

func main() {
	setup.LoadEnv()

	openAIClient := openai.NewOpenAIClient(os.Getenv("OPENAI_APIKEY"))
	options := openai.DefaultDalleOptions
	options.N = 4
	imageClient := openai.NewDalleClient(openAIClient, options)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", indexHandler())
	r.Get("/imagine", imagineHandler(imageClient))

	http.ListenAndServe(os.Getenv("HTTP_ADDR"), r)
}

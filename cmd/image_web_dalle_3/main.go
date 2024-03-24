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

func imagineHandler(imageClient *openai.DalleClient) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		prompt := req.FormValue("prompt")
		if prompt == "" {
			http.Error(w, "Missing or empty prompt", 400)
			return
		}

		ctx := req.Context()
		urls, err := imageClient.GenerateImage(ctx, prompt)
		if err != nil {
			http.Error(w, "Could not generate image", 500)
			return
		}
		url := urls[0]

		b, err := json.Marshal(&url)
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
	options.Model = "dall-e-3"
	options.N = 1
	imageClient := openai.NewDalleClient(openAIClient, options)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", indexHandler())
	r.Post("/imagine", imagineHandler(imageClient))

	link := color.New(color.Underline, color.FgHiBlue).SprintFunc()
	log.Printf("Listening on %s\n", link(os.Getenv("HTTP_ADDR")))
	http.ListenAndServe(os.Getenv("HTTP_ADDR"), r)
}
